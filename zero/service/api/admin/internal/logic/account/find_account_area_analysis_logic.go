package account

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/accountrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindAccountAreaAnalysisLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户分布地区
func NewFindAccountAreaAnalysisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindAccountAreaAnalysisLogic {
	return &FindAccountAreaAnalysisLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindAccountAreaAnalysisLogic) FindAccountAreaAnalysis(req *types.AccountQuery) (resp *types.PageResp, err error) {
	in := &accountrpc.EmptyReq{}
	// 查询用户数量
	users, err := l.svcCtx.AccountRpc.GetUserAreasAnalysis(l.ctx, in)
	if err != nil {
		return nil, err
	}

	// 分类
	var list []*types.AccountArea
	for _, item := range users.List {
		m := &types.AccountArea{
			Name:  item.Region,
			Value: item.Count,
		}

		list = append(list, m)
	}

	resp = &types.PageResp{
		Page:     0,
		PageSize: 0,
		Total:    0,
		List:     list,
	}

	return
}
