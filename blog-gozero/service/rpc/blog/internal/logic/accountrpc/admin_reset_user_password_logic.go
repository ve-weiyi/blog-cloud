package accountrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizcode"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizerr"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/cryptox"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminResetUserPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminResetUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminResetUserPasswordLogic {
	return &AdminResetUserPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 管理员重置用户密码
func (l *AdminResetUserPasswordLogic) AdminResetUserPassword(in *accountrpc.AdminResetUserPasswordReq) (*accountrpc.EmptyResp, error) {
	// 验证用户是否存在
	user, err := l.svcCtx.TUserModel.FindOneByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, bizerr.NewBizError(bizcode.CodeUserNotExist, err.Error())
	}

	// 更新密码
	user.Password = cryptox.BcryptHash(in.Password)

	_, err = l.svcCtx.TUserModel.Save(l.ctx, user)
	if err != nil {
		return nil, err
	}

	return &accountrpc.EmptyResp{}, nil
}
