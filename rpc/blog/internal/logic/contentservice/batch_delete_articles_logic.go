package contentservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/contentrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDeleteArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteArticlesLogic {
	return &BatchDeleteArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除文章
func (l *BatchDeleteArticlesLogic) BatchDeleteArticles(in *contentrpc.BatchDeleteArticlesRequest) (*contentrpc.BatchDeleteArticlesResponse, error) {
	rows, err := l.svcCtx.TArticleModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &contentrpc.BatchDeleteArticlesResponse{
		SuccessCount: rows,
	}, nil
}
