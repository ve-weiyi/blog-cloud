package friend

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/friendrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteFriendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除友链
func NewBatchDeleteFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteFriendLogic {
	return &BatchDeleteFriendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteFriendLogic) BatchDeleteFriend(req *types.IdsReq) (resp *types.BatchResp, err error) {
	in := &friendrpc.IdsReq{
		Ids: req.Ids,
	}

	out, err := l.svcCtx.FriendRpc.DeleteFriend(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
