package comment

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/messagerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCommentReviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新评论审核状态
func NewUpdateCommentReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCommentReviewLogic {
	return &UpdateCommentReviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCommentReviewLogic) UpdateCommentReview(req *types.CommentReviewReq) (resp *types.BatchResp, err error) {
	in := &messagerpc.UpdateCommentReviewReq{
		Ids:      req.Ids,
		IsReview: req.IsReview,
	}

	out, err := l.svcCtx.MessageRpc.UpdateCommentReview(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
