package statsservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/statsrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type GetStatsTrendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStatsTrendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStatsTrendLogic {
	return &GetStatsTrendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询统计趋势（日期范围内的每日数据）
func (l *GetStatsTrendLogic) GetStatsTrend(in *statsrpc.GetStatsTrendRequest) (*statsrpc.GetStatsTrendResponse, error) {
	records, err := l.svcCtx.TDailyStatsModel.FindALL(l.ctx,
		"`date` >= ? AND `date` <= ?",
		in.GetStartDate(), in.GetEndDate())
	if err != nil {
		return nil, err
	}

	list := make([]*statsrpc.DailyStats, 0, len(records))
	for _, r := range records {
		list = append(list, &statsrpc.DailyStats{
			Date:         r.Date,
			NewUsers:     r.NewUsers,
			TotalUsers:   r.TotalUsers,
			ActiveUsers:  r.ActiveUsers,
			UvCount:      r.UvCount,
			PvCount:      r.PvCount,
			TotalUvCount: r.TotalUvCount,
			TotalPvCount: r.TotalPvCount,
		})
	}

	return &statsrpc.GetStatsTrendResponse{List: list}, nil
}
