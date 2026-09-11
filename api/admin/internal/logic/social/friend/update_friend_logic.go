package friend

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/socialservice"
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

func (l *UpdateFriendLogic) UpdateFriend(req *types.UpdateFriendReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.SocialService.UpdateFriend(l.ctx, &socialservice.UpdateFriendRequest{
		Id:          req.Id,
		LinkName:    req.LinkName,
		LinkAvatar:  req.LinkAvatar,
		LinkAddress: req.LinkAddress,
		LinkIntro:   req.LinkIntro,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
