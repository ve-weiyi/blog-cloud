package accessservicelogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/model"
)

type UpdateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新角色
func (l *UpdateRoleLogic) UpdateRole(in *accessrpc.UpdateRoleRequest) (*accessrpc.UpdateRoleResponse, error) {
	count, err := l.svcCtx.TRoleModel.FindCount(l.ctx, "role_key = ? AND id != ?", in.RoleKey, in.Id)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, bizerr.NewBizError(bizcode.CodeResourceAlreadyExist, "角色标识已存在")
	}

	role := &model.TRole{
		Id:          in.Id,
		ParentId:    in.ParentId,
		RoleKey:     in.RoleKey,
		RoleLabel:   in.RoleLabel,
		RoleComment: in.RoleComment,
		IsDefault:   in.IsDefault,
		Status:      in.Status,
		UpdatedAt:   time.Now(),
	}

	_, err = l.svcCtx.TRoleModel.Update(l.ctx, role)
	if err != nil {
		return nil, err
	}

	// 更新后回读：返回体需要包含未在请求中的字段（创建时间等）
	entity, err := l.svcCtx.TRoleModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &accessrpc.UpdateRoleResponse{Success: true, Role: convertRoleOut(entity)}, nil
}
