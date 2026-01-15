package permissionrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindAllRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindAllRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindAllRoleLogic {
	return &FindAllRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询所有角色
func (l *FindAllRoleLogic) FindAllRole(in *permissionrpc.FindAllRoleReq) (*permissionrpc.FindAllRoleResp, error) {
	result, err := l.svcCtx.TRoleModel.FindALL(l.ctx, "")
	if err != nil {
		return nil, err
	}

	out := &permissionrpc.FindAllRoleResp{}
	for _, item := range result {
		out.List = append(out.List, convertRoleOut(item))
	}

	return out, nil
}
