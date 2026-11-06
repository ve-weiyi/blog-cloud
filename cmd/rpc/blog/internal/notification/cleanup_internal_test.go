package notification

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
)

// 删主题时连带回收其下评论产生的条目：主题没了，评论仍在库里，但页面已不可达。
func TestTerminateTopic_CascadesToCommentEntries(t *testing.T) {
	svc, mock, _ := newWorkerDeps(t)
	ctx := context.Background()

	// 主题自身的条目为空：这条用例只看评论那一轮
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `t_notification_event` SET")).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `t_notification` WHERE object_type = ? AND object_id IN")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	// 评论对象：子查询下推，命中一条条目，投递要跟着一起删
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `t_notification_event` SET")).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `t_notification` WHERE object_type = ? AND object_id IN (SELECT id FROM")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `t_notification_delivery` WHERE notification_id IN")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `t_notification` WHERE id IN")).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := svc.TerminateTopic(ctx, svc.db, notify.ObjectArticle, []int64{3}); err != nil {
		t.Fatalf("TerminateTopic: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 未按预期执行: %v", err)
	}
}

// 撤回一次互动只删它自己那条条目，投递一并删除，不留孤儿。
func TestRetractInteraction_DeletesEntryWithDelivery(t *testing.T) {
	svc, mock, _ := newWorkerDeps(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `t_notification` WHERE dedupe_key = ?")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `t_notification_delivery` WHERE notification_id IN")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `t_notification` WHERE id IN")).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := svc.RetractInteraction(context.Background(), svc.db,
		notify.SceneContentCommented, notify.ObjectArticle, 3, "u1")
	if err != nil {
		t.Fatalf("RetractInteraction: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 未按预期执行: %v", err)
	}
}

// 不按对象分组就没有「一次互动」的定位语义，直接跳过。
func TestRetractInteraction_SkipsNonGroupedScene(t *testing.T) {
	svc, _, _ := newWorkerDeps(t)

	// 没有任何 SQL 预期：一旦查询就会因 unexpected call 失败
	err := svc.RetractInteraction(context.Background(), svc.db,
		notify.SceneCommentReplied, notify.ObjectComment, 7, "u1")
	if err != nil {
		t.Fatalf("不分组场景应直接跳过，得到 %v", err)
	}
}
