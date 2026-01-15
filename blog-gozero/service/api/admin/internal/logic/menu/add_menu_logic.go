package menu

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/jsonconv"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/permissionrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建菜单
func NewAddMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMenuLogic {
	return &AddMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddMenuLogic) AddMenu(req *types.NewMenuReq) (resp *types.MenuBackVO, err error) {
	in := convertMenuPb(req)
	out, err := l.svcCtx.PermissionRpc.AddMenu(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = convertMenuTypes(out.Menu)
	return resp, nil
}

func convertMenuPb(in *types.NewMenuReq) (out *permissionrpc.AddMenuReq) {
	var children []*permissionrpc.AddMenuReq
	if in.Children != nil {
		for _, v := range in.Children {
			children = append(children, convertMenuPb(v))
		}
	}

	out = &permissionrpc.AddMenuReq{
		Id:        in.Id,
		ParentId:  in.ParentId,
		Path:      in.Path,
		Name:      in.Name,
		Component: in.Component,
		Redirect:  in.Redirect,
		Children:  children,
		Meta: &permissionrpc.MenuMeta{
			Type:       in.Type,
			Title:      in.Title,
			Icon:       in.Icon,
			Rank:       in.Rank,
			Perm:       in.Perm,
			Params:     jsonconv.AnyToJsonNE(in.Params),
			KeepAlive:  in.KeepAlive,
			AlwaysShow: in.AlwaysShow,
			Visible:    in.Visible,
			Status:     in.Status,
		},
	}

	return
}

func convertMenuTypes(in *permissionrpc.Menu) (out *types.MenuBackVO) {
	var children []*types.MenuBackVO
	if in.Children != nil {
		for _, v := range in.Children {
			children = append(children, convertMenuTypes(v))
		}
	}

	var params []*types.MenuMetaParams
	jsonconv.JsonToAny(in.Meta.Params, &params)

	out = &types.MenuBackVO{
		Id:        in.Id,
		ParentId:  in.ParentId,
		Path:      in.Path,
		Name:      in.Name,
		Component: in.Component,
		Redirect:  in.Redirect,
		MenuMeta: types.MenuMeta{
			Type:       in.Meta.Type,
			Title:      in.Meta.Title,
			Icon:       in.Meta.Icon,
			Rank:       in.Meta.Rank,
			Perm:       in.Meta.Perm,
			Params:     params,
			KeepAlive:  in.Meta.KeepAlive,
			AlwaysShow: in.Meta.AlwaysShow,
			Visible:    in.Meta.Visible,
			Status:     in.Meta.Status,
		},
		Children:  children,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}

	return
}
