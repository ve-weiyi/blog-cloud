package job

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func newSQLMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
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

	return db, mock
}

// 一批删满后继续下一批，不足一批才停。
func TestRunOnce_DeletesUntilShortBatch(t *testing.T) {
	db, mock := newSQLMockDB(t)
	retention := NewLogRetention(db, time.Hour)
	policies := []logPolicy{{table: "`t_visit_log`", keep: 90 * day}}

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `t_visit_log` WHERE created_at < ? ORDER BY created_at LIMIT ?")).
		WithArgs(sqlmock.AnyArg(), 2).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `t_visit_log` WHERE created_at < ? ORDER BY created_at LIMIT ?")).
		WithArgs(sqlmock.AnyArg(), 2).
		WillReturnResult(sqlmock.NewResult(0, 1))

	deleted, err := retention.RunOnce(context.Background(), time.Now(), policies, 2)
	if err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if deleted != 3 {
		t.Fatalf("删除 %d 行，期望 3 行", deleted)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 未按预期执行: %v", err)
	}
}

// 删除失败要回传错误，不能静默跳过整张表。
func TestRunOnce_DeleteError(t *testing.T) {
	db, mock := newSQLMockDB(t)
	retention := NewLogRetention(db, time.Hour)
	policies := []logPolicy{{table: "`t_operation_log`", keep: 365 * day}}

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `t_operation_log` WHERE created_at < ? ORDER BY created_at LIMIT ?")).
		WillReturnError(errors.New("db down"))

	if _, err := retention.RunOnce(context.Background(), time.Now(), policies, 10); err == nil {
		t.Fatal("删除失败时应返回错误")
	}
}
