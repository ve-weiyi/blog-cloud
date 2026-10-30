package notification_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/notification"
	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
	"github.com/ve-weiyi/blog-cloud/model"
)

// 未读计数对外按偏好分组下发：三个点赞场景要合并成 like 一个键，
// 否则前端查 group_counts["like"] 永远是 0。
func TestCountUnread_GroupsByPreferenceGroup(t *testing.T) {
	gdb, mock := newMockDB(t)

	// like + comment 两个分组展开成 5 个场景，参数个数与内容由目录决定
	mock.ExpectQuery(regexp.QuoteMeta("SELECT scene AS scene, COUNT(DISTINCT group_key) AS cnt FROM `t_notification` WHERE recipient_id = ? AND scene IN (?,?,?,?,?) AND read_at IS NULL GROUP BY `scene`")).
		WithArgs("u1", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"scene", "cnt"}).
			AddRow("article.liked", 2).
			AddRow("comment.replied", 1))

	svc := notification.NewService(notification.Deps{DB: gdb, Notifications: model.NewTNotificationModel(gdb)})
	total, counts, err := svc.CountUnread(context.Background(), "u1",
		[]notify.Group{notify.GroupLike, notify.GroupComment})
	if err != nil {
		t.Fatalf("CountUnread: %v", err)
	}

	if total != 3 {
		t.Errorf("未读总数 = %d，期望 3", total)
	}
	if counts[notify.GroupLike] != 2 {
		t.Errorf("点赞分组 = %d，期望 2；实际 %v", counts[notify.GroupLike], counts)
	}
	if counts[notify.GroupComment] != 1 {
		t.Errorf("评论分组 = %d，期望 1；实际 %v", counts[notify.GroupComment], counts)
	}
	if _, ok := counts[notify.Group("article.liked")]; ok {
		t.Errorf("不应以场景为键下发：%v", counts)
	}
}
