package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/permissionservice"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/client/userservice"
)

type GetUserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户详情
func NewGetUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserDetailLogic {
	return &GetUserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserDetailLogic) GetUserDetail(req *types.GetUserDetailReq) (resp *types.GetUserDetailResp, err error) {
	out, err := l.svcCtx.UserService.GetUser(l.ctx, &userservice.GetUserRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	var roleLabels []*types.UserRoleLabel
	rolesResp, err := l.svcCtx.PermissionService.GetUsersRoles(l.ctx, &permissionservice.GetUsersRolesRequest{
		UserIds: []string{req.UserId},
	})
	if err != nil {
		return nil, err
	}
	for _, ur := range rolesResp.List {
		for _, r := range ur.List {
			roleLabels = append(roleLabels, &types.UserRoleLabel{
				RoleId:    r.Id,
				RoleKey:   r.RoleKey,
				RoleLabel: r.RoleLabel,
			})
		}
	}

	return &types.GetUserDetailResp{
		UserVO: types.UserVO{
			Id:           out.User.Id,
			UserId:       out.User.UserId,
			Username:     out.User.Username,
			Nickname:     out.User.Nickname,
			Avatar:       out.User.Avatar,
			Mobile:       out.User.Mobile,
			Email:        out.User.Email,
			Status:       out.User.Status,
			RegisterType: out.User.RegisterType,
			IpAddress:    out.User.IpAddress,
			IpSource:     out.User.IpSource,
			CreatedAt:    out.User.CreatedAt,
			UpdatedAt:    out.User.UpdatedAt,
			RoleLabels:   roleLabels,
		},
	}, nil
}
