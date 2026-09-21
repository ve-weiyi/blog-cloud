package syslogservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type CreateUploadLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUploadLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUploadLogLogic {
	return &CreateUploadLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建文件上传日志
func (l *CreateUploadLogLogic) CreateUploadLog(in *syslogrpc.CreateUploadLogRequest) (*syslogrpc.CreateUploadLogResponse, error) {
	data := convertUploadLogIn(in)
	_, err := l.svcCtx.TUploadLogModel.Insert(l.ctx, data)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.CreateUploadLogResponse{
		LogId: data.Id,
	}, nil
}
