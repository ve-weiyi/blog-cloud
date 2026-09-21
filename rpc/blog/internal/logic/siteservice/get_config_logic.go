package siteservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/siterpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type GetConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConfigLogic {
	return &GetConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询配置
func (l *GetConfigLogic) GetConfig(in *siterpc.GetConfigRequest) (*siterpc.GetConfigResponse, error) {
	entity, err := l.svcCtx.TConfigModel.FindOneByKey(l.ctx, in.ConfigKey)
	if err != nil {
		return nil, err
	}

	return &siterpc.GetConfigResponse{
		ConfigValue: entity.Config,
	}, nil
}
