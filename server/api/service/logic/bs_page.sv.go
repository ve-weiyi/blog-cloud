package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/request"
	"github.com/ve-weiyi/ve-blog-golang/server/api/service/svc"
)

type PageService struct {
	svcCtx *svc.ServiceContext
}

func NewPageService(svcCtx *svc.ServiceContext) *PageService {
	return &PageService{
		svcCtx: svcCtx,
	}
}

// 创建Page记录
func (s *PageService) CreatePage(reqCtx *request.Context, page *entity.Page) (data *entity.Page, err error) {
	return s.svcCtx.PageRepository.Create(reqCtx, page)
}

// 更新Page记录
func (s *PageService) UpdatePage(reqCtx *request.Context, page *entity.Page) (data *entity.Page, err error) {
	return s.svcCtx.PageRepository.Update(reqCtx, page)
}

// 删除Page记录
func (s *PageService) DeletePage(reqCtx *request.Context, id int) (rows int64, err error) {
	return s.svcCtx.PageRepository.Delete(reqCtx, "id = ?", id)
}

// 查询Page记录
func (s *PageService) FindPage(reqCtx *request.Context, id int) (data *entity.Page, err error) {
	return s.svcCtx.PageRepository.First(reqCtx, "id = ?", id)
}

// 批量删除Page记录
func (s *PageService) DeletePageByIds(reqCtx *request.Context, ids []int) (rows int64, err error) {
	return s.svcCtx.PageRepository.Delete(reqCtx, "id in (?)", ids)
}

// 分页获取Page记录
func (s *PageService) FindPageList(reqCtx *request.Context, page *request.PageQuery) (list []*entity.Page, total int64, err error) {
	cond, args := page.ConditionClause()
	order := page.OrderClause()

	list, err = s.svcCtx.PageRepository.FindList(reqCtx, page.Page, page.PageSize, order, cond, args...)
	if err != nil {
		return nil, 0, err
	}
	total, err = s.svcCtx.PageRepository.Count(reqCtx, cond, args...)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
