package overview

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/statsservice"
)

type GetGeoStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户地理分布
func NewGetGeoStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGeoStatsLogic {
	return &GetGeoStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGeoStatsLogic) GetGeoStats(req *types.GetGeoStatsReq) (resp *types.GetGeoStatsResp, err error) {
	out, err := l.svcCtx.StatsService.GetGeoStats(l.ctx, &statsservice.GetGeoStatsRequest{})
	if err != nil {
		return nil, err
	}

	resp = &types.GetGeoStatsResp{
		Users:    make([]*types.RegionStatVO, 0, len(out.Users)),
		Visitors: make([]*types.RegionStatVO, 0, len(out.Visitors)),
	}
	for _, r := range out.Users {
		resp.Users = append(resp.Users, &types.RegionStatVO{
			Region: r.Region,
			Count:  r.Count,
		})
	}
	for _, v := range out.Visitors {
		resp.Visitors = append(resp.Visitors, &types.RegionStatVO{
			Region: v.Region,
			Count:  v.Count,
		})
	}
	return
}
