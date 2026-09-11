package permissionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/adapter/gormx/queryx"
)

type ListApisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListApisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListApisLogic {
	return &ListApisLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListApisLogic) ListApis(in *permissionrpc.ListApisRequest) (*permissionrpc.ListApisResponse, error) {
	var opts []queryx.Option

	if in.Name != nil {
		opts = append(opts, queryx.WithCondition("name like ?", "%"+*in.Name+"%"))
	}

	if in.Path != nil {
		opts = append(opts, queryx.WithCondition("path like ?", "%"+*in.Path+"%"))
	}

	if in.Method != nil {
		opts = append(opts, queryx.WithCondition("method like ?", "%"+*in.Method+"%"))
	}

	_, _, _, conditions, params := queryx.NewQueryBuilder(opts...).Build()
	result, err := l.svcCtx.TApiModel.FindALL(l.ctx, conditions, params...)
	if err != nil {
		return nil, err
	}

	out := &permissionrpc.ListApisResponse{
		PageResult: &permissionrpc.PageResult{
			Page:     1,
			PageSize: 0,
			Total:    0,
		},
	}
	for _, item := range result {
		isParent := true
		for _, v := range result {
			if item.ParentId == v.Id {
				isParent = false
			}
		}

		if isParent {
			root := convertApiOut(item)
			root.Children = appendApiChildren(root, result)
			out.List = append(out.List, root)
		}
	}

	out.PageResult.PageSize = int64(len(out.List))
	out.PageResult.Total = int64(len(out.List))
	return out, nil
}

func appendApiChildren(root *permissionrpc.Api, list []*model.TApi) (leafs []*permissionrpc.Api) {
	for _, item := range list {
		if item.ParentId == root.Id {
			leaf := convertApiOut(item)
			leaf.Children = appendApiChildren(leaf, list)
			leafs = append(leafs, leaf)
		}
	}
	return leafs
}
