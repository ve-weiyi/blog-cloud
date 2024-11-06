package article

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/articlerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindArticleRecommendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取首页推荐文章列表
func NewFindArticleRecommendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindArticleRecommendLogic {
	return &FindArticleRecommendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindArticleRecommendLogic) FindArticleRecommend(req *types.EmptyReq) (resp *types.PageResp, err error) {
	in := &articlerpc.FindArticleListReq{
		IsTop: 1,
	}
	out, err := l.svcCtx.ArticleRpc.FindArticlePublicList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	list := make([]*types.ArticleHome, 0)
	// 转换数据
	for _, v := range out.List {
		m := ConvertArticleHomeTypes(v)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Total = out.Total
	resp.List = list
	return
}
