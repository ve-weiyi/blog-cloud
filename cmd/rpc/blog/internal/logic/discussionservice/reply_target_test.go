package discussionservicelogic

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/model"
)

type fakeFinder struct {
	rows map[int64]*model.TComment
}

func (f fakeFinder) FindById(_ context.Context, id int64) (*model.TComment, error) {
	row, ok := f.rows[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return row, nil
}

// 优先按 ReplyId 反查；没有回复对象时退到父评论作者。
func TestResolveReplyUser(t *testing.T) {
	finder := fakeFinder{rows: map[int64]*model.TComment{
		11: {Id: 11, UserId: "replied-user"},
		22: {Id: 22, UserId: "parent-author"},
	}}

	ctx := context.Background()

	got, err := resolveReplyUser(ctx, finder, 11, 22)
	if err != nil || got != "replied-user" {
		t.Fatalf("ReplyId 优先：got=%q err=%v", got, err)
	}

	got, err = resolveReplyUser(ctx, finder, 0, 22)
	if err != nil || got != "parent-author" {
		t.Fatalf("缺省取父评论作者：got=%q err=%v", got, err)
	}

	got, err = resolveReplyUser(ctx, finder, 0, 0)
	if err != nil || got != "" {
		t.Fatalf("无回复对象时应为空：got=%q err=%v", got, err)
	}
}

// 指向不存在的评论必须报错，而不是静默留空——静默会让收件人变成「无人」。
func TestResolveReplyUser_TargetMissing(t *testing.T) {
	_, err := resolveReplyUser(context.Background(), fakeFinder{}, 99, 0)
	if !errors.Is(err, errReplyTargetNotFound) {
		t.Fatalf("期望 errReplyTargetNotFound，得到 %v", err)
	}
}
