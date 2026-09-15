package friend

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/siteservice"
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

func (l *BatchDeleteFriendLogic) BatchDeleteFriend(req *types.DeleteFriendReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.SiteService.BatchDeleteFriends(l.ctx, &siteservice.BatchDeleteFriendsRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
