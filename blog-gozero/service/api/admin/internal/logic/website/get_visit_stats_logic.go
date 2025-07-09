package website

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/websiterpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVisitStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取访客数据分析
func NewGetVisitStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVisitStatsLogic {
	return &GetVisitStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetVisitStatsLogic) GetVisitStats(req *types.EmptyReq) (resp *types.GetVisitStatsResp, err error) {
	// 查询用户浏览量
	visit, err := l.svcCtx.WebsiteRpc.AnalysisVisit(l.ctx, &websiterpc.AnalysisVisitReq{})
	if err != nil {
		return nil, err
	}

	return &types.GetVisitStatsResp{
		TodayUvCount: visit.TodayUvCount,
		TotalUvCount: visit.TotalUvCount,
		UvGrowthRate: visit.UvGrowthRate,
		TodayPvCount: visit.TodayPvCount,
		TotalPvCount: visit.TotalPvCount,
		PvGrowthRate: visit.PvGrowthRate,
	}, nil
}
