package permissionrpclogic

import (
	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/permissionrpc"
)

func convertRoleIn(in *permissionrpc.RoleNewReq) (out *model.TRole) {
	out = &model.TRole{
		Id:          in.Id,
		ParentId:    in.ParentId,
		RoleName:    in.RoleName,
		RoleLabel:   in.RoleLabel,
		RoleComment: in.RoleComment,
		IsDisable:   in.IsDisable,
		IsDefault:   in.IsDefault,
	}

	return out
}

func convertRoleOut(in *model.TRole) (out *permissionrpc.RoleDetails) {
	out = &permissionrpc.RoleDetails{
		Id:          in.Id,
		ParentId:    in.ParentId,
		RoleName:    in.RoleName,
		RoleLabel:   in.RoleLabel,
		RoleComment: in.RoleComment,
		IsDisable:   in.IsDisable,
		IsDefault:   in.IsDefault,
		CreatedAt:   in.CreatedAt.Unix(),
		UpdatedAt:   in.UpdatedAt.Unix(),
	}

	return out
}
