package logic

import (
	"encoding/json"

	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/request"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/response"
	"github.com/ve-weiyi/ve-blog-golang/server/api/service/svc"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/system"
)

type WebsiteService struct {
	svcCtx *svc.ServiceContext
}

func NewWebsiteService(svcCtx *svc.ServiceContext) *WebsiteService {
	return &WebsiteService{
		svcCtx: svcCtx,
	}
}

func (s *WebsiteService) GetBlogHomeInfo(reqCtx *request.Context, data interface{}) (resp *response.BlogHomeInfo, err error) {
	articleCount, _ := s.svcCtx.ArticleRepository.Count(reqCtx, "`is_delete` = ?", entity.False)
	categoryCount, _ := s.svcCtx.CategoryRepository.Count(reqCtx, "")
	tagCount, _ := s.svcCtx.TagRepository.Count(reqCtx, "")
	pages, _ := s.svcCtx.PageRepository.FindALL(reqCtx, "")
	config, err := s.svcCtx.WebsiteConfigRepository.First(reqCtx, "`key` = ?", "website_config")

	resp = &response.BlogHomeInfo{
		ArticleCount:  articleCount,
		CategoryCount: categoryCount,
		TagCount:      tagCount,
		ViewsCount:    "1",
		PageList:      convertPageList(pages),
	}

	json.Unmarshal([]byte(config.Config), &resp.WebsiteConfig)
	return resp, err
}

func (s *WebsiteService) GetAdminHomeInfo(reqCtx *request.Context, data interface{}) (resp *response.AdminHomeInfo, err error) {
	// 查询消息数量
	msgCount, err := s.svcCtx.RemarkRepository.Count(reqCtx, "")
	if err != nil {
		return nil, err
	}

	// 查询用户数量
	userCount, err := s.svcCtx.UserAccountRepository.Count(reqCtx, "")
	if err != nil {
		return nil, err
	}

	// 查询文章数量
	articles, err := s.svcCtx.ArticleRepository.FindALL(reqCtx, "")
	if err != nil {
		return nil, err
	}

	// 查询分类数量
	categories, err := s.svcCtx.CategoryRepository.FindALL(reqCtx, "")
	if err != nil {
		return nil, err
	}

	// 查询标签数量
	tags, err := s.svcCtx.TagRepository.FindALL(reqCtx, "")
	if err != nil {
		return nil, err
	}

	uniqueViews, err := s.svcCtx.UniqueViewRepository.FindALL(reqCtx, "")
	if err != nil {
		return nil, err
	}

	articleCount, err := s.svcCtx.ArticleRepository.Count(reqCtx, "")
	if err != nil {
		return nil, err
	}
	resp = &response.AdminHomeInfo{
		ViewsCount:            10,
		MessageCount:          msgCount,
		UserCount:             userCount,
		ArticleCount:          articleCount,
		CategoryDTOList:       convertCategoryList(categories),
		TagDTOList:            convertTagList(tags),
		ArticleStatisticsList: convertArticleStatisticsList(articles),
		UniqueViewDTOList:     convertUniqueViewList(uniqueViews),
		ArticleRankDTOList:    convertArticleRankList(articles),
	}

	return resp, err
}

func (s *WebsiteService) GetSystemState(reqCtx *request.Context, req interface{}) (server *system.Server, err error) {
	var sv system.Server
	sv.Os = system.InitOS()
	if sv.Cpu, err = system.InitCPU(); err != nil {
		return &sv, err
	}
	if sv.Ram, err = system.InitRAM(); err != nil {
		return &sv, err
	}
	if sv.Disk, err = system.InitDisk(); err != nil {
		return &sv, err
	}

	return &sv, nil
}

func (s *WebsiteService) GetAboutMe(reqCtx *request.Context, req interface{}) (resp string, err error) {
	config, err := s.svcCtx.WebsiteConfigRepository.First(reqCtx, "`key` = ?", "about")
	if err != nil {
		return "", err
	}

	return config.Config, err
}

func (s *WebsiteService) UpdateAboutMe(reqCtx *request.Context, req string) (resp string, err error) {
	config, err := s.svcCtx.WebsiteConfigRepository.First(reqCtx, "`key` = ?", "about")
	if err != nil {
		return "", err
	}
	// 更新
	config.Config = req
	_, err = s.svcCtx.WebsiteConfigRepository.Update(reqCtx, config)
	if err != nil {
		return "", err
	}

	return config.Config, err
}

func (s *WebsiteService) GetWebsiteConfig(reqCtx *request.Context, req interface{}) (resp string, err error) {
	config, err := s.svcCtx.WebsiteConfigRepository.First(reqCtx, "`key` = ?", "website_config")
	if err != nil {
		return "", err
	}

	return config.Config, err
}

func (s *WebsiteService) GetConfig(reqCtx *request.Context, req *request.WebsiteConfigReq) (resp string, err error) {
	config, err := s.svcCtx.WebsiteConfigRepository.First(reqCtx, "`key` = ?", req.Key)
	if err != nil {
		return "", err
	}

	return config.Config, err
}

func (s *WebsiteService) UpdateConfig(reqCtx *request.Context, req *request.WebsiteConfigReq) (resp string, err error) {
	config, err := s.svcCtx.WebsiteConfigRepository.First(reqCtx, "`key` = ?", req.Key)
	if err != nil {
		return "", err
	}
	// 更新
	config.Config = req.Value
	_, err = s.svcCtx.WebsiteConfigRepository.Update(reqCtx, config)
	if err != nil {
		return "", err
	}

	return config.Config, err
}
