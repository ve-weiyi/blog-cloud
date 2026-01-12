package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type ArticleLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewArticleLogic(svcCtx *svctx.ServiceContext) *ArticleLogic {
	return &ArticleLogic{
		svcCtx: svcCtx,
	}
}

// 添加文章
func (s *ArticleLogic) AddArticle(reqCtx *request.Context, in *types.NewArticleReq) (out *types.ArticleBackVO, err error) {
	// todo

	return
}

// 删除文章
func (s *ArticleLogic) DeleteArticle(reqCtx *request.Context, in *types.IdReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 导出文章列表
func (s *ArticleLogic) ExportArticleList(reqCtx *request.Context, in *types.IdsReq) (out *types.EmptyResp, err error) {
	// todo

	return
}

// 查询文章列表
func (s *ArticleLogic) FindArticleList(reqCtx *request.Context, in *types.QueryArticleReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 查询文章
func (s *ArticleLogic) GetArticle(reqCtx *request.Context, in *types.IdReq) (out *types.ArticleBackVO, err error) {
	// todo

	return
}

// 保存文章
func (s *ArticleLogic) UpdateArticle(reqCtx *request.Context, in *types.NewArticleReq) (out *types.ArticleBackVO, err error) {
	// todo

	return
}

// 更新文章删除状态
func (s *ArticleLogic) UpdateArticleDelete(reqCtx *request.Context, in *types.UpdateArticleDeleteReq) (out *types.EmptyResp, err error) {
	// todo

	return
}

// 更新文章置顶状态
func (s *ArticleLogic) UpdateArticleTop(reqCtx *request.Context, in *types.UpdateArticleTopReq) (out *types.EmptyResp, err error) {
	// todo

	return
}
