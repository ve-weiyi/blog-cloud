package config

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/siteservice"
)

type GetWebsiteConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取网站配置
func NewGetWebsiteConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWebsiteConfigLogic {
	return &GetWebsiteConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWebsiteConfigLogic) GetWebsiteConfig(req *types.EmptyReq) (resp *types.WebsiteConfigVO, err error) {
	out, err := l.svcCtx.SiteService.GetConfig(l.ctx, &siteservice.GetConfigRequest{
		ConfigKey: "website_config",
	})
	if err != nil {
		return nil, err
	}

	resp = &types.WebsiteConfigVO{}
	if out.ConfigValue != "" {
		if err = json.Unmarshal([]byte(out.ConfigValue), resp); err != nil {
			return nil, err
		}
	}
	return
}
