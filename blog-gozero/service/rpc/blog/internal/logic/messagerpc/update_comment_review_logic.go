package messagerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCommentReviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCommentReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCommentReviewLogic {
	return &UpdateCommentReviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新评论审核状态
func (l *UpdateCommentReviewLogic) UpdateCommentReview(in *messagerpc.UpdateCommentReviewReq) (*messagerpc.BatchResp, error) {
	rows, err := l.svcCtx.TCommentModel.Updates(l.ctx, map[string]interface{}{
		"is_review": in.IsReview,
	}, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &messagerpc.BatchResp{
		SuccessCount: rows,
	}, nil
}
