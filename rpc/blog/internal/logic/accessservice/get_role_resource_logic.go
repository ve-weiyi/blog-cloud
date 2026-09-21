package accessservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type GetRoleResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleResourceLogic {
	return &GetRoleResourceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询角色资源
func (l *GetRoleResourceLogic) GetRoleResource(in *accessrpc.GetRoleResourceRequest) (*accessrpc.GetRoleResourceResponse, error) {
	apiLinks, err := l.svcCtx.TRoleApiModel.FindALL(l.ctx, "role_id = ?", in.RoleId)
	if err != nil {
		return nil, err
	}
	menuLinks, err := l.svcCtx.TRoleMenuModel.FindALL(l.ctx, "role_id = ?", in.RoleId)
	if err != nil {
		return nil, err
	}

	out := &accessrpc.GetRoleResourceResponse{}
	for _, v := range apiLinks {
		out.ApiIds = append(out.ApiIds, v.ApiId)
	}
	for _, v := range menuLinks {
		out.MenuIds = append(out.MenuIds, v.MenuId)
	}

	return out, nil
}
