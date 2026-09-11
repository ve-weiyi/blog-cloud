package userservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/metax"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/userrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BindMePhoneLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindMePhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindMePhoneLogic {
	return &BindMePhoneLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 绑定当前用户手机号（验证码验证已在 API 层完成）
func (l *BindMePhoneLogic) BindMePhone(in *userrpc.BindMePhoneRequest) (*userrpc.BindMePhoneResponse, error) {
	// 从上下文获取用户ID
	userID, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// 更新手机号
	_, err = l.svcCtx.TUserModel.UpdateFields(l.ctx, map[string]interface{}{
		"mobile": in.Mobile,
	}, "user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	return &userrpc.BindMePhoneResponse{
		Success: true,
	}, nil
}
