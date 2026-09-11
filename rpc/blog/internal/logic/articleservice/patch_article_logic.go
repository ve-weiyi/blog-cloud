package articleservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type PatchArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPatchArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PatchArticleLogic {
	return &PatchArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 部分更新文章
func (l *PatchArticleLogic) PatchArticle(in *articlerpc.PatchArticleRequest) (*articlerpc.PatchArticleResponse, error) {
	record, err := l.svcCtx.TArticleModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	if in.IsDelete != nil {
		record.IsDelete = *in.IsDelete
	}
	if in.IsTop != nil {
		record.IsTop = *in.IsTop
	}

	_, err = l.svcCtx.TArticleModel.Save(l.ctx, record)
	if err != nil {
		return nil, err
	}

	return &articlerpc.PatchArticleResponse{Success: true}, nil
}
