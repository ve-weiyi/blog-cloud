package permissionrpclogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/query"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
)

type FindApiListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindApiListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindApiListLogic {
	return &FindApiListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页获取接口列表
func (l *FindApiListLogic) FindApiList(in *permissionrpc.FindApiListReq) (*permissionrpc.FindApiListResp, error) {
	var opts []query.Option

	if in.Name != "" {
		opts = append(opts, query.WithCondition("name like ?", "%"+in.Name+"%"))
	}

	if in.Path != "" {
		opts = append(opts, query.WithCondition("path like ?", "%"+in.Path+"%"))
	}

	if in.Method != "" {
		opts = append(opts, query.WithCondition("method like ?", "%"+in.Method+"%"))
	}

	_, _, _, conditions, params := query.NewQueryBuilder(opts...).Build()
	result, err := l.svcCtx.TApiModel.FindALL(l.ctx, conditions, params...)
	if err != nil {
		return nil, err
	}

	out := &permissionrpc.FindApiListResp{}
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
			root := convertApiOut(item)
			root.Children = appendApiChildren(root, result)
			out.List = append(out.List, root)
		}
	}

	return out, nil
}

func appendApiChildren(root *permissionrpc.ApiDetailsResp, list []*model.TApi) (leafs []*permissionrpc.ApiDetailsResp) {
	for _, item := range list {
		if item.ParentId == root.Id {
			leaf := convertApiOut(item)
			leaf.Children = appendApiChildren(leaf, list)
			leafs = append(leafs, leaf)
		}
	}
	return leafs
}
