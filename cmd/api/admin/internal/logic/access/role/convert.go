package role

import (
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/accessservice"
)

func convertRoleVO(in *accessservice.Role) *types.RoleVO {
	if in == nil {
		return nil
	}

	return &types.RoleVO{
		Id:          in.Id,
		ParentId:    in.ParentId,
		RoleKey:     in.RoleKey,
		RoleLabel:   in.RoleLabel,
		RoleComment: in.RoleComment,
		IsDefault:   in.IsDefault,
		Status:      in.Status,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	}
}
