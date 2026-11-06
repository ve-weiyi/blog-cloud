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

type CreateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建角色
func (l *CreateRoleLogic) CreateRole(in *accessrpc.CreateRoleRequest) (*accessrpc.CreateRoleResponse, error) {
	count, err := l.svcCtx.TRoleModel.FindCount(l.ctx, "role_key = ?", in.RoleKey)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, bizerr.NewBizError(bizcode.CodeResourceAlreadyExist, "角色标识已存在")
	}

	role := &model.TRole{
		ParentId:    in.ParentId,
		RoleKey:     in.RoleKey,
		RoleLabel:   in.RoleLabel,
		RoleComment: in.RoleComment,
		IsDefault:   in.IsDefault,
		Status:      in.Status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = l.svcCtx.TRoleModel.Insert(l.ctx, role)
	if err != nil {
		return nil, err
	}

	return &accessrpc.CreateRoleResponse{
		Id:   role.Id,
		Role: convertRoleOut(role),
	}, nil
}
