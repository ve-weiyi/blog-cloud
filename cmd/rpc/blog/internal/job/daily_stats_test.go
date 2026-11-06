package job

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

// 聚合结果写入用户指标，且不覆盖 UV/PV。
func TestRunOnce_UpsertsUserCounts(t *testing.T) {
	db, mock := newSQLMockDB(t)
	stats := NewDailyStats(db)
	day := time.Date(2026, 11, 2, 0, 0, 0, 0, time.Local)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT GET_LOCK(?, 0)")).
		WithArgs("blog:daily-stats").
		WillReturnRows(sqlmock.NewRows([]string{"acquired"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `t_user` WHERE created_at >= ? AND created_at < ?")).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `t_user` WHERE created_at < ?")).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(100))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(DISTINCT")).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(7))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT total_uv_count, total_pv_count FROM `t_daily_stats` WHERE date < ?")).
		WillReturnRows(sqlmock.NewRows([]string{"total_uv_count", "total_pv_count"}).AddRow(10, 20))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `t_daily_stats`")).
		WithArgs(sqlmock.AnyArg(), 3, 100, 7, 10, 20).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("SELECT RELEASE_LOCK(?)")).
		WithArgs("blog:daily-stats").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := stats.RunOnce(context.Background(), day); err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 未按预期执行: %v", err)
	}
}

// 到达执行时刻后顺延到第二天。
func TestNextDaily(t *testing.T) {
	loc := time.Local
	now := time.Date(2026, 11, 3, 0, 4, 0, 0, loc)
	if got, want := nextDaily(now, 0, 5), time.Date(2026, 11, 3, 0, 5, 0, 0, loc); !got.Equal(want) {
		t.Fatalf("nextDaily = %v, want %v", got, want)
	}

	now = time.Date(2026, 11, 3, 0, 5, 0, 0, loc)
	if got, want := nextDaily(now, 0, 5), time.Date(2026, 11, 4, 0, 5, 0, 0, loc); !got.Equal(want) {
		t.Fatalf("nextDaily = %v, want %v", got, want)
	}
}
