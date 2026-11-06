package job

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// AcquireLock 获取命名锁，返回释放函数；同一后台任务在多副本中单活。
//
// 命名锁是连接级的：必须在同一个 *sql.Conn 上获取与释放。走连接池的话，
// RELEASE_LOCK 可能落到另一条连接，锁会一直留在原连接上。
func AcquireLock(ctx context.Context, db *gorm.DB, name string) (func(), bool, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, false, fmt.Errorf("job: get sql db: %w", err)
	}
	conn, err := sqlDB.Conn(ctx)
	if err != nil {
		return nil, false, fmt.Errorf("job: pin lock connection: %w", err)
	}

	var acquired sql.NullInt64
	if err := conn.QueryRowContext(ctx, "SELECT GET_LOCK(?, 0)", name).Scan(&acquired); err != nil {
		_ = conn.Close()
		return nil, false, fmt.Errorf("job: acquire %s: %w", name, err)
	}
	if !acquired.Valid || acquired.Int64 != 1 {
		_ = conn.Close()
		return nil, false, nil
	}

	release := func() {
		// 释放不跟随任务 ctx：任务被取消时仍要把锁还回去
		releaseCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _ = conn.ExecContext(releaseCtx, "SELECT RELEASE_LOCK(?)", name)
		_ = conn.Close()
	}
	return release, true, nil
}
