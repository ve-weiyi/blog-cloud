package syslogservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type UpdateLoginLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLoginLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLoginLogLogic {
	return &UpdateLoginLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新登出时间
func (l *UpdateLoginLogLogic) UpdateLoginLog(in *syslogrpc.UpdateLoginLogRequest) (*syslogrpc.UpdateLoginLogResponse, error) {
	logoutAt := timeFromMilli(in.LogoutAt)
	fields := map[string]interface{}{
		"logout_at": logoutAt,
	}

	_, err := l.svcCtx.TLoginLogModel.UpdateFields(l.ctx, fields, "id = ?", in.LogId)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.UpdateLoginLogResponse{
		Success: true,
	}, nil
}
