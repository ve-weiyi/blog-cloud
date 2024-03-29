package logic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
)

func (s *TagRepository) FindArticleTagList(ctx context.Context, articleId int) (list []*entity.Tag, err error) {
	// 创建db
	db := s.DbEngin.WithContext(ctx)
	var ats []*entity.ArticleTag
	var tags []*entity.Tag

	err = db.Where("article_id = ?", articleId).Find(&ats).Error
	if err != nil {
		return nil, err
	}

	var tagIds []int
	for _, at := range ats {
		tagIds = append(tagIds, at.TagID)
	}

	err = db.Where("id in (?)", tagIds).Find(&tags).Error
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (s *TagRepository) FindArticleTagMap(ctx context.Context, articleIds []int) (mp map[int][]*entity.Tag, err error) {
	// 创建db
	db := s.DbEngin.WithContext(ctx)
	var ats []*entity.ArticleTag
	var tags []*entity.Tag

	// 查找所有文章关联的tag
	err = db.Where("article_id in ?", articleIds).Find(&ats).Error
	if err != nil {
		return nil, err
	}

	// 收集id
	var tagIds []int
	for _, at := range ats {
		tagIds = append(tagIds, at.TagID)
	}

	// 查询tag
	err = db.Where("id in (?)", tagIds).Find(&tags).Error
	if err != nil {
		return nil, err
	}

	// tag map
	var tagMap = make(map[int]*entity.Tag)
	for _, tag := range tags {
		tagMap[tag.ID] = tag
	}

	// article []tag map
	var atMap = make(map[int][]*entity.Tag)
	for _, at := range ats {
		if _, ok := atMap[at.ArticleID]; ok {
			atMap[at.ArticleID] = append(atMap[at.ArticleID], tagMap[at.TagID])
		} else {
			atMap[at.ArticleID] = []*entity.Tag{tagMap[at.TagID]}
		}
	}

	return atMap, nil
}

// 批量创建tag,当tag_name不存在时创建
func (s *TagRepository) BatchCreateTagNotExist(ctx context.Context, tagNames []string) (list []*entity.Tag, err error) {
	// 创建db
	db := s.DbEngin.WithContext(ctx)
	var tags []*entity.Tag
	// 查找所有文章关联的tag
	err = db.Where("tag_name in ?", tagNames).Find(&tags).Error
	if err != nil {
		return nil, err
	}

	// 收集id
	var tagNameMap = make(map[string]*entity.Tag)
	for _, tag := range tags {
		tagNameMap[tag.TagName] = tag
	}

	// 创建不存在的tag
	for _, tagName := range tagNames {
		if _, ok := tagNameMap[tagName]; !ok {
			tag := &entity.Tag{
				TagName: tagName,
			}
			err = db.Create(tag).Error
			if err != nil {
				return nil, err
			}
			tagNameMap[tagName] = tag
		}
	}

	// tag map
	for _, tag := range tagNameMap {
		list = append(list, tag)
	}

	return list, nil
}
