package websiterpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/global/constant"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/websiterpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindVisitTrendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindVisitTrendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindVisitTrendLogic {
	return &FindVisitTrendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询用户访问趋势
func (l *FindVisitTrendLogic) FindVisitTrend(in *websiterpc.FindVisitTrendReq) (*websiterpc.FindVisitTrendResp, error) {
	uvs, err := l.svcCtx.TVisitDailyStatsModel.FindALL(l.ctx, "date >= ? and date <= ? and visit_type = ?", in.StartDate, in.EndDate, constant.VisitTypeUv)
	if err != nil {
		return nil, err
	}

	pvs, err := l.svcCtx.TVisitDailyStatsModel.FindALL(l.ctx, "date >= ? and date <= ? and visit_type = ?", in.StartDate, in.EndDate, constant.VisitTypePv)
	if err != nil {
		return nil, err
	}

	uvTrend := make([]*websiterpc.VisitDailyStatistics, 0)
	pvTrend := make([]*websiterpc.VisitDailyStatistics, 0)
	for _, uv := range uvs {
		uvTrend = append(uvTrend, &websiterpc.VisitDailyStatistics{
			Date:  uv.Date,
			Count: uv.ViewCount,
		})
	}

	for _, pv := range pvs {
		pvTrend = append(pvTrend, &websiterpc.VisitDailyStatistics{
			Date:  pv.Date,
			Count: pv.ViewCount,
		})
	}

	return &websiterpc.FindVisitTrendResp{
		UvTrend: uvTrend,
		PvTrend: pvTrend,
	}, nil
}
