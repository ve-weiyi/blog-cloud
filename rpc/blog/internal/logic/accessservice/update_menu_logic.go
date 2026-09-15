package accessservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
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
func (l *UpdateMenuLogic) UpdateMenu(in *accessrpc.UpdateMenuRequest) (*accessrpc.UpdateMenuResponse, error) {
	data := &model.TMenu{
		Id:        in.Id,
		ParentId:  in.ParentId,
		Path:      in.Path,
		Name:      in.Name,
		Component: in.Component,
		Redirect:  in.Redirect,
	}
	if in.Meta != nil {
		data.Type = in.Meta.Type
		data.Title = in.Meta.Title
		data.Icon = in.Meta.Icon
		data.Rank = in.Meta.Rank
		data.Perm = in.Meta.Perm
		data.Params = in.Meta.Params
		data.KeepAlive = in.Meta.KeepAlive
		data.AlwaysShow = in.Meta.AlwaysShow
		data.Visible = in.Meta.Visible
		data.Status = in.Meta.Status
	}

	_, err := l.svcCtx.TMenuModel.Update(l.ctx, data)
	if err != nil {
		return nil, err
	}

	return &accessrpc.UpdateMenuResponse{Success: true}, nil
}
