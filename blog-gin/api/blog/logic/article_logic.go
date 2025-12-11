package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
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

// 文章归档(时间轴)
func (s *ArticleLogic) FindArticleArchives(reqCtx *request.Context, in *types.ArticleArchivesQueryReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 通过分类获取文章列表
func (s *ArticleLogic) FindArticleClassifyCategory(reqCtx *request.Context, in *types.ArticleClassifyQueryReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 通过标签获取文章列表
func (s *ArticleLogic) FindArticleClassifyTag(reqCtx *request.Context, in *types.ArticleClassifyQueryReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 获取文章详情
func (s *ArticleLogic) GetArticleDetails(reqCtx *request.Context, in *types.IdReq) (out *types.ArticleDetails, err error) {
	// todo

	return
}

// 获取首页文章列表
func (s *ArticleLogic) FindArticleHomeList(reqCtx *request.Context, in *types.ArticleHomeQueryReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 获取首页推荐文章列表
func (s *ArticleLogic) FindArticleRecommend(reqCtx *request.Context, in *types.EmptyReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 点赞文章
func (s *ArticleLogic) LikeArticle(reqCtx *request.Context, in *types.IdReq) (out *types.EmptyResp, err error) {
	// todo

	return
}
