package notification

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func newReconcileDeps(t *testing.T) (*Service, sqlmock.Sqlmock, *fakeEnqueuer) {
	t.Helper()

	outbox := &fakeEnqueuer{}
	svc, mock := newNotifyService(t, outbox)
	return svc, mock, outbox
}

// 停在待扇出/扇出中的事件要补发唤醒；补发后刷新时间，避免每轮重复入队。
func TestReconcileStalled_NudgesStalledEvents(t *testing.T) {
	svc, mock, outbox := newReconcileDeps(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `t_notification_event` WHERE status IN")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).
			AddRow(1, StatusPending).
			AddRow(2, StatusFanning))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `t_notification_event` SET")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `t_notification_event` SET")).
		WillReturnResult(sqlmock.NewResult(0, 1))

	n, err := svc.ReconcileStalled(context.Background(), time.Now(), time.Minute, 10)
	if err != nil {
		t.Fatalf("ReconcileStalled: %v", err)
	}
	if n != 2 {
		t.Fatalf("补发 %d 条，期望 2 条", n)
	}
	if outbox.enqueued != 2 {
		t.Fatalf("入队 %d 条，期望 2 条", outbox.enqueued)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 未按预期执行: %v", err)
	}
}

// 查询失败要显式返回，不能静默当作「没有停滞事件」。
func TestReconcileStalled_QueryError(t *testing.T) {
	svc, mock, _ := newReconcileDeps(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `t_notification_event` WHERE status IN")).
		WillReturnError(errors.New("db down"))

	if _, err := svc.ReconcileStalled(context.Background(), time.Now(), time.Minute, 10); err == nil {
		t.Fatal("查询失败时应返回错误")
	}
}

// 全部入队失败必须回传错误，不能把系统性故障伪装成「没有停滞事件」。
func TestReconcileStalled_AllEnqueuesFail(t *testing.T) {
	svc, mock, _ := newReconcileDeps(t)
	svc.outbox = &fakeEnqueuer{err: errors.New("mq down")}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `t_notification_event` WHERE status IN")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(1, StatusFanning))

	n, err := svc.ReconcileStalled(context.Background(), time.Now(), time.Minute, 10)
	if err == nil {
		t.Fatal("全部入队失败时应返回错误")
	}
	if n != 0 {
		t.Fatalf("补发 %d 条，期望 0 条", n)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 未按预期执行: %v", err)
	}
}

// 没有 outbox 就没有补发通道，应明确失败而不是假装完成。
func TestReconcileStalled_RequiresOutbox(t *testing.T) {
	svc, _, _ := newWorkerDeps(t)

	if _, err := svc.ReconcileStalled(context.Background(), time.Now(), time.Minute, 10); err == nil {
		t.Fatal("缺少 outbox 时应返回错误")
	}
}
