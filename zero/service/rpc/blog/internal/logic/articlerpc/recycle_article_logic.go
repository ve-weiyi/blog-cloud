package articlerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecycleArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecycleArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecycleArticleLogic {
	return &RecycleArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 回收文章
func (l *RecycleArticleLogic) RecycleArticle(in *articlerpc.RecycleArticleReq) (*articlerpc.EmptyResp, error) {
	record, err := l.svcCtx.TArticleModel.FindOne(l.ctx, in.ArticleId)
	if err != nil {
		return nil, err
	}

	record.IsDelete = in.IsDelete
	_, err = l.svcCtx.TArticleModel.Save(l.ctx, record)
	if err != nil {
		return nil, err
	}

	return &articlerpc.EmptyResp{}, nil
}
