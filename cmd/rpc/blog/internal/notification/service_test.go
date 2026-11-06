package notification_test

import (
	"context"
	"testing"

	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/notification"
	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
)

// Emit 在落库之前完成全部校验：无事务、未知场景、非法受众都在这里拦住。
func TestEmit_RejectsBeforeWrite(t *testing.T) {
	ctx := context.Background()
	svc := notification.NewService(notification.Deps{})
	tx := &gorm.DB{}

	if err := svc.Emit(ctx, nil, notification.Request{Scene: notify.SceneArticleLiked, Audience: notification.AudienceUserIDs}); err == nil {
		t.Error("缺少事务应报错")
	}
	if err := svc.Emit(ctx, tx, notification.Request{Scene: "unknown.scene", Audience: notification.AudienceAll}); err == nil {
		t.Error("未知场景应报错")
	}
	if err := svc.Emit(ctx, tx, notification.Request{Scene: notify.SceneArticleLiked, Audience: ""}); err == nil {
		t.Error("空受众应报错")
	}
	if err := svc.Emit(ctx, tx, notification.Request{Scene: notify.SceneArticleLiked, Audience: notification.AudienceUserIDs}); err == nil {
		t.Error("user_ids 受众为空应报错")
	}
}

// 只有终态事件允许清理，进行中的事件不能回收。
func TestIsTerminal(t *testing.T) {
	for _, s := range []string{notification.StatusPending, notification.StatusFanning} {
		if notification.IsTerminal(s) {
			t.Errorf("状态 %q 不应视为终态", s)
		}
	}
	for _, s := range []string{notification.StatusDone, notification.StatusFailed, notification.StatusTerminated} {
		if !notification.IsTerminal(s) {
			t.Errorf("状态 %q 应视为终态", s)
		}
	}
}
