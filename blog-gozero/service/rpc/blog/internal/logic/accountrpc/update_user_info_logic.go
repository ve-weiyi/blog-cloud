package accountrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rpcutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改用户信息
func (l *UpdateUserInfoLogic) UpdateUserInfo(in *accountrpc.UpdateUserInfoReq) (*accountrpc.UpdateUserInfoResp, error) {
	userId, err := rpcutils.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	user, err := l.svcCtx.TUserModel.FindOneByUserId(l.ctx, userId)
	if err != nil {
		return nil, err
	}

	user.Nickname = in.Nickname
	user.Info = in.Info

	_, err = l.svcCtx.TUserModel.Save(l.ctx, user)
	if err != nil {
		return nil, err
	}

	return &accountrpc.UpdateUserInfoResp{}, nil
}
