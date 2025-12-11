package websiterpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/websiterpc"

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
func (l *UpdateFriendLogic) UpdateFriend(in *websiterpc.FriendNewReq) (*websiterpc.FriendDetailsResp, error) {
	entity := convertFriendIn(in)

	_, err := l.svcCtx.TFriendModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertFriendOut(entity), nil
}
