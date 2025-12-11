package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type MenuLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewMenuLogic(svcCtx *svctx.ServiceContext) *MenuLogic {
	return &MenuLogic{
		svcCtx: svcCtx,
	}
}

// 创建菜单
func (s *MenuLogic) AddMenu(reqCtx *request.Context, in *types.MenuNewReq) (out *types.MenuBackVO, err error) {
	// todo

	return
}

// 清空菜单列表
func (s *MenuLogic) CleanMenuList(reqCtx *request.Context, in *types.EmptyReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 删除菜单
func (s *MenuLogic) DeletesMenu(reqCtx *request.Context, in *types.IdsReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 分页获取菜单列表
func (s *MenuLogic) FindMenuList(reqCtx *request.Context, in *types.MenuQuery) (out *types.PageResp, err error) {
	// todo

	return
}

// 同步菜单列表
func (s *MenuLogic) SyncMenuList(reqCtx *request.Context, in *types.SyncMenuReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 更新菜单
func (s *MenuLogic) UpdateMenu(reqCtx *request.Context, in *types.MenuNewReq) (out *types.MenuBackVO, err error) {
	// todo

	return
}
