package permissionrpclogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"
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
	var (
		conditions string
		params     []interface{}
	)

	if in.Name != "" {
		conditions += "name like ?"
		params = append(params, "%"+in.Name+"%")
	}

	if in.Path != "" {
		if conditions != "" {
			conditions += " and "
		}
		conditions += "path like ?"
		params = append(params, "%"+in.Path+"%")
	}

	if in.Method != "" {
		if conditions != "" {
			conditions += " and "
		}
		conditions += "method like ?"
		params = append(params, "%"+in.Method+"%")
	}

	result, err := l.svcCtx.TApiModel.FindALL(l.ctx, conditions, params...)
	if err != nil {
		return nil, err
	}

	out := &permissionrpc.FindApiListResp{}
	if conditions == "" {
		var root permissionrpc.ApiDetails
		root.Children = appendApiChildren(&root, result)
		out.List = root.Children
	} else {
		for _, item := range result {
			out.List = append(out.List, convertApiOut(item))
		}
	}

	return out, nil
}

func appendApiChildren(root *permissionrpc.ApiDetails, list []*model.TApi) (leafs []*permissionrpc.ApiDetails) {
	for _, item := range list {
		if item.ParentId == root.Id {
			leaf := convertApiOut(item)
			leaf.Children = appendApiChildren(leaf, list)
			leafs = append(leafs, leaf)
		}
	}
	return leafs
}
