package permissionrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新菜单
func (l *UpdateMenuLogic) UpdateMenu(in *permissionrpc.UpdateMenuReq) (*permissionrpc.UpdateMenuResp, error) {
	entity, err := l.svcCtx.TMenuModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	entity.ParentId = in.ParentId
	entity.Path = in.Path
	entity.Name = in.Name
	entity.Component = in.Component
	entity.Redirect = in.Redirect
	// 更新 meta 信息
	if in.Meta != nil {
		entity.Type = in.Meta.Type
		entity.Title = in.Meta.Title
		entity.Icon = in.Meta.Icon
		entity.Rank = in.Meta.Rank
		entity.Perm = in.Meta.Perm
		entity.Params = in.Meta.Params
		entity.KeepAlive = in.Meta.KeepAlive
		entity.AlwaysShow = in.Meta.AlwaysShow
		entity.Visible = in.Meta.Visible
		entity.Status = in.Meta.Status
	}

	_, err = l.svcCtx.TMenuModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &permissionrpc.UpdateMenuResp{
		Menu: convertMenuOut(entity),
	}, nil
}
