package notification

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
)

// 过期条目的投递必须与条目在同一事务里删除，不能留下孤儿投递。
func TestDeleteExpiredBatch_DeletesEntriesWithDeliveries(t *testing.T) {
	svc, mock, _ := newWorkerDeps(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM `t_notification` WHERE scene = ? AND (read_at IS NOT NULL AND read_at < ?)")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `t_notification_delivery` WHERE notification_id IN")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `t_notification` WHERE id IN")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	n, err := svc.deleteExpiredBatch(context.Background(), "article.liked", true, time.Now().Add(-90*24*time.Hour), 10)
	if err != nil {
		t.Fatalf("deleteExpiredBatch: %v", err)
	}
	if n != 1 {
		t.Fatalf("删除 %d 条，期望 1 条", n)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 未按预期执行: %v", err)
	}
}

// 未读条目按创建时间过期，条件不能落到 read_at 上。
func TestDeleteExpiredBatch_UnreadUsesCreatedAt(t *testing.T) {
	svc, mock, _ := newWorkerDeps(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM `t_notification` WHERE scene = ? AND (read_at IS NULL AND created_at < ?)")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	// 空批次不应开启事务

	n, err := svc.deleteExpiredBatch(context.Background(), "article.liked", false, time.Now().Add(-180*24*time.Hour), 10)
	if err != nil {
		t.Fatalf("deleteExpiredBatch: %v", err)
	}
	if n != 0 {
		t.Fatalf("空批次删除 %d 条，期望 0 条", n)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 未按预期执行: %v", err)
	}
}

// 只回收终态且已无引用条目的旧事件。
func TestPruneEvents_TerminalOnlyAndUnreferenced(t *testing.T) {
	svc, mock, _ := newWorkerDeps(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM `t_notification_event` WHERE status IN")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(7).AddRow(8))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `t_notification_event` WHERE id IN")).
		WillReturnResult(sqlmock.NewResult(0, 2))

	n, err := svc.pruneEvents(context.Background(), time.Now().Add(-7*24*time.Hour), 10)
	if err != nil {
		t.Fatalf("pruneEvents: %v", err)
	}
	if n != 2 {
		t.Fatalf("删除 %d 条，期望 2 条", n)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 未按预期执行: %v", err)
	}
}

// 场景遍历使用 Key 作为 scene 值，并为已读/未读各跑一次。
func TestRunRetention_UsesSceneKeyAndBothTTLs(t *testing.T) {
	svc, mock, _ := newWorkerDeps(t)
	scene := notify.Scene{Key: "article.liked", ReadTTL: 24 * time.Hour, UnreadTTL: 48 * time.Hour}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM `t_notification` WHERE scene = ? AND (read_at IS NOT NULL AND read_at < ?)")).
		WithArgs("article.liked", sqlmock.AnyArg(), 10).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM `t_notification` WHERE scene = ? AND (read_at IS NULL AND created_at < ?)")).
		WithArgs("article.liked", sqlmock.AnyArg(), 10).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM `t_notification_event` WHERE status IN")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	entries, events, err := svc.runRetention(context.Background(), time.Now(), 10, []notify.Scene{scene})
	if err != nil {
		t.Fatalf("runRetention: %v", err)
	}
	if entries != 0 || events != 0 {
		t.Fatalf("空库删除条目 %d、事件 %d，期望均为 0", entries, events)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 未按预期执行: %v", err)
	}
}

// 一批删满后必须继续下一批，不足一批才停。
func TestDrainScene_ContinuesUntilShortBatch(t *testing.T) {
	svc, mock, _ := newWorkerDeps(t)
	cutoff := time.Now().Add(-24 * time.Hour)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM `t_notification` WHERE scene = ? AND (read_at IS NOT NULL AND read_at < ?)")).
		WithArgs("article.liked", sqlmock.AnyArg(), 2).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `t_notification_delivery` WHERE notification_id IN")).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `t_notification` WHERE id IN")).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectCommit()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM `t_notification` WHERE scene = ? AND (read_at IS NOT NULL AND read_at < ?)")).
		WithArgs("article.liked", sqlmock.AnyArg(), 2).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `t_notification_delivery` WHERE notification_id IN")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `t_notification` WHERE id IN")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	n, err := svc.drainScene(context.Background(), "article.liked", true, cutoff, 2)
	if err != nil {
		t.Fatalf("drainScene: %v", err)
	}
	if n != 3 {
		t.Fatalf("删除 %d 条，期望 3 条", n)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 未按预期执行: %v", err)
	}
}
