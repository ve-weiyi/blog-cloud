package permissionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

type UpdateRoleApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRoleApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleApiLogic {
	return &UpdateRoleApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新角色 API 资源
func (l *UpdateRoleApiLogic) UpdateRoleApi(in *permissionrpc.UpdateRoleApiRequest) (*permissionrpc.UpdateRoleApiResponse, error) {
	_, err := l.svcCtx.TRoleApiModel.DeleteBatch(l.ctx, "role_id = ?", in.RoleId)
	if err != nil {
		return nil, err
	}

	if len(in.ApiIds) > 0 {
		var batch []*model.TRoleApi
		for _, id := range in.ApiIds {
			batch = append(batch, &model.TRoleApi{
				RoleId: in.RoleId,
				ApiId:  id,
			})
		}
		_, err = l.svcCtx.TRoleApiModel.InsertBatch(l.ctx, batch...)
		if err != nil {
			return nil, err
		}
	}

	return &permissionrpc.UpdateRoleApiResponse{Success: true}, nil
}
