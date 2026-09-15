package overview

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/statsservice"
)

type GetVisitTrendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取访客数据趋势
func NewGetVisitTrendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVisitTrendLogic {
	return &GetVisitTrendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetVisitTrendLogic) GetVisitTrend(req *types.GetVisitTrendReq) (resp *types.GetVisitTrendResp, err error) {
	out, err := l.svcCtx.StatsService.GetStatsTrend(l.ctx, &statsservice.GetStatsTrendRequest{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.GetVisitTrendResp{
		List: make([]*types.DashboardStats, 0, len(out.List)),
	}
	for _, d := range out.List {
		resp.List = append(resp.List, &types.DashboardStats{
			Date:         d.Date,
			NewUsers:     d.NewUsers,
			TotalUsers:   d.TotalUsers,
			ActiveUsers:  d.ActiveUsers,
			UvCount:      d.UvCount,
			PvCount:      d.PvCount,
			TotalUvCount: d.TotalUvCount,
			TotalPvCount: d.TotalPvCount,
		})
	}
	return
}
