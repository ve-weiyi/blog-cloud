package socialrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/socialrpc"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFriendLogic {
	return &UpdateFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新友链
func (l *UpdateFriendLogic) UpdateFriend(in *socialrpc.UpdateFriendReq) (*socialrpc.UpdateFriendResp, error) {
	entity, err := l.svcCtx.TFriendModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	entity.LinkName = in.LinkName
	entity.LinkAvatar = in.LinkAvatar
	entity.LinkAddress = in.LinkAddress
	entity.LinkIntro = in.LinkIntro

	_, err = l.svcCtx.TFriendModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &socialrpc.UpdateFriendResp{
		Friend: convertFriendOut(entity),
	}, nil
}
