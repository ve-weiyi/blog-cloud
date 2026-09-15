package home

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/siteservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/statsservice"
	"github.com/ve-weiyi/vkit/x/jsonconv"
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
	dashboard, err := l.svcCtx.StatsService.GetDashboardStats(l.ctx, &statsservice.GetDashboardStatsRequest{})
	if err != nil {
		return nil, err
	}

	articleStats, err := l.svcCtx.StatsService.GetArticleStats(l.ctx, &statsservice.GetArticleStatsRequest{})
	if err != nil {
		return nil, err
	}

	pages, err := l.svcCtx.SiteService.ListPages(l.ctx, &siteservice.ListPagesRequest{})
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

	conf, err := l.svcCtx.SiteService.GetConfig(l.ctx, &siteservice.GetConfigRequest{
		ConfigKey: "website_config",
	})
	if err != nil {
		return nil, err
	}

	config := types.WebsiteConfigVO{}
	err = jsonconv.JsonToAny(conf.ConfigValue, &config)
	if err != nil {
		return nil, err
	}

	// 查询已发布的全局通知
	noticeResp, err := l.svcCtx.NotificationService.ListNotifyMessages(l.ctx, &notificationservice.ListNotifyMessagesRequest{
		ListQuery: &notificationservice.ListQuery{Page: 1, PageSize: 5},
		Status:    stringPtr("published"),
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("查询通知列表失败: %v", err)
	}

	notices := make([]*types.NoticeVO, 0)
	if noticeResp != nil {
		for _, v := range noticeResp.List {
			notices = append(notices, &types.NoticeVO{
				Id:          v.Id,
				Title:       v.Title,
				Content:     v.Content,
				Type:        v.Category,
				Level:       v.Level,
				PublishTime: v.PublishedAt,
			})
		}
	}

	resp = &types.GetBlogHomeInfoResp{
		ArticleCount:       dashboard.ArticleCount,
		CategoryCount:      articleStats.CategoryCount,
		TagCount:           articleStats.TagCount,
		TotalUserViewCount: dashboard.Today.TotalUvCount,
		TotalPageViewCount: dashboard.Today.TotalPvCount,
		WebsiteConfig:      config,
		PageList:           ps,
		NoticeList:         notices,
	}

	return resp, nil
}

func stringPtr(s string) *string {
	return &s
}
