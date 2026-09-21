package api

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/accessservice"
)

type CleanApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 清空接口列表
func NewCleanApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CleanApiLogic {
	return &CleanApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CleanApiLogic) CleanApi(req *types.EmptyReq) (resp *types.CleanApiResp, err error) {
	out, err := l.svcCtx.AccessService.CleanApis(l.ctx, &accessservice.CleanApisRequest{})
	if err != nil {
		return nil, err
	}

	return &types.CleanApiResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
