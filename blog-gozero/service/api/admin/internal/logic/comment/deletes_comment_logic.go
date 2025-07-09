package comment

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/messagerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除评论
func NewDeletesCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesCommentLogic {
	return &DeletesCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletesCommentLogic) DeletesComment(req *types.IdsReq) (resp *types.BatchResp, err error) {
	in := &messagerpc.IdsReq{
		Ids: req.Ids,
	}

	out, err := l.svcCtx.MessageRpc.DeletesComment(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
