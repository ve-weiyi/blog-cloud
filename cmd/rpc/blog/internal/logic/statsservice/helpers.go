package statsservicelogic

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/statsrpc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/cachekey"
	"github.com/ve-weiyi/blog-cloud/model"
)

type realtimeStatsDeps struct {
	ctx             context.Context
	redis           *redis.Client
	gormDB          *gorm.DB
	dailyStatsModel model.TDailyStatsModel
	userModel       model.TUserModel
}

func buildRealtimeDailyStats(deps realtimeStatsDeps, date string) *statsrpc.DailyStats {
	return &statsrpc.DailyStats{
		Date:         date,
		NewUsers:     getNewUsersCount(deps, date),
		TotalUsers:   getTotalUsersCount(deps),
		ActiveUsers:  getActiveUsersCount(deps, date),
		UvCount:      getDayCount(deps, cachekey.DailyUserViewCountKey, date),
		PvCount:      getDayCount(deps, cachekey.DailyPageViewCountKey, date),
		TotalUvCount: getTotalCount(deps, cachekey.TotalUserViewCountKey),
		TotalPvCount: getTotalCount(deps, cachekey.TotalPageViewCountKey),
	}
}

func getNewUsersCount(deps realtimeStatsDeps, date string) int64 {
	count, _ := deps.userModel.FindCount(deps.ctx, "DATE(created_at) = ?", date)
	return count
}

func getTotalUsersCount(deps realtimeStatsDeps) int64 {
	count, _ := deps.userModel.FindCount(deps.ctx, "")
	return count
}

func getActiveUsersCount(deps realtimeStatsDeps, date string) int64 {
	var count int64
	deps.gormDB.Table("t_visit_log").
		Where("DATE(created_at) = ?", date).
		Distinct("user_id").
		Count(&count)
	return count
}

// parseDay 把 YYYY-MM-DD 解析为当天零点；解析失败返回零值，查询自然查不到。
func parseDay(day string) time.Time {
	t, err := time.Parse("2006-01-02", day)
	if err != nil {
		return time.Time{}
	}
	return t
}

func getDayCount(deps realtimeStatsDeps, key string, day string) int64 {
	count, err := deps.redis.ZScore(deps.ctx, key, day).Result()
	if err != nil {
		record, err := deps.dailyStatsModel.FindOneByDate(deps.ctx, parseDay(day))
		if err != nil {
			return 0
		}
		if key == cachekey.DailyUserViewCountKey {
			count = float64(record.UvCount)
		} else {
			count = float64(record.PvCount)
		}
		_, _ = deps.redis.ZIncrBy(deps.ctx, key, count, day).Result()
		logx.Errorf("%s ZIncrBy fallback error: %v", key, err)
		return int64(count)
	}
	return int64(count)
}

func getTotalCount(deps realtimeStatsDeps, key string) int64 {
	count, err := deps.redis.Get(deps.ctx, key).Int64()
	if err != nil {
		column := "uv_count"
		if key != cachekey.TotalUserViewCountKey {
			column = "pv_count"
		}
		deps.gormDB.Raw(
			"SELECT COALESCE(SUM(" + column + "), 0) FROM t_daily_stats",
		).Scan(&count)
		_, _ = deps.redis.Set(deps.ctx, key, count, 0).Result()
		logx.Errorf("%s Set fallback error: %v", key, err)
		return count
	}
	return count
}
