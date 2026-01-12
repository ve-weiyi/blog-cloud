package permissionrpclogic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
)

func convertRoleIn(in *permissionrpc.NewRoleReq) (out *model.TRole) {
	out = &model.TRole{
		Id:          in.Id,
		ParentId:    in.ParentId,
		RoleKey:     in.RoleKey,
		RoleLabel:   in.RoleLabel,
		RoleComment: in.RoleComment,
		IsDisable:   in.IsDisable,
		IsDefault:   in.IsDefault,
	}

	return out
}

func convertRoleOut(in *model.TRole) (out *permissionrpc.RoleDetailsResp) {
	out = &permissionrpc.RoleDetailsResp{
		Id:          in.Id,
		ParentId:    in.ParentId,
		RoleKey:     in.RoleKey,
		RoleLabel:   in.RoleLabel,
		RoleComment: in.RoleComment,
		IsDisable:   in.IsDisable,
		IsDefault:   in.IsDefault,
		CreatedAt:   in.CreatedAt.Unix(),
		UpdatedAt:   in.UpdatedAt.Unix(),
	}

	return out
}
