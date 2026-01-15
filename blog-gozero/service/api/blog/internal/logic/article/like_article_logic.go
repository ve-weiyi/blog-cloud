package article

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 点赞文章
func NewLikeArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeArticleLogic {
	return &LikeArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeArticleLogic) LikeArticle(req *types.IdReq) (resp *types.EmptyResp, err error) {
	in := &articlerpc.LikeArticleReq{Id: req.Id}

	_, err = l.svcCtx.ArticleRpc.LikeArticle(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return
}
