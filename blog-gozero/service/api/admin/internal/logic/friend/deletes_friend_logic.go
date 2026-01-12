package friend

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/socialrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesFriendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除友链
func NewDeletesFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesFriendLogic {
	return &DeletesFriendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletesFriendLogic) DeletesFriend(req *types.IdsReq) (resp *types.BatchResp, err error) {
	in := &socialrpc.DeletesFriendReq{
		Ids: req.Ids,
	}

	out, err := l.svcCtx.SocialRpc.DeletesFriend(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
