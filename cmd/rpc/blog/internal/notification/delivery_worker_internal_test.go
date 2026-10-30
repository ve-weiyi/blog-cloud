package notification

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
	"github.com/ve-weiyi/blog-cloud/model"
)

// recordingSender 记录是否真的发出了邮件。
type recordingSender struct {
	sent int
}

func (s *recordingSender) SendEmail(_ context.Context, _, _, _ string) error {
	s.sent++
	return nil
}

func newWorkerDeps(t *testing.T) (*Service, sqlmock.Sqlmock, *recordingSender) {
	t.Helper()

	gdb, mock := newSQLMockDB(t)

	svc := NewService(Deps{
		DB:            gdb,
		Events:        model.NewTNotificationEventModel(gdb),
		Notifications: model.NewTNotificationModel(gdb),
		Deliveries:    model.NewTNotificationDeliveryModel(gdb),
		Preferences:   model.NewTNotificationPreferenceModel(gdb),
		Announcements: model.NewTAnnouncementModel(gdb),
		Users:         model.NewTUserModel(gdb),
		Comments:      model.NewTCommentModel(gdb),
	})
	return svc, mock, &recordingSender{}
}

func pendingDelivery() *model.TNotificationDelivery {
	return &model.TNotificationDelivery{
		Id:             9,
		NotificationId: 5,
		Channel:        "email",
		Recipient:      "a@b.c",
		Status:         DeliveryPending,
		Attempts:       1,
	}
}

func pendingNotif() *model.TNotification {
	return &model.TNotification{Id: 5, EventId: 7, RecipientId: "u1", Scene: notify.SceneArticleLiked}
}

func eventOf(status string) *model.TNotificationEvent {
	return &model.TNotificationEvent{Id: 7, Scene: notify.SceneArticleLiked, Status: status}
}

// announcementNotif 构造一条公告条目：公告是强制分组，投递时不该查偏好。
func announcementNotif() (*model.TNotification, *model.TNotificationEvent) {
	payload := `{"announcement_id":9}`
	n := &model.TNotification{Id: 5, EventId: 7, RecipientId: "u1", Scene: notify.SceneAnnouncement}
	ev := &model.TNotificationEvent{Id: 7, Scene: notify.SceneAnnouncement, Status: StatusDone, Payload: &payload}
	return n, ev
}

// expectNoPreference 返回空偏好（沿用场景默认）。
func expectNoPreference(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `t_notification_preference` WHERE user_id = ? AND pref_group = ? AND channel IN")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "pref_group", "channel", "enabled"}))
}

// expectStatusUpdate 匹配一次投递记录状态回写。
func expectStatusUpdate(mock sqlmock.Sqlmock) {
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `t_notification_delivery` SET")).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

// 扇出完成（done）是正常可投递状态，邮件必须真的发出去。
func TestProcess_EventDoneIsDeliverable(t *testing.T) {
	svc, mock, sender := newWorkerDeps(t)
	expectNoPreference(mock)
	expectStatusUpdate(mock)

	worker := NewDeliveryWorker(svc, sender, 0)
	if err := worker.process(context.Background(), pendingDelivery(), pendingNotif(), eventOf(StatusDone)); err != nil {
		t.Fatalf("process: %v", err)
	}
	if sender.sent != 1 {
		t.Fatalf("事件 done 时应发送邮件，实际发送 %d 次", sender.sent)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 未按预期执行: %v", err)
	}
}

// 已终止的事件不该再发。
func TestProcess_EventTerminatedIsCancelled(t *testing.T) {
	svc, mock, sender := newWorkerDeps(t)
	// 已终止：不查偏好、直接置取消
	expectStatusUpdate(mock)

	worker := NewDeliveryWorker(svc, sender, 0)
	if err := worker.process(context.Background(), pendingDelivery(), pendingNotif(), eventOf(StatusTerminated)); err != nil {
		t.Fatalf("process: %v", err)
	}
	if sender.sent != 0 {
		t.Fatalf("事件已终止时不应发送，实际发送 %d 次", sender.sent)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 未按预期执行: %v", err)
	}
}

// 公告内容已删：取消而不是重试——重试永远不会成功。
func TestProcess_AnnouncementGoneIsCancelled(t *testing.T) {
	svc, mock, sender := newWorkerDeps(t)
	notif, ev := announcementNotif()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `t_announcement` WHERE `id` = ?")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}))
	expectStatusUpdate(mock)

	worker := NewDeliveryWorker(svc, sender, 0)
	if err := worker.process(context.Background(), pendingDelivery(), notif, ev); err != nil {
		t.Fatalf("process: %v", err)
	}
	if sender.sent != 0 {
		t.Fatalf("公告已删除时不应发送，实际发送 %d 次", sender.sent)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 未按预期执行: %v", err)
	}
}

// 强制分组的通知不看偏好：脏数据里的关闭行不能拦下公告。
func TestProcess_MandatoryIgnoresPreference(t *testing.T) {
	svc, mock, sender := newWorkerDeps(t)
	notif, ev := announcementNotif()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `t_announcement` WHERE `id` = ?")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content"}).AddRow(9, "维护通知", "正文"))
	// 没有偏好查询的预期：一旦查询就会因「unexpected call」失败
	expectStatusUpdate(mock)

	worker := NewDeliveryWorker(svc, sender, 0)
	if err := worker.process(context.Background(), pendingDelivery(), notif, ev); err != nil {
		t.Fatalf("process: %v", err)
	}
	if sender.sent != 1 {
		t.Fatalf("强制分组应忽略偏好并发送，实际发送 %d 次", sender.sent)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 未按预期执行: %v", err)
	}
}

// 批量预取：一批投递最多两次查询，不随条数增长。
func TestLoadDeliveryTargets_BatchesQueries(t *testing.T) {
	svc, mock, _ := newWorkerDeps(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `t_notification` WHERE id IN")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "event_id", "recipient_id", "scene"}).
			AddRow(5, 7, "u1", "article.liked").
			AddRow(6, 8, "u2", "article.liked"))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `t_notification_event` WHERE id IN")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "scene", "status"}).
			AddRow(7, "article.liked", StatusDone).
			AddRow(8, "article.liked", StatusDone))

	rows := []*model.TNotificationDelivery{
		{Id: 1, NotificationId: 5},
		{Id: 2, NotificationId: 5},
		{Id: 3, NotificationId: 6},
	}
	notifs, events, err := svc.loadDeliveryTargets(context.Background(), rows)
	if err != nil {
		t.Fatalf("loadDeliveryTargets: %v", err)
	}
	if len(notifs) != 2 || len(events) != 2 {
		t.Fatalf("条目 %d 条、事件 %d 条，期望各 2 条", len(notifs), len(events))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 未按预期执行: %v", err)
	}
}

// 判据本身：done/failed 可投递，只有 terminated 不可。
func TestUndeliverable(t *testing.T) {
	if undeliverable(StatusDone) || undeliverable(StatusFailed) || undeliverable(StatusFanning) {
		t.Error("done / failed / fanning 都应允许投递")
	}
	if !undeliverable(StatusTerminated) {
		t.Error("terminated 应拒绝投递")
	}
}
