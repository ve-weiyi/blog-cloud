package analyticsservicelogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/analyticsrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type GetDashboardStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDashboardStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDashboardStatsLogic {
	return &GetDashboardStatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDashboardStatsLogic) GetDashboardStats(in *analyticsrpc.GetDashboardStatsRequest) (*analyticsrpc.GetDashboardStatsResponse, error) {
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	deps := realtimeStatsDeps{
		ctx:             l.ctx,
		redis:           l.svcCtx.Redis,
		gormDB:          l.svcCtx.GormDB,
		dailyStatsModel: l.svcCtx.TDailyStatsModel,
		userModel:       l.svcCtx.TUserModel,
	}

	userCount, _ := l.svcCtx.TUserModel.FindCount(l.ctx, "")
	articleCount, _ := l.svcCtx.TArticleModel.FindCount(l.ctx, "")
	messageCount, _ := l.svcCtx.TMessageModel.FindCount(l.ctx, "")

	todayStats := buildRealtimeDailyStats(deps, today)

	yesterdayRecord, err := l.svcCtx.TDailyStatsModel.FindOneByDate(l.ctx, yesterday)
	var uvGrowthRate, pvGrowthRate, userGrowthRate float64
	if err == nil {
		if yesterdayRecord.UvCount != 0 {
			uvGrowthRate = float64(todayStats.UvCount-yesterdayRecord.UvCount) / float64(yesterdayRecord.UvCount)
		}
		if yesterdayRecord.PvCount != 0 {
			pvGrowthRate = float64(todayStats.PvCount-yesterdayRecord.PvCount) / float64(yesterdayRecord.PvCount)
		}
		if yesterdayRecord.NewUsers != 0 {
			userGrowthRate = float64(todayStats.NewUsers-yesterdayRecord.NewUsers) / float64(yesterdayRecord.NewUsers)
		}
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &analyticsrpc.GetDashboardStatsResponse{
		UserCount:      userCount,
		ArticleCount:   articleCount,
		MessageCount:   messageCount,
		Today:          todayStats,
		UvGrowthRate:   uvGrowthRate,
		PvGrowthRate:   pvGrowthRate,
		UserGrowthRate: userGrowthRate,
	}, nil
}
