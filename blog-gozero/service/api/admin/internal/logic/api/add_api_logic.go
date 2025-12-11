package api

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/permissionrpc"

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

func (l *AddApiLogic) AddApi(req *types.ApiNewReq) (resp *types.ApiBackVO, err error) {
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

func convertApiTypes(req *permissionrpc.ApiDetailsResp) *types.ApiBackVO {

	children := make([]*types.ApiBackVO, 0)
	for _, v := range req.Children {
		m := convertApiTypes(v)
		children = append(children, m)
	}

	out := &types.ApiBackVO{
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
