package auth

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/constants/enums"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/authservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 邮箱注册
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// 验证邮箱验证码（API 层负责）
	_, err = l.svcCtx.NotificationService.VerifyEmailCode(l.ctx, &notificationservice.VerifyEmailCodeRequest{
		Email: req.Email,
		Scene: enums.CodeSceneRegister,
		Code:  req.Code,
		BizId: "", // 空值，使用默认规则 scene:email
	})
	if err != nil {
		return nil, err
	}

	// 调用 RPC 创建用户
	in := &authservice.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Username: &req.Username,
		Nickname: &req.Nickname,
	}

	_, err = l.svcCtx.AuthService.Register(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return &types.RegisterResp{}, nil
}
