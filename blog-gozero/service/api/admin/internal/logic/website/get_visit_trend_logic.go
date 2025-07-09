package website

import (
	"context"
	"sort"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/websiterpc"

	"github.com/zeromicro/go-zero/core/logx"
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
	// 查询用户浏览量
	visit, err := l.svcCtx.WebsiteRpc.FindVisitTrend(l.ctx, &websiterpc.FindVisitTrendReq{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	})
	if err != nil {
		return nil, err
	}

	// 使用map临时存储数据
	trendMap := make(map[string]types.VisitTrendVO)
	dates := make([]string, 0)

	// 处理 PV 趋势数据
	for _, v := range visit.PvTrend {
		trendMap[v.Date] = types.VisitTrendVO{
			Date:    v.Date,
			PvCount: v.Count,
			UvCount: 0, // 初始化为0
		}
		dates = append(dates, v.Date)
	}

	// 处理 UV 趋势数据
	for _, v := range visit.UvTrend {
		if data, exists := trendMap[v.Date]; exists {
			data.UvCount = v.Count
			trendMap[v.Date] = data
		} else {
			trendMap[v.Date] = types.VisitTrendVO{
				Date:    v.Date,
				PvCount: 0, // 初始化为0
				UvCount: v.Count,
			}
			dates = append(dates, v.Date)
		}
	}

	// 对日期进行排序
	sort.Strings(dates)

	// 构建返回数据
	trendList := make([]types.VisitTrendVO, 0, len(dates))
	for _, date := range dates {
		trendList = append(trendList, trendMap[date])
	}

	return &types.GetVisitTrendResp{
		VisitTrend: trendList,
	}, nil
}
