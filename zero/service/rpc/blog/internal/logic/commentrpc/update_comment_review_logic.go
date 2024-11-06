package commentrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/commentrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

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
func (l *UpdateCommentReviewLogic) UpdateCommentReview(in *commentrpc.UpdateCommentReviewReq) (*commentrpc.BatchResp, error) {
	rows, err := l.svcCtx.TCommentModel.Updates(l.ctx, map[string]interface{}{
		"is_review": in.IsReview,
	}, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &commentrpc.BatchResp{
		SuccessCount: rows,
	}, nil
}
