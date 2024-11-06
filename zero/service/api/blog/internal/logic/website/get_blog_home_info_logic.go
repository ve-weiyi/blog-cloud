package website

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/configrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBlogHomeInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取博客前台首页信息
func NewGetBlogHomeInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBlogHomeInfoLogic {
	return &GetBlogHomeInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBlogHomeInfoLogic) GetBlogHomeInfo(req *types.GetBlogHomeInfoReq) (resp *types.GetBlogHomeInfoResp, err error) {
	analysis, err := l.svcCtx.ArticleRpc.AnalysisArticle(l.ctx, &articlerpc.EmptyReq{})
	if err != nil {
		return nil, err
	}

	in := &configrpc.FindConfigReq{
		ConfigKey: "website_config",
	}

	out, err := l.svcCtx.ConfigRpc.FindConfig(l.ctx, in)
	if err != nil {
		return nil, err
	}

	config := &types.WebsiteConfigDTO{}
	jsonconv.JsonToAny(out.ConfigValue, &config)

	resp = &types.GetBlogHomeInfoResp{
		ArticleCount:  analysis.ArticleCount,
		CategoryCount: analysis.CategoryCount,
		TagCount:      analysis.TagCount,
		ViewsCount:    0,
		WebsiteConfig: *config,
		PageList:      make([]*types.PageDTO, 0),
	}

	return resp, nil
}
