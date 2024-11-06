package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/dto"
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

// 文章归档(时间轴)
func (s *ArticleService) FindArticleArchives(reqCtx *request.Context, in *dto.ArticleArchivesQueryReq) (out *dto.PageResp, err error) {
	// todo

	return
}

// 通过分类获取文章列表
func (s *ArticleService) FindArticleClassifyCategory(reqCtx *request.Context, in *dto.ArticleClassifyQueryReq) (out *dto.PageResp, err error) {
	// todo

	return
}

// 通过标签获取文章列表
func (s *ArticleService) FindArticleClassifyTag(reqCtx *request.Context, in *dto.ArticleClassifyQueryReq) (out *dto.PageResp, err error) {
	// todo

	return
}

// 获取文章详情
func (s *ArticleService) GetArticleDetails(reqCtx *request.Context, in *dto.IdReq) (out *dto.ArticleDeatils, err error) {
	// todo

	return
}

// 获取首页文章列表
func (s *ArticleService) FindArticleHomeList(reqCtx *request.Context, in *dto.ArticleHomeQueryReq) (out *dto.PageResp, err error) {
	// todo

	return
}

// 获取首页推荐文章列表
func (s *ArticleService) FindArticleRecommend(reqCtx *request.Context, in *dto.EmptyReq) (out *dto.PageResp, err error) {
	// todo

	return
}

// 点赞文章
func (s *ArticleService) LikeArticle(reqCtx *request.Context, in *dto.IdReq) (out *dto.EmptyResp, err error) {
	// todo

	return
}
