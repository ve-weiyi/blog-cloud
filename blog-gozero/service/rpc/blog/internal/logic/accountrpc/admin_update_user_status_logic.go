package accountrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateUserStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminUpdateUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateUserStatusLogic {
	return &AdminUpdateUserStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改用户状态
func (l *AdminUpdateUserStatusLogic) AdminUpdateUserStatus(in *accountrpc.AdminUpdateUserStatusReq) (*accountrpc.EmptyResp, error) {
	user, err := l.svcCtx.TUserModel.FindOneByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	user.Status = in.Status

	_, err = l.svcCtx.TUserModel.Save(l.ctx, user)
	if err != nil {
		return nil, err
	}

	return &accountrpc.EmptyResp{}, nil
}
