package article

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/articlerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindArticleClassifyCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 通过分类获取文章列表
func NewFindArticleClassifyCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindArticleClassifyCategoryLogic {
	return &FindArticleClassifyCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindArticleClassifyCategoryLogic) FindArticleClassifyCategory(req *types.ArticleClassifyQueryReq) (resp *types.PageResp, err error) {
	in := &articlerpc.FindArticleListReq{
		CategoryName: req.ClassifyName,
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
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	resp.Total = out.Total
	resp.List = list
	return
}
