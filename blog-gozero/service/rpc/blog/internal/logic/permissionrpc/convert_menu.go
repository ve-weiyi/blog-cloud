package permissionrpclogic

import (
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
)

func convertMenuIn(in *permissionrpc.MenuNewReq) (out *model.TMenu) {
	out = &model.TMenu{
		Id:         in.Id,
		ParentId:   in.ParentId,
		Path:       in.Path,
		Name:       in.Name,
		Component:  in.Component,
		Redirect:   in.Redirect,
		Type:       in.Meta.Type,
		Title:      in.Meta.Title,
		Icon:       in.Meta.Icon,
		Rank:       in.Meta.Rank,
		Perm:       in.Meta.Perm,
		Params:     in.Meta.Params,
		KeepAlive:  in.Meta.KeepAlive,
		AlwaysShow: in.Meta.AlwaysShow,
		IsHidden:   in.Meta.IsHidden,
		IsDisable:  in.Meta.IsDisable,
		Extra:      jsonconv.AnyToJsonNE(in.Meta),
	}

	return out
}

func convertMenuOut(in *model.TMenu) (out *permissionrpc.MenuDetailsResp) {

	out = &permissionrpc.MenuDetailsResp{
		Id:        in.Id,
		ParentId:  in.ParentId,
		Path:      in.Path,
		Name:      in.Name,
		Component: in.Component,
		Redirect:  in.Redirect,
		CreatedAt: in.CreatedAt.Unix(),
		UpdatedAt: in.UpdatedAt.Unix(),
		Children:  nil,
		Meta: &permissionrpc.MenuMeta{
			Type:       in.Type,
			Title:      in.Title,
			Icon:       in.Icon,
			Rank:       in.Rank,
			Perm:       in.Perm,
			Params:     in.Params,
			KeepAlive:  in.KeepAlive,
			AlwaysShow: in.AlwaysShow,
			IsHidden:   in.IsHidden,
			IsDisable:  in.IsDisable,
		},
	}
	return out
}
