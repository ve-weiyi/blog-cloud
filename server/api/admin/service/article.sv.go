package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type ArticleService struct {
	svcCtx *svctx.ServiceContext
}

func NewArticleService(svcCtx *svctx.ServiceContext) *ArticleService {
	return &ArticleService{
		svcCtx: svcCtx,
	}
}

// 添加文章
func (s *ArticleService) AddArticle(reqCtx *request.Context, in *dto.ArticleNewReq) (out *dto.ArticleBackDTO, err error) {
	// todo

	return
}

// 删除文章
func (s *ArticleService) DeleteArticle(reqCtx *request.Context, in *dto.IdReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 导出文章列表
func (s *ArticleService) ExportArticleList(reqCtx *request.Context, in *dto.IdsReq) (out *dto.EmptyResp, err error) {
	// todo

	return
}

// 查询文章列表
func (s *ArticleService) FindArticleList(reqCtx *request.Context, in *dto.ArticleQuery) (out *dto.PageResp, err error) {
	// todo

	return
}

// 查询文章
func (s *ArticleService) GetArticle(reqCtx *request.Context, in *dto.IdReq) (out *dto.ArticleBackDTO, err error) {
	// todo

	return
}

// 回收文章
func (s *ArticleService) RecycleArticle(reqCtx *request.Context, in *dto.ArticleRecycleReq) (out *dto.EmptyResp, err error) {
	// todo

	return
}

// 置顶文章
func (s *ArticleService) TopArticle(reqCtx *request.Context, in *dto.ArticleTopReq) (out *dto.EmptyResp, err error) {
	// todo

	return
}

// 保存文章
func (s *ArticleService) UpdateArticle(reqCtx *request.Context, in *dto.ArticleNewReq) (out *dto.ArticleBackDTO, err error) {
	// todo

	return
}
