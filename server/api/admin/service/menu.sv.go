package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type MenuService struct {
	svcCtx *svctx.ServiceContext
}

func NewMenuService(svcCtx *svctx.ServiceContext) *MenuService {
	return &MenuService{
		svcCtx: svcCtx,
	}
}

// 创建菜单
func (s *MenuService) AddMenu(reqCtx *request.Context, in *dto.MenuNewReq) (out *dto.MenuBackDTO, err error) {
	// todo

	return
}

// 批量删除菜单
func (s *MenuService) BatchDeleteMenu(reqCtx *request.Context, in *dto.IdsReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 清空菜单列表
func (s *MenuService) CleanMenuList(reqCtx *request.Context, in *dto.EmptyReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 删除菜单
func (s *MenuService) DeleteMenu(reqCtx *request.Context, in *dto.IdReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 分页获取菜单列表
func (s *MenuService) FindMenuList(reqCtx *request.Context, in *dto.MenuQuery) (out *dto.PageResp, err error) {
	// todo

	return
}

// 同步菜单列表
func (s *MenuService) SyncMenuList(reqCtx *request.Context, in *dto.SyncMenuReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 更新菜单
func (s *MenuService) UpdateMenu(reqCtx *request.Context, in *dto.MenuNewReq) (out *dto.MenuBackDTO, err error) {
	// todo

	return
}
