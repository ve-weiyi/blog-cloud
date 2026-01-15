package socialrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/socialrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFriendLogic {
	return &AddFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建友链
func (l *AddFriendLogic) AddFriend(in *socialrpc.AddFriendReq) (*socialrpc.AddFriendResp, error) {
	entity := &model.TFriend{
		Id:          in.Id,
		LinkName:    in.LinkName,
		LinkAvatar:  in.LinkAvatar,
		LinkAddress: in.LinkAddress,
		LinkIntro:   in.LinkIntro,
		//CreatedAt:   time.Time{},
		//UpdatedAt:   time.Time{},
	}

	_, err := l.svcCtx.TFriendModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &socialrpc.AddFriendResp{
		Friend: convertFriendOut(entity),
	}, nil
}

func convertFriendOut(in *model.TFriend) (out *socialrpc.Friend) {
	out = &socialrpc.Friend{
		Id:          in.Id,
		LinkName:    in.LinkName,
		LinkAvatar:  in.LinkAvatar,
		LinkAddress: in.LinkAddress,
		LinkIntro:   in.LinkIntro,
		CreatedAt:   in.CreatedAt.UnixMilli(),
		UpdatedAt:   in.UpdatedAt.UnixMilli(),
	}

	return out
}
