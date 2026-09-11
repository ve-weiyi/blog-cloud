package articleservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type DeleteArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteArticleLogic {
	return &DeleteArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除文章
func (l *DeleteArticleLogic) DeleteArticle(in *articlerpc.DeleteArticleRequest) (*articlerpc.DeleteArticleResponse, error) {
	rows, err := l.svcCtx.TArticleModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &articlerpc.DeleteArticleResponse{
		SuccessCount: rows,
	}, nil
}
