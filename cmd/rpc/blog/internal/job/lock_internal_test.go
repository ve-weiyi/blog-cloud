package job

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// 锁在获取它的连接上释放。
func TestAcquireLock_ReleasesLock(t *testing.T) {
	db, mock := newSQLMockDB(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT GET_LOCK(?, 0)")).
		WithArgs("test-lock").
		WillReturnRows(sqlmock.NewRows([]string{"acquired"}).AddRow(1))
	mock.ExpectExec(regexp.QuoteMeta("SELECT RELEASE_LOCK(?)")).
		WithArgs("test-lock").
		WillReturnResult(sqlmock.NewResult(0, 1))

	release, ok, err := AcquireLock(context.Background(), db, "test-lock")
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	if !ok || release == nil {
		t.Fatal("锁应获取成功")
	}
	release()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 未按预期执行: %v", err)
	}
}

// 别的实例持有锁时直接跳过，不执行清理。
func TestAcquireLock_SkipsWhenHeld(t *testing.T) {
	db, mock := newSQLMockDB(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT GET_LOCK(?, 0)")).
		WithArgs("test-lock").
		WillReturnRows(sqlmock.NewRows([]string{"acquired"}).AddRow(0))

	release, ok, err := AcquireLock(context.Background(), db, "test-lock")
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	if ok || release != nil {
		t.Fatal("锁被占用时应返回 false")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 未按预期执行: %v", err)
	}
}
