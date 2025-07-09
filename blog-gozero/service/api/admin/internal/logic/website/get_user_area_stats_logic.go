package website

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserAreaStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户分布地区
func NewGetUserAreaStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserAreaStatsLogic {
	return &GetUserAreaStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserAreaStatsLogic) GetUserAreaStats(req *types.GetUserAreaStatsReq) (resp *types.GetUserAreaStatsResp, err error) {
	in := &accountrpc.AnalysisUserAreasReq{
		UserType: req.UserType,
	}
	// 查询用户数量
	users, err := l.svcCtx.AccountRpc.AnalysisUserAreas(l.ctx, in)
	if err != nil {
		return nil, err
	}

	// 分类
	var list []*types.UserAreaVO
	for _, v := range users.List {
		m := &types.UserAreaVO{
			Name:  v.Area,
			Value: v.Count,
		}

		list = append(list, m)
	}

	return &types.GetUserAreaStatsResp{
		UserAreas:    list,
		TouristAreas: list,
	}, nil
}
