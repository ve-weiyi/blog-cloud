package api

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/permissionrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建api路由
func NewAddApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddApiLogic {
	return &AddApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddApiLogic) AddApi(req *types.ApiNewReq) (resp *types.ApiBackDTO, err error) {
	in := &permissionrpc.ApiNewReq{
		Id:        req.Id,
		ParentId:  req.ParentId,
		Path:      req.Path,
		Name:      req.Name,
		Method:    req.Method,
		Traceable: req.Traceable,
		IsDisable: req.IsDisable,
		Children:  nil,
	}

	out, err := l.svcCtx.PermissionRpc.AddApi(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return convertApiTypes(out), nil
}

func convertApiTypes(req *permissionrpc.ApiDetails) *types.ApiBackDTO {

	children := make([]*types.ApiBackDTO, 0)
	for _, v := range req.Children {
		m := convertApiTypes(v)
		children = append(children, m)
	}

	out := &types.ApiBackDTO{
		Id:        req.Id,
		ParentId:  req.ParentId,
		Name:      req.Name,
		Path:      req.Path,
		Method:    req.Method,
		Traceable: req.Traceable,
		IsDisable: req.IsDisable,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
		Children:  children,
	}

	return out
}
