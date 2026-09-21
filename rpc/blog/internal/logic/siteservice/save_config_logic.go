package siteservicelogic

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/siterpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

type SaveConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveConfigLogic {
	return &SaveConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 保存配置
func (l *SaveConfigLogic) SaveConfig(in *siterpc.SaveConfigRequest) (*siterpc.SaveConfigResponse, error) {
	entity := &model.TConfig{
		Key:    in.ConfigKey,
		Config: in.ConfigValue,
	}

	// 已存在则带主键走更新，不存在则主键留零走新增。
	// 必须区分"记录不存在"与真实查询错误：后者若被当作不存在，
	// 会让 Save 走新增，撞上 t_config 的 uk_key 唯一键，
	// 把真正的失败原因掩盖成一个误导性的重复键错误。
	result, err := l.svcCtx.TConfigModel.FindOneByKey(l.ctx, in.ConfigKey)
	switch {
	case err == nil:
		entity.Id = result.Id
	case errors.Is(err, gorm.ErrRecordNotFound):
		// 新增，Id 保持零值
	default:
		return nil, err
	}

	if _, err = l.svcCtx.TConfigModel.Save(l.ctx, entity); err != nil {
		return nil, err
	}

	return &siterpc.SaveConfigResponse{
		Success: true,
	}, nil
}
