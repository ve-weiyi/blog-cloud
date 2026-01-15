package permissionrpclogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/query"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
)

type FindRoleListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindRoleListLogic {
	return &FindRoleListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页获取角色列表
func (l *FindRoleListLogic) FindRoleList(in *permissionrpc.FindRoleListReq) (*permissionrpc.FindRoleListResp, error) {
	var opts []query.Option
	if in.Paginate != nil {
		opts = append(opts, query.WithPage(int(in.Paginate.Page)))
		opts = append(opts, query.WithSize(int(in.Paginate.PageSize)))
		opts = append(opts, query.WithSorts(in.Paginate.Sorts...))
	}
	page, size, sorts, conditions, params := query.NewQueryBuilder(opts...).Build()

	result, _, err := l.svcCtx.TRoleModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var root permissionrpc.Role
	root.Children = appendRoleChildren(&root, result)

	return &permissionrpc.FindRoleListResp{
		List: root.Children,
		Pagination: &permissionrpc.PageResp{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    int64(len(root.Children)),
		},
	}, nil
}

func appendRoleChildren(root *permissionrpc.Role, list []*model.TRole) (leafs []*permissionrpc.Role) {
	for _, item := range list {
		if item.ParentId == root.Id {
			leaf := convertRoleOut(item)
			leaf.Children = appendRoleChildren(leaf, list)
			leafs = append(leafs, leaf)
		}
	}
	return leafs
}
