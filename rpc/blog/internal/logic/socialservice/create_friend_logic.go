package socialservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/pb/socialrpc"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

type CreateFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFriendLogic {
	return &CreateFriendLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *CreateFriendLogic) CreateFriend(in *socialrpc.CreateFriendRequest) (*socialrpc.CreateFriendResponse, error) {
	entity := &model.TFriend{
		LinkName:    in.LinkName,
		LinkAvatar:  in.LinkAvatar,
		LinkAddress: in.LinkAddress,
		LinkIntro:   in.LinkIntro,
	}
	_, err := l.svcCtx.TFriendModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}
	return &socialrpc.CreateFriendResponse{Id: entity.Id}, nil
}
