package accessservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/queryx"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

type ListMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMenusLogic {
	return &ListMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListMenusLogic) ListMenus(in *accessrpc.ListMenusRequest) (*accessrpc.ListMenusResponse, error) {
	var opts []queryx.Option

	if in.Name != nil {
		opts = append(opts, queryx.WithCondition("name like ?", "%"+*in.Name+"%"))
	}

	if in.Title != nil {
		opts = append(opts, queryx.WithCondition("title like ?", "%"+*in.Title+"%"))
	}

	_, _, _, conditions, params := queryx.NewQueryBuilder(opts...).Build()

	result, err := l.svcCtx.TMenuModel.FindALL(l.ctx, conditions, params...)
	if err != nil {
		return nil, err
	}

	out := &accessrpc.ListMenusResponse{
		ListResult: &accessrpc.ListResult{
			Page:     1,
			PageSize: 0,
			Total:    0,
		},
		List: []*accessrpc.Menu{},
	}
	for _, item := range result {
		isParent := true
		for _, v := range result {
			if item.ParentId == v.Id {
				isParent = false
			}
		}

		if isParent {
			root := convertMenuOut(item)
			root.Children = appendMenuChildren(root, result)
			out.List = append(out.List, root)
		}
	}
	out.ListResult.PageSize = int64(len(out.List))
	out.ListResult.Total = int64(len(out.List))
	return out, nil
}

func appendMenuChildren(root *accessrpc.Menu, list []*model.TMenu) (leafs []*accessrpc.Menu) {
	for _, item := range list {
		if item.ParentId == root.Id {
			leaf := convertMenuOut(item)
			leaf.Children = appendMenuChildren(leaf, list)
			leafs = append(leafs, leaf)
		}
	}
	return leafs
}
