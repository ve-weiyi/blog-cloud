package website

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/global/constant"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/configrpc"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"

	"github.com/zeromicro/go-zero/core/logx"
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
	in := &configrpc.SaveConfigReq{
		ConfigKey:   constant.ConfigKeyWebsite,
		ConfigValue: jsonconv.AnyToJsonNE(req),
	}

	_, err = l.svcCtx.ConfigRpc.SaveConfig(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return
}
