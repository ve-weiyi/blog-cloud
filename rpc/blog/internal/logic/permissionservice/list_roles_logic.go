package permissionservicelogic

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type ListRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRolesLogic {
	return &ListRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListRolesLogic) ListRoles(in *permissionrpc.ListRolesRequest) (*permissionrpc.ListRolesResponse, error) {
	var conditions string
	var args []interface{}
	if in.RoleKey != nil {
		conditions = appendCondRole(conditions, "`role_key` like ?")
		args = append(args, "%"+*in.RoleKey+"%")
	}
	if in.RoleLabel != nil {
		conditions = appendCondRole(conditions, "`role_label` like ?")
		args = append(args, "%"+*in.RoleLabel+"%")
	}
	if in.Status != nil {
		conditions = appendCondRole(conditions, "`status` = ?")
		args = append(args, in.Status)
	}

	page := int(in.PageQuery.GetPage())
	size := int(in.PageQuery.GetPageSize())

	if size <= 0 {
		list, err := l.svcCtx.TRoleModel.FindALL(l.ctx, conditions, args...)
		if err != nil {
			return nil, err
		}
		out := &permissionrpc.ListRolesResponse{
			PageResult: &permissionrpc.PageResult{
				Page:     1,
				PageSize: int64(len(list)),
				Total:    int64(len(list)),
			},
		}
		for _, r := range list {
			out.List = append(out.List, convertRoleOut(r))
		}
		return out, nil
	}

	sorts := strings.Join(in.PageQuery.GetSorts(), ",")
	list, total, err := l.svcCtx.TRoleModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, args...)
	if err != nil {
		return nil, err
	}

	out := &permissionrpc.ListRolesResponse{
		PageResult: &permissionrpc.PageResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
	}
	for _, r := range list {
		out.List = append(out.List, convertRoleOut(r))
	}

	return out, nil
}

func appendCondRole(base string, next string) string {
	if len(base) == 0 {
		return next
	}
	return base + " and " + next
}
