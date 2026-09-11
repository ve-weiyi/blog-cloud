package socialservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/pb/socialrpc"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type UpdateFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFriendLogic {
	return &UpdateFriendLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateFriendLogic) UpdateFriend(in *socialrpc.UpdateFriendRequest) (*socialrpc.UpdateFriendResponse, error) {
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
	return &socialrpc.UpdateFriendResponse{Success: true}, nil
}
