package stats

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/analyticsservice"
)

type GetStatsDashboardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetStatsDashboardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStatsDashboardLogic {
	return &GetStatsDashboardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetStatsDashboardLogic) GetStatsDashboard(req *types.EmptyReq) (resp *types.GetStatsDashboardResp, err error) {
	out, err := l.svcCtx.AnalyticsService.GetDashboardStats(l.ctx, &analyticsservice.GetDashboardStatsRequest{})
	if err != nil {
		return nil, err
	}

	return &types.GetStatsDashboardResp{
		UserCount:    out.UserCount,
		ArticleCount: out.ArticleCount,
		MessageCount: out.MessageCount,
		Today: &types.DashboardStats{
			Date:         out.Today.Date,
			NewUsers:     out.Today.NewUsers,
			TotalUsers:   out.Today.TotalUsers,
			ActiveUsers:  out.Today.ActiveUsers,
			UvCount:      out.Today.UvCount,
			PvCount:      out.Today.PvCount,
			TotalUvCount: out.Today.TotalUvCount,
			TotalPvCount: out.Today.TotalPvCount,
		},
		UvGrowthRate:   out.UvGrowthRate,
		PvGrowthRate:   out.PvGrowthRate,
		UserGrowthRate: out.UserGrowthRate,
	}, nil
}
