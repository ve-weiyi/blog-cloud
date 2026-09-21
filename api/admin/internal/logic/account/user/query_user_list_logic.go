package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/accessservice"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/userservice"
)

type QueryUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户列表
func NewQueryUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryUserListLogic {
	return &QueryUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryUserListLogic) QueryUserList(req *types.QueryUserListReq) (resp *types.ListResult, err error) {
	out, err := l.svcCtx.UserService.ListUsers(l.ctx, &userservice.ListUsersRequest{
		ListQuery: &userservice.ListQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
	})
	if err != nil {
		return nil, err
	}

	var userIds []string
	for _, v := range out.List {
		userIds = append(userIds, v.UserId)
	}

	roleMap := make(map[string][]*types.UserRoleLabel)
	if len(userIds) > 0 {
		rolesResp, err := l.svcCtx.AccessService.BatchGetUserRoles(l.ctx, &accessservice.BatchGetUserRolesRequest{
			UserIds: userIds,
		})
		if err != nil {
			return nil, err
		}
		for _, ur := range rolesResp.List {
			var labels []*types.UserRoleLabel
			for _, r := range ur.List {
				labels = append(labels, &types.UserRoleLabel{
					RoleId:    r.Id,
					RoleKey:   r.RoleKey,
					RoleLabel: r.RoleLabel,
				})
			}
			roleMap[ur.UserId] = labels
		}
	}

	var list []*types.UserVO
	for _, v := range out.List {
		list = append(list, &types.UserVO{
			Id:           v.Id,
			UserId:       v.UserId,
			Username:     v.Username,
			Nickname:     v.Nickname,
			Avatar:       v.Avatar,
			Mobile:       v.Mobile,
			Email:        v.Email,
			Status:       v.Status,
			RegisterType: v.RegisterType,
			IpAddress:    v.IpAddress,
			IpSource:     v.IpSource,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
			RoleLabels:   roleMap[v.UserId],
		})
	}

	return &types.ListResult{
		Page:     out.ListResult.Page,
		PageSize: out.ListResult.PageSize,
		Total:    out.ListResult.Total,
		List:     list,
	}, nil
}
