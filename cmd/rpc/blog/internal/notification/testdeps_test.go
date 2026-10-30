package notification

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	mqx "github.com/ve-weiyi/vkit/adapter/mqx"
)

// newSQLMockDB 返回一个不连真实库的 gorm 实例。
//
// 关掉默认事务：单条读写不需要外层事务，省掉 BEGIN/COMMIT 让用例只看语句本身。
func newSQLMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		TranslateError:         true,
		SkipDefaultTransaction: true,
		NamingStrategy:         schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		t.Fatalf("gorm: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	return gdb, mock
}

// fakeEnqueuer 记录唤醒消息是否入队。
type fakeEnqueuer struct {
	enqueued int
}

func (e *fakeEnqueuer) Enqueue(_ context.Context, _ *gorm.DB, _ mqx.Message) error {
	e.enqueued++
	return nil
}

// fakeAnnouncer 记录管理端广播次数。
type fakeAnnouncer struct {
	calls int
}

func (a *fakeAnnouncer) NotifyAnnouncement(_ context.Context, _ bool, _ int64, _, _, _ string, _ time.Time) {
	a.calls++
}
