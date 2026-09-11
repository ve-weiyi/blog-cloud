package article

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
)

type ExportArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 导出文章列表
func NewExportArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportArticleLogic {
	return &ExportArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportArticleLogic) ExportArticle(req *types.ExportArticleReq) (resp *types.EmptyResp, err error) {
	return &types.EmptyResp{}, nil
}
