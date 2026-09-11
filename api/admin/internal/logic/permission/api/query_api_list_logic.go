package api

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/permissionservice"
)

type QueryApiListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取接口列表
func NewQueryApiListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryApiListLogic {
	return &QueryApiListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryApiListLogic) QueryApiList(req *types.QueryApiListReq) (resp *types.PageResult, err error) {
	out, err := l.svcCtx.PermissionService.ListApis(l.ctx, &permissionservice.ListApisRequest{
		Name:   req.Name,
		Path:   req.Path,
		Method: req.Method,
	})
	if err != nil {
		return nil, err
	}

	var apis []*types.ApiVO
	for _, v := range out.List {
		apis = append(apis, convertApiVO(v))
	}

	return &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     apis,
	}, nil
}

func convertApiVO(in *permissionservice.Api) *types.ApiVO {
	if in == nil {
		return nil
	}
	out := &types.ApiVO{
		Id:        in.Id,
		ParentId:  in.ParentId,
		Name:      in.Name,
		Path:      in.Path,
		Method:    in.Method,
		Traceable: in.Traceable,
		Status:    in.Status,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
	for _, child := range in.Children {
		out.Children = append(out.Children, convertApiVO(child))
	}
	return out
}
