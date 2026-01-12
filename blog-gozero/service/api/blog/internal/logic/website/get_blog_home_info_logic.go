package website

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/constant"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/configrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/resourcerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/websiterpc"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/jsonconv"

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
	_, err = l.svcCtx.WebsiteRpc.AddVisit(l.ctx, &websiterpc.AddVisitReq{})
	if err != nil {
		return nil, err
	}

	analysis, err := l.svcCtx.ArticleRpc.AnalysisArticle(l.ctx, &articlerpc.AnalysisArticleReq{})
	if err != nil {
		return nil, err
	}

	visit, err := l.svcCtx.WebsiteRpc.AnalysisVisit(l.ctx, &websiterpc.AnalysisVisitReq{})
	if err != nil {
		return nil, err
	}

	pages, err := l.svcCtx.ResourceRpc.FindPageList(l.ctx, &resourcerpc.FindPageListReq{})
	if err != nil {
		return nil, err
	}

	ps := make([]*types.PageVO, 0)
	for _, v := range pages.List {
		p := &types.PageVO{
			Id:         v.Id,
			PageName:   v.PageName,
			PageLabel:  v.PageLabel,
			PageCover:  v.PageCover,
			IsCarousel: v.IsCarousel,
		}
		ps = append(ps, p)
	}

	conf, err := l.svcCtx.ConfigRpc.FindConfig(l.ctx, &configrpc.FindConfigReq{
		ConfigKey: constant.ConfigKeyWebsite,
	})
	if err != nil {
		return nil, err
	}

	config := types.WebsiteConfigVO{}
	err = jsonconv.JsonToAny(conf.ConfigValue, &config)
	if err != nil {
		return nil, err
	}

	resp = &types.GetBlogHomeInfoResp{
		ArticleCount:       analysis.ArticleCount,
		CategoryCount:      analysis.CategoryCount,
		TagCount:           analysis.TagCount,
		TotalUserViewCount: visit.TotalUvCount,
		TotalPageViewCount: visit.TotalPvCount,
		WebsiteConfig:      config,
		PageList:           ps,
	}

	return resp, nil
}
