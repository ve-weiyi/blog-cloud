package notification

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/ve-weiyi/blog-cloud/model"
	"github.com/ve-weiyi/vkit/adapter/mqx"
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

// newNotifyService 装配测试用的通知服务；outbox 为 nil 时表示未配置发件箱。
func newNotifyService(t *testing.T, outbox Enqueuer) (*Service, sqlmock.Sqlmock) {
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
		Outbox:        outbox,
	})
	return svc, mock
}

// fakeEnqueuer 记录唤醒消息是否入队；err 非空时模拟入队失败。
type fakeEnqueuer struct {
	enqueued int
	err      error
}

func (e *fakeEnqueuer) Enqueue(_ context.Context, _ *gorm.DB, _ mqx.Message) error {
	if e.err != nil {
		return e.err
	}
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
