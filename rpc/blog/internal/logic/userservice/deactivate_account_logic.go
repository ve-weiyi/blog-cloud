package userservicelogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/constants/enums"
	"github.com/ve-weiyi/blog-cloud/infra/metax"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/userrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type DeactivateAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeactivateAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeactivateAccountLogic {
	return &DeactivateAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 停用账号（进入冷静期）
func (l *DeactivateAccountLogic) DeactivateAccount(in *userrpc.DeactivateAccountRequest) (*userrpc.DeactivateAccountResponse, error) {
	// 从上下文获取用户ID
	userID, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// 查询用户信息
	user, err := l.svcCtx.TUserModel.FindOneByUserId(l.ctx, userID)
	if err != nil {
		return nil, err
	}

	// 设置冷静期（30天）
	const coolingPeriodDays = 30
	now := time.Now()
	canReactivateUntil := now.AddDate(0, 0, coolingPeriodDays)

	// 更新用户状态
	user.Status = enums.UserStatusLogoff
	user.DeletedAt = &now
	_, err = l.svcCtx.TUserModel.Save(l.ctx, user)
	if err != nil {
		return nil, err
	}

	return &userrpc.DeactivateAccountResponse{
		Success:            true,
		CoolingPeriodDays:  coolingPeriodDays,
		CanReactivateUntil: canReactivateUntil.UnixMilli(),
	}, nil
}
