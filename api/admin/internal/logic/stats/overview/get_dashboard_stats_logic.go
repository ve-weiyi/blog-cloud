package overview

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/statsservice"
)

type GetDashboardStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDashboardStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDashboardStatsLogic {
	return &GetDashboardStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDashboardStatsLogic) GetDashboardStats(req *types.EmptyReq) (resp *types.GetStatsDashboardResp, err error) {
	out, err := l.svcCtx.StatsService.GetDashboardStats(l.ctx, &statsservice.GetDashboardStatsRequest{})
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
