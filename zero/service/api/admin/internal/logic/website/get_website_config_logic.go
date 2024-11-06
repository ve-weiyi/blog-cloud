package website

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/configrpc"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *GetWebsiteConfigLogic) GetWebsiteConfig(req *types.EmptyReq) (resp *types.WebsiteConfig, err error) {
	in := &configrpc.FindConfigReq{
		ConfigKey: "website_config",
	}

	out, err := l.svcCtx.ConfigRpc.FindConfig(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.WebsiteConfig{}
	jsonconv.JsonToAny(out.ConfigValue, &resp)
	return resp, nil
}
