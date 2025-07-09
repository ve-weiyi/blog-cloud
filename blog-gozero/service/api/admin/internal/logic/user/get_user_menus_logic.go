package user

import (
	"context"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/restx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/permissionrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户菜单权限
func NewGetUserMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserMenusLogic {
	return &GetUserMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserMenusLogic) GetUserMenus(req *types.EmptyReq) (resp *types.UserMenusResp, err error) {
	in := &permissionrpc.UserIdReq{
		UserId: cast.ToString(l.ctx.Value(restx.HeaderUid)),
	}

	out, err := l.svcCtx.PermissionRpc.FindUserMenus(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var list []*types.UserMenu
	for _, v := range out.List {
		list = append(list, ConvertUserMenuTypes(v))
	}

	resp = &types.UserMenusResp{}
	resp.List = list
	return
}

func ConvertUserMenuTypes(in *permissionrpc.MenuDetailsResp) (out *types.UserMenu) {
	children := make([]*types.UserMenu, 0)
	if in.Children != nil {
		for _, v := range in.Children {
			children = append(children, ConvertUserMenuTypes(v))
		}
	}

	out = &types.UserMenu{
		Id:        in.Id,
		ParentId:  in.ParentId,
		Path:      in.Path,
		Name:      in.Name,
		Component: in.Component,
		Redirect:  in.Redirect,
		Meta: types.UserMenuMeta{
			Title:      in.Meta.Title,
			Icon:       in.Meta.Icon,
			Hidden:     in.Meta.IsHidden == 1,
			AlwaysShow: in.Meta.AlwaysShow == 1,
			Affix:      false,
			KeepAlive:  false,
			Breadcrumb: false,
		},
		Children:  children,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}

	return
}
