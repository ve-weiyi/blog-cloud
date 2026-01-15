package friend

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/socialrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFriendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建友链
func NewAddFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFriendLogic {
	return &AddFriendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddFriendLogic) AddFriend(req *types.NewFriendReq) (resp *types.FriendBackVO, err error) {
	in := &socialrpc.AddFriendReq{
		Id:          req.Id,
		LinkName:    req.LinkName,
		LinkAvatar:  req.LinkAvatar,
		LinkAddress: req.LinkAddress,
		LinkIntro:   req.LinkIntro,
	}

	out, err := l.svcCtx.SocialRpc.AddFriend(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return convertFriendTypes(out.Friend), nil
}

func convertFriendTypes(out *socialrpc.Friend) *types.FriendBackVO {
	return &types.FriendBackVO{
		Id:          out.Id,
		LinkName:    out.LinkName,
		LinkAvatar:  out.LinkAvatar,
		LinkAddress: out.LinkAddress,
		LinkIntro:   out.LinkIntro,
		CreatedAt:   out.CreatedAt,
		UpdatedAt:   out.UpdatedAt,
	}
}
