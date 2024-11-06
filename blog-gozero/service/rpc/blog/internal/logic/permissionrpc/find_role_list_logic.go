package permissionrpclogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"

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
	var (
		page       int
		size       int
		sorts      string
		conditions string
		params     []interface{}
	)

	page = int(in.Page)
	size = int(in.PageSize)

	result, err := l.svcCtx.TRoleModel.FindList(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var root permissionrpc.RoleDetails
	root.Children = appendRoleChildren(&root, result)

	out := &permissionrpc.FindRoleListResp{}
	out.Total = int64(len(root.Children))
	out.List = root.Children

	return out, nil
}

func appendRoleChildren(root *permissionrpc.RoleDetails, list []*model.TRole) (leafs []*permissionrpc.RoleDetails) {
	for _, item := range list {
		if item.ParentId == root.Id {
			leaf := convertRoleOut(item)
			leaf.Children = appendRoleChildren(leaf, list)
			leafs = append(leafs, leaf)
		}
	}
	return leafs
}
