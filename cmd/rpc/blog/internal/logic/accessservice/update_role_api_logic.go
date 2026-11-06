package accessservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/model"
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
func (l *UpdateRoleApiLogic) UpdateRoleApi(in *accessrpc.UpdateRoleApiRequest) (*accessrpc.UpdateRoleApiResponse, error) {
	err := l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
		if _, err := l.svcCtx.TRoleApiModel.WithTx(tx).DeleteBatch(l.ctx, "role_id = ?", in.RoleId); err != nil {
			return err
		}

		if len(in.ApiIds) == 0 {
			return nil
		}

		batch := make([]*model.TRoleApi, 0, len(in.ApiIds))
		for _, id := range in.ApiIds {
			batch = append(batch, &model.TRoleApi{RoleId: in.RoleId, ApiId: id})
		}
		_, err := l.svcCtx.TRoleApiModel.WithTx(tx).InsertBatch(l.ctx, batch...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return &accessrpc.UpdateRoleApiResponse{Success: true}, nil
}
