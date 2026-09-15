package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/accessservice"
)

type BatchDeleteRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除角色
func NewBatchDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteRoleLogic {
	return &BatchDeleteRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteRoleLogic) BatchDeleteRole(req *types.DeleteRoleReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.AccessService.BatchDeleteRoles(l.ctx, &accessservice.BatchDeleteRolesRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
