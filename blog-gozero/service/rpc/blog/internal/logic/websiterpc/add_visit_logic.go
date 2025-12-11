package websiterpclogic

import (
	"context"
	"time"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/global/constant"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rediskey"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rpcutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/websiterpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddVisitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddVisitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddVisitLogic {
	return &AddVisitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 添加用户访问记录
func (l *AddVisitLogic) AddVisit(in *websiterpc.AddVisitReq) (*websiterpc.AddVisitResp, error) {
	visitor, err := rpcutils.GetTerminalIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	day := time.Now().Format("2006-01-02")
	_, err = l.addUserVisit(day, visitor)
	if err != nil {
		return nil, err
	}

	_, err = l.addPageView(day)
	if err != nil {
		return nil, err
	}

	return &websiterpc.AddVisitResp{}, nil
}

// 添加日访客数
func (l *AddVisitLogic) addUserVisit(day string, visitor string) (any, error) {

	key := rediskey.GetDailyUserVisitKey(day)
	ok, err := l.svcCtx.Redis.SIsMember(l.ctx, key, visitor).Result()
	if err != nil {
		return nil, err
	}

	if !ok {
		// 添加当天的用户访问量+1 mysql
		result := l.svcCtx.Gorm.Exec("UPDATE t_visit_daily_stats SET view_count = view_count + 1 WHERE date = ? and visit_type = ?", day, constant.VisitTypeUv)
		if result.RowsAffected == 0 {
			_, err := l.svcCtx.TVisitDailyStatsModel.Insert(l.ctx, &model.TVisitDailyStats{
				Id:        0,
				Date:      day,
				ViewCount: 1,
				VisitType: constant.VisitTypeUv,
			})
			if err != nil {
				return nil, err
			}
		}

		// 添加总访问量
		totalKey := rediskey.GetTotalUserViewCountKey()
		_, err = l.svcCtx.Redis.Incr(l.ctx, totalKey).Result()
		if err != nil {
			return nil, err
		}

		// 添加日访问量
		dailyKey := rediskey.GetDailyUserViewCountKey()
		_, err = l.svcCtx.Redis.ZIncrBy(l.ctx, dailyKey, 1, day).Result()
		if err != nil {
			return nil, err
		}

		// 保存用户访问标识
		_, err = l.svcCtx.Redis.SAdd(l.ctx, key, visitor).Result()
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// 添加页面浏览量
func (l *AddVisitLogic) addPageView(day string) (any, error) {
	// 添加当天的用户访问量+1
	result := l.svcCtx.Gorm.Exec("UPDATE t_visit_daily_stats SET view_count = view_count + 1 WHERE date = ? and visit_type = ?", day, constant.VisitTypePv)
	if result.RowsAffected == 0 {
		_, err := l.svcCtx.TVisitDailyStatsModel.Insert(l.ctx, &model.TVisitDailyStats{
			Id:        0,
			Date:      day,
			ViewCount: 1,
			VisitType: constant.VisitTypePv,
		})
		if err != nil {
			return nil, err
		}
	}

	// 添加总访问量
	totalKey := rediskey.GetTotalPageViewCountKey()
	_, err := l.svcCtx.Redis.Incr(l.ctx, totalKey).Result()
	if err != nil {
		return nil, err
	}

	// 添加日访问量
	dailyKey := rediskey.GetDailyPageViewCountKey()
	_, err = l.svcCtx.Redis.ZIncrBy(l.ctx, dailyKey, 1, day).Result()
	if err != nil {
		return nil, err
	}
	return nil, nil
}
