package userservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/metax"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/userrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type GetMeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMeLogic {
	return &GetMeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取当前用户详细资料
func (l *GetMeLogic) GetMe(in *userrpc.GetMeRequest) (*userrpc.GetMeResponse, error) {
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

	return &userrpc.GetMeResponse{
		MeInfo: convertTUserToMeInfo(l.ctx, l.svcCtx, user),
	}, nil
}
