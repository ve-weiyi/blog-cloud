package permissionrpclogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
)

type UpdateRoleApisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRoleApisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleApisLogic {
	return &UpdateRoleApisLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新角色资源
func (l *UpdateRoleApisLogic) UpdateRoleApis(in *permissionrpc.UpdateRoleApisReq) (*permissionrpc.EmptyResp, error) {
	// 删除
	_, err := l.svcCtx.TRoleApiModel.Deletes(l.ctx, "role_id in (?)", in.RoleId)
	if err != nil {
		return nil, err
	}

	var roleApis []*model.TRoleApi
	for _, apiId := range in.ApiIds {
		m := &model.TRoleApi{
			RoleId: in.RoleId,
			ApiId:  apiId,
		}

		roleApis = append(roleApis, m)
	}

	// 添加
	_, err = l.svcCtx.TRoleApiModel.Inserts(l.ctx, roleApis...)
	if err != nil {
		return nil, err
	}

	return &permissionrpc.EmptyResp{}, nil
}
