package discussionservicelogic

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
)

type ReviewCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReviewCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReviewCommentLogic {
	return &ReviewCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 审核评论
//
// 通知的触发点是「审核通过」这个状态迁移：评论在通过之前对读者不可见，
// 提前发出通知会指向一条看不到的内容。只有状态真的变化才发事件，重复审核是空操作。
func (l *ReviewCommentLogic) ReviewComment(in *discussionrpc.ReviewCommentRequest) (*discussionrpc.ReviewCommentResponse, error) {
	if in.IsReview != enums.CommentIsReviewNo && in.IsReview != enums.CommentIsReviewYes {
		return nil, bizerr.NewBizError(bizcode.CodeParamValueInvalid, "审核状态不合法")
	}

	comment, err := l.svcCtx.TCommentModel.FindById(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerr.NewBizError(bizcode.CodeResourceNotFound, "评论不存在")
		}
		return nil, err
	}
	if comment.IsReview == in.IsReview {
		return &discussionrpc.ReviewCommentResponse{Success: true}, nil
	}

	err = l.svcCtx.GormDB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		// WHERE 带上原状态：并发审核由影响行数裁决
		rows, err := l.svcCtx.TCommentModel.WithTx(tx).UpdateFields(l.ctx,
			map[string]interface{}{"is_review": in.IsReview},
			"id = ? AND is_review <> ?", in.Id, in.IsReview)
		if err != nil {
			return err
		}
		if rows == 0 {
			return nil
		}
		return emitReviewNotifications(l.ctx, l.svcCtx, tx, comment, in.IsReview == enums.CommentIsReviewYes)
	})
	if err != nil {
		return nil, err
	}

	return &discussionrpc.ReviewCommentResponse{Success: true}, nil
}
