package friend

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/socialrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFriendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新友链
func NewUpdateFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFriendLogic {
	return &UpdateFriendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateFriendLogic) UpdateFriend(req *types.NewFriendReq) (resp *types.FriendBackVO, err error) {
	in := &socialrpc.UpdateFriendReq{
		Id:          req.Id,
		LinkName:    req.LinkName,
		LinkAvatar:  req.LinkAvatar,
		LinkAddress: req.LinkAddress,
		LinkIntro:   req.LinkIntro,
	}

	out, err := l.svcCtx.SocialRpc.UpdateFriend(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return convertFriendTypes(out.Friend), nil
}
