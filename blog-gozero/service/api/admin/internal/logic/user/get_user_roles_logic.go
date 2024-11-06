package user

import (
	"context"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/permissionrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户角色
func NewGetUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserRolesLogic {
	return &GetUserRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserRolesLogic) GetUserRoles(req *types.EmptyReq) (resp *types.UserRolesResp, err error) {
	in := &permissionrpc.UserIdReq{
		UserId: cast.ToString(l.ctx.Value("uid")),
	}

	out, err := l.svcCtx.PermissionRpc.FindUserRoles(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var list []*types.UserRole
	for _, v := range out.List {
		list = append(list, convertUserRole(v))
	}

	resp = &types.UserRolesResp{}
	resp.List = list
	return
}

func convertUserRole(in *permissionrpc.RoleDetails) (out *types.UserRole) {
	out = &types.UserRole{
		Id:          in.Id,
		ParentId:    in.ParentId,
		RoleName:    in.RoleName,
		RoleLabel:   in.RoleLabel,
		RoleComment: in.RoleComment,
	}

	return out
}
