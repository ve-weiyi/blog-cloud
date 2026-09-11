package syslogservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type CreateOperationLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOperationLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOperationLogLogic {
	return &CreateOperationLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建操作日志
func (l *CreateOperationLogLogic) CreateOperationLog(in *syslogrpc.CreateOperationLogRequest) (*syslogrpc.CreateOperationLogResponse, error) {
	data := convertOperationLogIn(in)
	_, err := l.svcCtx.TOperationLogModel.Insert(l.ctx, data)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.CreateOperationLogResponse{
		LogId: data.Id,
	}, nil
}
