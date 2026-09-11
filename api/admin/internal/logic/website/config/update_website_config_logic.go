package config

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/configservice"
)

type UpdateWebsiteConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新网站配置
func NewUpdateWebsiteConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWebsiteConfigLogic {
	return &UpdateWebsiteConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateWebsiteConfigLogic) UpdateWebsiteConfig(req *types.WebsiteConfigVO) (resp *types.EmptyResp, err error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.ConfigService.SaveConfig(l.ctx, &configservice.SaveConfigRequest{
		ConfigKey:   "website_config",
		ConfigValue: string(data),
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
