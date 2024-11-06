package articlerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type TopArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTopArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TopArticleLogic {
	return &TopArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 置顶文章
func (l *TopArticleLogic) TopArticle(in *articlerpc.TopArticleReq) (*articlerpc.EmptyResp, error) {

	record, err := l.svcCtx.TArticleModel.FindOne(l.ctx, in.ArticleId)
	if err != nil {
		return nil, err
	}

	record.IsTop = in.IsTop
	_, err = l.svcCtx.TArticleModel.Save(l.ctx, record)
	if err != nil {
		return nil, err
	}

	return &articlerpc.EmptyResp{}, nil
}
