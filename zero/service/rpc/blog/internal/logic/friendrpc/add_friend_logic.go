package friendrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/friendrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

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
func (l *AddFriendLogic) AddFriend(in *friendrpc.FriendNewReq) (*friendrpc.FriendDetails, error) {
	entity := convertFriendIn(in)

	_, err := l.svcCtx.TFriendModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertFriendOut(entity), nil
}

func convertFriendIn(in *friendrpc.FriendNewReq) (out *model.TFriend) {
	out = &model.TFriend{
		Id:          in.Id,
		LinkName:    in.LinkName,
		LinkAvatar:  in.LinkAvatar,
		LinkAddress: in.LinkAddress,
		LinkIntro:   in.LinkIntro,
		//CreatedAt:   time.Time{},
		//UpdatedAt:   time.Time{},
	}

	return out
}

func convertFriendOut(in *model.TFriend) (out *friendrpc.FriendDetails) {
	out = &friendrpc.FriendDetails{
		Id:          in.Id,
		LinkName:    in.LinkName,
		LinkAvatar:  in.LinkAvatar,
		LinkAddress: in.LinkAddress,
		LinkIntro:   in.LinkIntro,
		CreatedAt:   in.CreatedAt.Unix(),
		UpdatedAt:   in.UpdatedAt.Unix(),
	}

	return out
}
