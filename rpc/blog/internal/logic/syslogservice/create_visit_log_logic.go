package syslogservicelogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type CreateVisitLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateVisitLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateVisitLogLogic {
	return &CreateVisitLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建访问日志
func (l *CreateVisitLogLogic) CreateVisitLog(in *syslogrpc.CreateVisitLogRequest) (*syslogrpc.CreateVisitLogResponse, error) {
	data := convertVisitLogIn(in)
	_, err := l.svcCtx.TVisitLogModel.Insert(l.ctx, data)
	if err != nil {
		return nil, err
	}

	day := time.Now().Format("2006-01-02")
	visitor := in.DeviceId

	if err := l.addUserVisit(day, visitor); err != nil {
		l.Errorf("addUserVisit error: %v", err)
	}
	if err := l.addPageView(day); err != nil {
		l.Errorf("addPageView error: %v", err)
	}

	return &syslogrpc.CreateVisitLogResponse{
		LogId: data.Id,
	}, nil
}

// getPrevCumulativeTotals 获取指定日期之前（不含当天）的累计 UV/PV
func (l *CreateVisitLogLogic) getPrevCumulativeTotals(day string) (prevTotalUv int64, prevTotalPv int64) {
	row := l.svcCtx.GormDB.Raw(
		"SELECT COALESCE(SUM(uv_count), 0), COALESCE(SUM(pv_count), 0) FROM t_daily_stats WHERE date < ?",
		day,
	).Row()
	row.Scan(&prevTotalUv, &prevTotalPv)
	return
}

// 添加日访客数
func (l *CreateVisitLogLogic) addUserVisit(day string, visitor string) error {
	key := cachekey.GetDailyUserVisitKey(day)

	// 用 SAdd 返回值判断是否为当天首次访问，避免 SIsMember+SAdd 的竞态
	added, err := l.svcCtx.Redis.SAdd(l.ctx, key, visitor).Result()
	if err != nil {
		return err
	}

	if added == 0 {
		return nil
	}

	// 首次添加时设置 30 天过期
	if size, _ := l.svcCtx.Redis.SCard(l.ctx, key).Result(); size == 1 {
		l.svcCtx.Redis.Expire(l.ctx, key, 30*24*time.Hour)
	}

	prevUv, prevPv := l.getPrevCumulativeTotals(day)

	// 新行: total_uv_count 继承前一天累计值+1, total_pv_count 继承前一天累计值
	result := l.svcCtx.GormDB.Exec(
		"INSERT INTO t_daily_stats (date, uv_count, total_uv_count, pv_count, total_pv_count) VALUES (?, 1, ?, 0, ?) ON DUPLICATE KEY UPDATE uv_count = uv_count + 1, total_uv_count = total_uv_count + 1",
		day, prevUv+1, prevPv,
	)
	if result.Error != nil {
		return result.Error
	}

	if err := l.svcCtx.Redis.Incr(l.ctx, cachekey.TotalUserViewCountKey).Err(); err != nil {
		l.Errorf("addUserVisit Incr total error: %v", err)
	}
	if err := l.svcCtx.Redis.ZIncrBy(l.ctx, cachekey.DailyUserViewCountKey, 1, day).Err(); err != nil {
		l.Errorf("addUserVisit ZIncrBy daily error: %v", err)
	}

	return nil
}

// 添加页面浏览量
func (l *CreateVisitLogLogic) addPageView(day string) error {
	prevUv, prevPv := l.getPrevCumulativeTotals(day)

	// 新行: total_uv_count 继承前一天累计值, total_pv_count 继承前一天累计值+1
	result := l.svcCtx.GormDB.Exec(
		"INSERT INTO t_daily_stats (date, uv_count, total_uv_count, pv_count, total_pv_count) VALUES (?, 0, ?, 1, ?) ON DUPLICATE KEY UPDATE pv_count = pv_count + 1, total_pv_count = total_pv_count + 1",
		day, prevUv, prevPv+1,
	)
	if result.Error != nil {
		return result.Error
	}

	if err := l.svcCtx.Redis.Incr(l.ctx, cachekey.TotalPageViewCountKey).Err(); err != nil {
		l.Errorf("addPageView Incr total error: %v", err)
	}
	if err := l.svcCtx.Redis.ZIncrBy(l.ctx, cachekey.DailyPageViewCountKey, 1, day).Err(); err != nil {
		l.Errorf("addPageView ZIncrBy daily error: %v", err)
	}

	return nil
}
