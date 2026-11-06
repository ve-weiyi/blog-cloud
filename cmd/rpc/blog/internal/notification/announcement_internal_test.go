package notification

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/ve-weiyi/blog-cloud/model"
)

// 广播只在事务提交之后发出：事务回滚时管理端不该看到任何状态变化。
func TestPublishAnnouncement_RollbackDoesNotBroadcast(t *testing.T) {
	gdb, mock := newSQLMockDB(t)
	outbox := &fakeEnqueuer{}
	announcer := &fakeAnnouncer{}
	svc := NewService(Deps{
		DB:                   gdb,
		Events:               model.NewTNotificationEventModel(gdb),
		Announcements:        model.NewTAnnouncementModel(gdb),
		Outbox:               outbox,
		AnnouncementNotifier: announcer,
	})

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `t_announcement` WHERE `id` = ?")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "target_type", "status"}).
			AddRow(9, "维护通知", AudienceAll, AnnouncementDraft))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `t_announcement` SET")).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `t_notification_event`")).
		WillReturnError(errors.New("insert failed"))
	mock.ExpectRollback()

	published, err := svc.PublishAnnouncement(context.Background(), 9, "admin")
	if err == nil {
		t.Fatal("事件写入失败时应返回错误")
	}
	if published {
		t.Fatal("事务失败时不应报告已发布")
	}
	if announcer.calls != 0 {
		t.Fatalf("事务回滚时不应广播，实际广播 %d 次", announcer.calls)
	}
	if outbox.enqueued != 0 {
		t.Fatalf("事务回滚时不应入队唤醒消息，实际 %d 条", outbox.enqueued)
	}
}
