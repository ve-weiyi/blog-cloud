package comment

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/discussionservice"
)

type BatchDeleteCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除评论
func NewBatchDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteCommentLogic {
	return &BatchDeleteCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteCommentLogic) BatchDeleteComment(req *types.DeleteCommentReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.DiscussionService.BatchDeleteComments(l.ctx, &discussionservice.BatchDeleteCommentsRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
