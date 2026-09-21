package page

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/siteservice"
)

type BatchDeletePageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除页面
func NewBatchDeletePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeletePageLogic {
	return &BatchDeletePageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeletePageLogic) BatchDeletePage(req *types.DeletePageReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.SiteService.BatchDeletePages(l.ctx, &siteservice.BatchDeletePagesRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
