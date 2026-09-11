package syslogservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type CreateLoginLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLoginLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLoginLogLogic {
	return &CreateLoginLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建用户登录日志
func (l *CreateLoginLogLogic) CreateLoginLog(in *syslogrpc.CreateLoginLogRequest) (*syslogrpc.CreateLoginLogResponse, error) {
	data := convertLoginLogIn(in)
	_, err := l.svcCtx.TLoginLogModel.Insert(l.ctx, data)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.CreateLoginLogResponse{
		LogId: data.Id,
	}, nil
}
