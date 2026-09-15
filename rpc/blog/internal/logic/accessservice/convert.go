package accessservicelogic

import (
	"time"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

func convertApiIn(in *accessrpc.CreateApiRequest) *model.TApi {
	if in == nil {
		return nil
	}
	return &model.TApi{
		ParentId:  in.ParentId,
		Path:      in.Path,
		Name:      in.Name,
		Method:    in.Method,
		Traceable: in.Traceable,
		Status:    in.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func convertApiOut(in *model.TApi) *accessrpc.Api {
	if in == nil {
		return nil
	}
	return &accessrpc.Api{
		Id:        in.Id,
		ParentId:  in.ParentId,
		Path:      in.Path,
		Name:      in.Name,
		Method:    in.Method,
		Traceable: in.Traceable,
		Status:    in.Status,
		CreatedAt: in.CreatedAt.UnixMilli(),
		UpdatedAt: in.UpdatedAt.UnixMilli(),
	}
}

func convertMenuIn(in *accessrpc.CreateMenuRequest) *model.TMenu {
	if in == nil {
		return nil
	}
	out := &model.TMenu{
		ParentId:  in.ParentId,
		Path:      in.Path,
		Name:      in.Name,
		Component: in.Component,
		Redirect:  in.Redirect,
	}
	if in.Meta != nil {
		out.Type = in.Meta.Type
		out.Title = in.Meta.Title
		out.Icon = in.Meta.Icon
		out.Rank = in.Meta.Rank
		out.Perm = in.Meta.Perm
		out.Params = in.Meta.Params
		out.KeepAlive = in.Meta.KeepAlive
		out.AlwaysShow = in.Meta.AlwaysShow
		out.Visible = in.Meta.Visible
		out.Status = in.Meta.Status
	}
	return out
}

func convertMenuOut(in *model.TMenu) *accessrpc.Menu {
	if in == nil {
		return nil
	}
	return &accessrpc.Menu{
		Id:        in.Id,
		ParentId:  in.ParentId,
		Path:      in.Path,
		Name:      in.Name,
		Component: in.Component,
		Redirect:  in.Redirect,
		Meta: &accessrpc.MenuMeta{
			Type:       in.Type,
			Title:      in.Title,
			Icon:       in.Icon,
			Rank:       in.Rank,
			Perm:       in.Perm,
			Params:     in.Params,
			KeepAlive:  in.KeepAlive,
			AlwaysShow: in.AlwaysShow,
			Visible:    in.Visible,
			Status:     in.Status,
		},
		CreatedAt: in.CreatedAt.UnixMilli(),
		UpdatedAt: in.UpdatedAt.UnixMilli(),
	}
}

func convertRoleOut(in *model.TRole) *accessrpc.Role {
	if in == nil {
		return nil
	}
	return &accessrpc.Role{
		Id:          in.Id,
		ParentId:    in.ParentId,
		RoleKey:     in.RoleKey,
		RoleLabel:   in.RoleLabel,
		RoleComment: in.RoleComment,
		IsDefault:   in.IsDefault,
		Status:      in.Status,
		CreatedAt:   in.CreatedAt.UnixMilli(),
		UpdatedAt:   in.UpdatedAt.UnixMilli(),
	}
}
