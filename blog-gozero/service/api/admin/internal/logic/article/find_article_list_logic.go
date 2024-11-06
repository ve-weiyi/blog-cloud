package article

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindArticleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询文章列表
func NewFindArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindArticleListLogic {
	return &FindArticleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindArticleListLogic) FindArticleList(req *types.ArticleQuery) (resp *types.PageResp, err error) {
	in := &articlerpc.FindArticleListReq{
		Page:         req.Page,
		PageSize:     req.PageSize,
		Sorts:        req.Sorts,
		ArticleTitle: req.ArticleTitle,
		ArticleType:  req.ArticleType,
		CategoryName: req.CategoryName,
		TagName:      req.TagName,
		IsTop:        req.IsTop,
		IsDelete:     req.IsDelete,
		Status:       req.Status,
	}

	out, err := l.svcCtx.ArticleRpc.FindArticleList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var list []*types.ArticleBackDTO
	for _, v := range out.List {
		m := ConvertArticleTypes(v)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = in.Page
	resp.PageSize = in.PageSize
	resp.Total = out.Total
	resp.List = list
	return resp, nil
}
