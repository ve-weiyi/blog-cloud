package userservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/userrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type ResetUserPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetUserPasswordLogic {
	return &ResetUserPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 重置用户密码
func (l *ResetUserPasswordLogic) ResetUserPassword(in *userrpc.ResetUserPasswordRequest) (*userrpc.ResetUserPasswordResponse, error) {
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 更新密码
	_, err = l.svcCtx.TUserModel.UpdateFields(l.ctx, map[string]interface{}{
		"password": string(hashedPassword),
	}, "user_id = ?", in.UserId)
	if err != nil {
		return nil, err
	}

	return &userrpc.ResetUserPasswordResponse{
		Success: true,
	}, nil
}
