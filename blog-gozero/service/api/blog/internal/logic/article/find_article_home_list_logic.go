package article

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/constant"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"

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

func (l *FindArticleHomeListLogic) FindArticleHomeList(req *types.QueryArticleHomeReq) (resp *types.PageResp, err error) {
	in := &articlerpc.FindArticleListReq{
		Paginate: &articlerpc.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		ArticleTitle: req.ArticleTitle,
		IsTop:        constant.ArticleIsTopALL,
		IsDelete:     constant.ArticleIsDeleteNo,
		Status:       constant.ArticleStatusPublic,
	}

	out, err := l.svcCtx.ArticleRpc.FindArticleList(l.ctx, in)
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
	resp.Page = out.Pagination.Page
	resp.PageSize = out.Pagination.PageSize
	resp.Total = out.Pagination.Total
	resp.List = list
	return
}
