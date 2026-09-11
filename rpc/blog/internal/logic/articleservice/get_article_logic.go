package articleservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

type GetArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleLogic {
	return &GetArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询文章详情
func (l *GetArticleLogic) GetArticle(in *articlerpc.GetArticleRequest) (*articlerpc.GetArticleResponse, error) {
	helper := NewArticleHelper(l.ctx, l.svcCtx)
	record, err := l.svcCtx.TArticleModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	list, err := helper.convertArticleOut([]*model.TArticle{record})
	if err != nil {
		return nil, err
	}

	return &articlerpc.GetArticleResponse{
		Article: list[0],
	}, nil
}
