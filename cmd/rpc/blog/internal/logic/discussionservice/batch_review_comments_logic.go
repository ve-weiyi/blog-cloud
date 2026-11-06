package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/model"
)

type BatchReviewCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchReviewCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchReviewCommentsLogic {
	return &BatchReviewCommentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量审核评论
//
// 逐条判断是否真的发生了状态迁移：只有迁移的那些才发通知。
func (l *BatchReviewCommentsLogic) BatchReviewComments(in *discussionrpc.BatchReviewCommentsRequest) (*discussionrpc.BatchReviewCommentsResponse, error) {
	if len(in.Ids) == 0 {
		return &discussionrpc.BatchReviewCommentsResponse{}, nil
	}
	if in.IsReview != enums.CommentIsReviewNo && in.IsReview != enums.CommentIsReviewYes {
		return nil, bizerr.NewBizError(bizcode.CodeParamValueInvalid, "审核状态不合法")
	}

	var affected int64
	approved := in.IsReview == enums.CommentIsReviewYes

	err := l.svcCtx.GormDB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		// 锁住这批评论再读状态：并发的另一个审核请求会在这里排队，拿到锁时看到的已是
		// 提交后的结果，它算出的待迁移集合为空，不会重复发通知。
		// 先读后写之间没有锁的话，「影响行数为 0」只能说明没生效，却拦不住已经发出的通知
		var comments []*model.TComment
		if err := tx.Table(l.svcCtx.TCommentModel.TableName()).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id IN ?", in.Ids).
			Find(&comments).Error; err != nil {
			return err
		}

		// 只保留真的会发生状态迁移的评论：重复审核是空操作
		changed := make([]*model.TComment, 0, len(comments))
		changedIDs := make([]int64, 0, len(comments))
		for _, comment := range comments {
			if comment.IsReview == in.IsReview {
				continue
			}
			changed = append(changed, comment)
			changedIDs = append(changedIDs, comment.Id)
		}
		if len(changedIDs) == 0 {
			return nil
		}

		// WHERE 仍带上原状态：锁之外还可能有别处改写状态
		rows, err := l.svcCtx.TCommentModel.WithTx(tx).UpdateFields(l.ctx,
			map[string]interface{}{"is_review": in.IsReview},
			"id IN ? AND is_review <> ?", changedIDs, in.IsReview)
		if err != nil {
			return err
		}
		affected = rows

		for _, comment := range changed {
			if err := emitReviewNotifications(l.ctx, l.svcCtx, tx, comment, approved); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &discussionrpc.BatchReviewCommentsResponse{SuccessCount: affected}, nil
}
