package me

import (
	"context"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizheader"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/accessservice"
)

type GetMeRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户角色
func NewGetMeRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMeRolesLogic {
	return &GetMeRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMeRolesLogic) GetMeRoles(req *types.EmptyReq) (resp *types.GetMeRolesResp, err error) {
	userId := cast.ToString(l.ctx.Value(bizheader.HeaderUid))

	out, err := l.svcCtx.AccessService.ListUserRoles(l.ctx, &accessservice.ListUserRolesRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	var list []*types.UserRole
	for _, v := range out.List {
		list = append(list, &types.UserRole{
			Id:          v.Id,
			ParentId:    v.ParentId,
			RoleKey:     v.RoleKey,
			RoleLabel:   v.RoleLabel,
			RoleComment: v.RoleComment,
		})
	}

	return &types.GetMeRolesResp{List: list}, nil
}
