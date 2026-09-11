package userservicelogic

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/constants/enums"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/userrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

type ReactivateAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReactivateAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReactivateAccountLogic {
	return &ReactivateAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 恢复账号（冷静期内）
func (l *ReactivateAccountLogic) ReactivateAccount(in *userrpc.ReactivateAccountRequest) (*userrpc.ReactivateAccountResponse, error) {
	var user *model.TUser
	var err error
	var identifier string

	if in.Email != nil && *in.Email != "" {
		identifier = *in.Email
		user, err = l.svcCtx.TUserModel.FindOneByEmail(l.ctx, identifier)
		if err != nil {
			return nil, err
		}
	} else if in.Mobile != nil && *in.Mobile != "" {
		identifier = *in.Mobile
		user, err = l.svcCtx.TUserModel.FindOneByMobile(l.ctx, identifier)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("email or phone is required")
	}

	if user.Status != enums.UserStatusLogoff {
		return nil, fmt.Errorf("account is not deactivated")
	}

	user.Status = enums.UserStatusNormal
	user.DeletedAt = nil
	_, err = l.svcCtx.TUserModel.Save(l.ctx, user)
	if err != nil {
		return nil, err
	}

	return &userrpc.ReactivateAccountResponse{Success: true}, nil
}
