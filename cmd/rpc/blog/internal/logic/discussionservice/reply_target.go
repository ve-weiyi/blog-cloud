package discussionservicelogic

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/model"
)

// errReplyTargetNotFound 表示请求指向了一条不存在的评论。
var errReplyTargetNotFound = errors.New("discussionservice: reply target not found")

// commentFinder 是反查被回复评论所需的最小能力，便于单测替换。
type commentFinder interface {
	FindById(ctx context.Context, id int64) (*model.TComment, error)
}

// resolveReplyUser 由服务端反查被回复者：优先 ReplyId，缺省 ParentId。
//
// 请求体里的 reply_user_id 一律不采信——它同时用于通知收件人与「回复 @某人」的展示，
// 采信客户端传入值等于把收件关系交给调用方决定。
func resolveReplyUser(ctx context.Context, finder commentFinder, replyID, parentID int64) (string, error) {
	target := replyID
	if target <= 0 {
		target = parentID
	}
	if target <= 0 {
		return "", nil
	}

	row, err := finder.FindById(ctx, target)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", errReplyTargetNotFound
	}
	if err != nil {
		return "", err
	}
	return row.UserId, nil
}
