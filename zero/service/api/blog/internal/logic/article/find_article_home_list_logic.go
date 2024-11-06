package article

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/articlerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindArticleHomeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取首页文章列表
func NewFindArticleHomeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindArticleHomeListLogic {
	return &FindArticleHomeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindArticleHomeListLogic) FindArticleHomeList(req *types.ArticleHomeQueryReq) (resp *types.PageResp, err error) {
	in := &articlerpc.FindArticleListReq{
		Page:         req.Page,
		PageSize:     req.PageSize,
		Sorts:        req.Sorts,
		ArticleTitle: req.ArticleTitle,
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
