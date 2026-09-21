package api

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/accessservice"
)

type BatchDeleteApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除接口
func NewBatchDeleteApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteApiLogic {
	return &BatchDeleteApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteApiLogic) BatchDeleteApi(req *types.DeleteApiReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.AccessService.BatchDeleteApis(l.ctx, &accessservice.BatchDeleteApisRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
