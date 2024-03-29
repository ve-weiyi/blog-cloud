package logic

import (
	"math/rand"

	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/response"
	"github.com/ve-weiyi/ve-blog-golang/server/config/properties"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/oauth"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/jsonconv"
)

func convertAuthConfig(conf properties.AuthConfig) *oauth.AuthConfig {
	return &oauth.AuthConfig{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		RedirectUri:  conf.RedirectUri,
	}
}

func convertLoginHistory(entity *entity.UserLoginHistory) *response.LoginHistory {
	return &response.LoginHistory{
		ID:        entity.ID,
		LoginType: entity.LoginType,
		Agent:     entity.Agent,
		IpAddress: entity.IpAddress,
		IpSource:  entity.IpSource,
		LoginTime: entity.CreatedAt.String(),
	}
}

func convertArticle(article *entity.Article) response.ArticleDTO {
	out := response.ArticleDTO{
		ID:             article.ID,
		ArticleCover:   article.ArticleCover,
		ArticleTitle:   article.ArticleTitle,
		ArticleContent: article.ArticleContent,
		LikeCount:      100,
		ViewsCount:     200,
		IsTop:          article.IsTop,
		IsDelete:       article.IsDelete,
		Type:           article.Type,
		Status:         article.Status,
		OriginalURL:    article.OriginalURL,
		CreatedAt:      article.CreatedAt,
		UpdatedAt:      article.UpdatedAt,
	}
	return out
}

func convertArticlePreviewList(list []*entity.Article) []*response.ArticlePreviewDTO {
	var out []*response.ArticlePreviewDTO
	for _, item := range list {
		at := &response.ArticlePreviewDTO{
			ID:           item.ID,
			ArticleCover: item.ArticleCover,
			ArticleTitle: item.ArticleTitle,
			CreatedAt:    item.CreatedAt,
		}
		out = append(out, at)
	}

	return out
}

func convertArticlePreview(article *entity.Article) *response.ArticlePreviewDTO {
	if article == nil {
		return nil
	}
	return &response.ArticlePreviewDTO{
		ID:           article.ID,
		ArticleCover: article.ArticleCover,
		ArticleTitle: article.ArticleTitle,
	}
}

func convertArticleStatisticsList(list []*entity.Article) []*response.ArticleStatisticsDTO {
	var dateMap = make(map[string]int64)
	for _, item := range list {
		key := item.CreatedAt.Format("2006-01-02")
		if _, ok := dateMap[key]; ok {
			dateMap[key]++
		} else {
			dateMap[key] = 1
		}
	}
	// 日期统计
	var out []*response.ArticleStatisticsDTO
	for k, v := range dateMap {
		at := &response.ArticleStatisticsDTO{
			Day:   k,
			Count: v,
		}
		out = append(out, at)
	}

	return out
}

func convertUniqueViewList(list []*entity.UniqueView) []*response.UniqueViewDTO {
	var out []*response.UniqueViewDTO
	for _, item := range list {
		at := &response.UniqueViewDTO{
			Day:   item.CreatedAt.Format("2006-01-02"),
			Count: rand.Int63n(100),
		}
		out = append(out, at)
	}

	return out
}

func convertArticleRankList(list []*entity.Article) []*response.ArticleRankDTO {
	var out []*response.ArticleRankDTO
	for _, item := range list {
		at := &response.ArticleRankDTO{
			ID:           item.ID,
			ArticleTitle: item.ArticleTitle,
			Count:        rand.Int63n(100),
		}
		out = append(out, at)
	}

	return out
}

func convertCategory(entity *entity.Category) *response.CategoryDTO {
	if entity == nil {
		return &response.CategoryDTO{}
	}
	return &response.CategoryDTO{
		ID:           entity.ID,
		CategoryName: entity.CategoryName,
	}
}

func convertTagList(list []*entity.Tag) []*response.TagDTO {
	var tagList []*response.TagDTO
	for _, tag := range list {
		tagList = append(tagList, &response.TagDTO{
			ID:      tag.ID,
			TagName: tag.TagName,
		})
	}
	return tagList
}

func convertCategoryList(list []*entity.Category) []*response.CategoryDTO {
	var categoryList []*response.CategoryDTO

	for _, in := range list {
		data := &response.CategoryDTO{
			ID:           in.ID,
			CategoryName: in.CategoryName,
		}
		categoryList = append(categoryList, data)
	}
	return categoryList
}

func convertPageList(list []*entity.Page) []*response.PageDTO {
	var pageList []*response.PageDTO

	for _, in := range list {
		data := &response.PageDTO{
			ID:        in.ID,
			PageName:  in.PageName,
			PageLabel: in.PageLabel,
			PageCover: in.PageCover,
		}
		pageList = append(pageList, data)
	}
	return pageList
}

func convertRoleList(list []*entity.Role) []*response.RoleDTO {
	var roleList []*response.RoleDTO

	for _, in := range list {
		data := &response.RoleDTO{
			RoleName:    in.RoleName,
			RoleComment: in.RoleComment,
		}
		roleList = append(roleList, data)
	}
	return roleList
}

func convertMenu(entity *entity.Menu) *response.MenuDetailsDTO {
	if entity == nil {
		return nil
	}

	meta := response.Meta{}
	jsonconv.JsonToObject(entity.Meta, &meta)

	return &response.MenuDetailsDTO{
		ID:        entity.ID,
		Type:      entity.Type,
		Name:      entity.Name,
		Path:      entity.Path,
		Component: entity.Component,
		Title:     entity.Title,
		ParentID:  entity.ParentID,
		Meta:      meta,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		Children:  nil,
	}
}
