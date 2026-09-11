package permissionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

type UpdateUserRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserRolesLogic {
	return &UpdateUserRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新用户角色
func (l *UpdateUserRolesLogic) UpdateUserRoles(in *permissionrpc.UpdateUserRolesRequest) (*permissionrpc.UpdateUserRolesResponse, error) {
	_, err := l.svcCtx.TUserRoleModel.DeleteBatch(l.ctx, "user_id = ?", in.UserId)
	if err != nil {
		return nil, err
	}

	if len(in.RoleIds) > 0 {
		var batch []*model.TUserRole
		for _, id := range in.RoleIds {
			batch = append(batch, &model.TUserRole{
				UserId: in.UserId,
				RoleId: id,
			})
		}
		_, err = l.svcCtx.TUserRoleModel.InsertBatch(l.ctx, batch...)
		if err != nil {
			return nil, err
		}
	}

	return &permissionrpc.UpdateUserRolesResponse{Success: true}, nil
}
