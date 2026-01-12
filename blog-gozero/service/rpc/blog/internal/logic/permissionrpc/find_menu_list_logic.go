package permissionrpclogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/query"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
)

type FindMenuListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindMenuListLogic {
	return &FindMenuListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页获取菜单列表
func (l *FindMenuListLogic) FindMenuList(in *permissionrpc.FindMenuListReq) (*permissionrpc.FindMenuListResp, error) {
	var opts []query.Option

	if in.Name != "" {
		opts = append(opts, query.WithCondition("name like ?", "%"+in.Name+"%"))
	}

	if in.Title != "" {
		opts = append(opts, query.WithCondition("title like ?", "%"+in.Title+"%"))
	}

	_, _, _, conditions, params := query.NewQueryBuilder(opts...).Build()

	result, err := l.svcCtx.TMenuModel.FindALL(l.ctx, conditions, params...)
	if err != nil {
		return nil, err
	}

	out := &permissionrpc.FindMenuListResp{}
	for _, item := range result {
		// parentId不在当前菜单id列表，说明为父级菜单id，根据此id作为递归的开始条件节点
		isParent := true
		for _, v := range result {
			if item.ParentId == v.Id {
				isParent = false
			}
		}

		// parentId为0，说明为父级菜单
		if isParent {
			root := convertMenuOut(item)
			root.Children = appendMenuChildren(root, result)
			out.List = append(out.List, root)
		}
	}
	return out, nil
}

func appendMenuChildren(root *permissionrpc.Menu, list []*model.TMenu) (leafs []*permissionrpc.Menu) {
	for _, item := range list {
		if item.ParentId == root.Id {
			leaf := convertMenuOut(item)
			leaf.Children = appendMenuChildren(leaf, list)
			leafs = append(leafs, leaf)
		}
	}
	return leafs
}
