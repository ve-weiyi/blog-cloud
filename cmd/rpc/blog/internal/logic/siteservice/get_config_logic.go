package siteservicelogic

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/siterpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
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
	entity, err := l.svcCtx.TConfigModel.FindOneByConfigKey(l.ctx, in.ConfigKey)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 配置缺失时返回空对象，调用方按零值降级，而不是整个页面 500
		return &siterpc.GetConfigResponse{ConfigValue: "{}"}, nil
	}
	if err != nil {
		return nil, err
	}

	return &siterpc.GetConfigResponse{
		ConfigValue: entity.Config,
	}, nil
}
