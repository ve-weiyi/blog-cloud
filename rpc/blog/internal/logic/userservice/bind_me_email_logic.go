package userservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/metax"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/userrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BindMeEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindMeEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindMeEmailLogic {
	return &BindMeEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 绑定当前用户邮箱（验证码验证已在 API 层完成）
func (l *BindMeEmailLogic) BindMeEmail(in *userrpc.BindMeEmailRequest) (*userrpc.BindMeEmailResponse, error) {
	// 从上下文获取用户ID
	userID, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// 更新邮箱
	_, err = l.svcCtx.TUserModel.UpdateFields(l.ctx, map[string]interface{}{
		"email": in.Email,
	}, "user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	return &userrpc.BindMeEmailResponse{
		Success: true,
	}, nil
}
