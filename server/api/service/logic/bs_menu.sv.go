package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/request"
	"github.com/ve-weiyi/ve-blog-golang/server/api/service/svc"
)

type MenuService struct {
	svcCtx *svc.ServiceContext
}

func NewMenuService(svcCtx *svc.ServiceContext) *MenuService {
	return &MenuService{
		svcCtx: svcCtx,
	}
}

// 创建Menu记录
func (s *MenuService) CreateMenu(reqCtx *request.Context, menu *entity.Menu) (data *entity.Menu, err error) {
	return s.svcCtx.MenuRepository.Create(reqCtx, menu)
}

// 更新Menu记录
func (s *MenuService) UpdateMenu(reqCtx *request.Context, menu *entity.Menu) (data *entity.Menu, err error) {
	return s.svcCtx.MenuRepository.Update(reqCtx, menu)
}

// 删除Menu记录
func (s *MenuService) DeleteMenu(reqCtx *request.Context, id int) (rows int64, err error) {
	return s.svcCtx.MenuRepository.Delete(reqCtx, "id = ?", id)
}

// 查询Menu记录
func (s *MenuService) FindMenu(reqCtx *request.Context, id int) (data *entity.Menu, err error) {
	return s.svcCtx.MenuRepository.First(reqCtx, "id = ?", id)
}

// 批量删除Menu记录
func (s *MenuService) DeleteMenuByIds(reqCtx *request.Context, ids []int) (rows int64, err error) {
	return s.svcCtx.MenuRepository.Delete(reqCtx, "id in (?)", ids)
}

// 分页获取Menu记录
func (s *MenuService) FindMenuList(reqCtx *request.Context, page *request.PageQuery) (list []*entity.Menu, total int64, err error) {
	cond, args := page.ConditionClause()
	order := page.OrderClause()

	list, err = s.svcCtx.MenuRepository.FindList(reqCtx, page.Page, page.PageSize, order, cond, args...)
	if err != nil {
		return nil, 0, err
	}
	total, err = s.svcCtx.MenuRepository.Count(reqCtx, cond, args...)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
