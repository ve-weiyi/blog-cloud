package comment

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/commentrpc"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *BatchDeleteCommentLogic) BatchDeleteComment(req *types.IdsReq) (resp *types.BatchResp, err error) {
	in := &commentrpc.IdsReq{
		Ids: req.Ids,
	}

	out, err := l.svcCtx.CommentRpc.DeleteComment(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
