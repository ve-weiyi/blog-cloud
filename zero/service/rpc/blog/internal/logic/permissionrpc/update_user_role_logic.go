package permissionrpclogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"
)

type UpdateUserRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserRoleLogic {
	return &UpdateUserRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改用户角色
func (l *UpdateUserRoleLogic) UpdateUserRole(in *permissionrpc.UpdateUserRoleReq) (*permissionrpc.EmptyResp, error) {
	ua, err := l.svcCtx.TUserModel.First(l.ctx, "id = ?", in.UserId)
	if err != nil {
		return nil, err
	}

	// 删除用户角色
	_, err = l.svcCtx.TUserRoleModel.Deletes(l.ctx, "user_id = ?", ua.UserId)
	if err != nil {
		return nil, err
	}

	var userRoles []*model.TUserRole
	for _, roleId := range in.RoleIds {
		m := &model.TUserRole{
			UserId: in.UserId,
			RoleId: roleId,
		}
		userRoles = append(userRoles, m)
	}

	// 添加用户角色
	_, err = l.svcCtx.TUserRoleModel.Inserts(l.ctx, userRoles...)
	if err != nil {
		return nil, err
	}

	return &permissionrpc.EmptyResp{}, nil
}
