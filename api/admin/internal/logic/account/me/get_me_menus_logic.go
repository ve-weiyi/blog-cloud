package me

import (
	"context"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizheader"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/accessservice"
)

type GetMeMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户菜单权限
func NewGetMeMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMeMenusLogic {
	return &GetMeMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMeMenusLogic) GetMeMenus(req *types.EmptyReq) (resp *types.GetMeMenusResp, err error) {
	userId := cast.ToString(l.ctx.Value(bizheader.HeaderUid))

	out, err := l.svcCtx.AccessService.GetUserMenus(l.ctx, &accessservice.GetUserMenusRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	var list []*types.UserMenu
	for _, v := range out.List {
		list = append(list, convertUserMenu(v))
	}

	return &types.GetMeMenusResp{List: list}, nil
}

func convertUserMenu(in *accessservice.Menu) *types.UserMenu {
	children := make([]*types.UserMenu, 0)
	if in.Children != nil {
		for _, v := range in.Children {
			children = append(children, convertUserMenu(v))
		}
	}

	return &types.UserMenu{
		Id:        in.Id,
		ParentId:  in.ParentId,
		Path:      in.Path,
		Name:      in.Name,
		Component: in.Component,
		Redirect:  in.Redirect,
		Meta: types.UserMenuMeta{
			Title:      in.Meta.Title,
			Icon:       in.Meta.Icon,
			Hidden:     in.Meta.Visible == 1,
			AlwaysShow: in.Meta.AlwaysShow == 1,
			Affix:      false,
			KeepAlive:  false,
			Breadcrumb: false,
		},
		Children:  children,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
}
