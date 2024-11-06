package menu

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"

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

func (l *AddMenuLogic) AddMenu(req *types.MenuNewReq) (resp *types.MenuBackDTO, err error) {
	in := ConvertMenuPb(req)
	out, err := l.svcCtx.PermissionRpc.AddMenu(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = ConvertMenuTypes(out)
	return resp, nil
}

func ConvertMenuPb(in *types.MenuNewReq) (out *permissionrpc.MenuNewReq) {
	var children []*permissionrpc.MenuNewReq
	if in.Children != nil {
		for _, v := range in.Children {
			children = append(children, ConvertMenuPb(v))
		}
	}

	out = &permissionrpc.MenuNewReq{
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
			IsHidden:   in.IsHidden,
			IsDisable:  in.IsDisable,
		},
	}

	return
}

func ConvertMenuTypes(in *permissionrpc.MenuDetails) (out *types.MenuBackDTO) {
	var children []*types.MenuBackDTO
	if in.Children != nil {
		for _, v := range in.Children {
			children = append(children, ConvertMenuTypes(v))
		}
	}

	var params []*types.MenuMetaParams
	jsonconv.JsonToAny(in.Meta.Params, &params)

	out = &types.MenuBackDTO{
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
			IsHidden:   in.Meta.IsHidden,
			IsDisable:  in.Meta.IsDisable,
		},
		Children:  children,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}

	return
}
