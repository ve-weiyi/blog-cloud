package permissionrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/jsonconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMenuLogic {
	return &AddMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建菜单
func (l *AddMenuLogic) AddMenu(in *permissionrpc.AddMenuReq) (*permissionrpc.AddMenuResp, error) {
	entity := convertMenuIn(in)

	_, err := l.svcCtx.TMenuModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &permissionrpc.AddMenuResp{
		Menu: convertMenuOut(entity),
	}, nil
}

func convertMenuIn(in *permissionrpc.AddMenuReq) (out *model.TMenu) {
	out = &model.TMenu{
		Id:         in.Id,
		ParentId:   in.ParentId,
		Path:       in.Path,
		Name:       in.Name,
		Component:  in.Component,
		Redirect:   in.Redirect,
		Type:       in.Meta.Type,
		Title:      in.Meta.Title,
		Icon:       in.Meta.Icon,
		Rank:       in.Meta.Rank,
		Perm:       in.Meta.Perm,
		Params:     in.Meta.Params,
		KeepAlive:  in.Meta.KeepAlive,
		AlwaysShow: in.Meta.AlwaysShow,
		Visible:    in.Meta.Visible,
		Status:     in.Meta.Status,
		Extra:      jsonconv.AnyToJsonNE(in.Meta),
	}

	return out
}

func convertMenuOut(in *model.TMenu) (out *permissionrpc.Menu) {

	out = &permissionrpc.Menu{
		Id:        in.Id,
		ParentId:  in.ParentId,
		Path:      in.Path,
		Name:      in.Name,
		Component: in.Component,
		Redirect:  in.Redirect,
		CreatedAt: in.CreatedAt.UnixMilli(),
		UpdatedAt: in.UpdatedAt.UnixMilli(),
		Children:  nil,
		Meta: &permissionrpc.MenuMeta{
			Type:       in.Type,
			Title:      in.Title,
			Icon:       in.Icon,
			Rank:       in.Rank,
			Perm:       in.Perm,
			Params:     in.Params,
			KeepAlive:  in.KeepAlive,
			AlwaysShow: in.AlwaysShow,
			Visible:    in.Visible,
			Status:     in.Status,
		},
	}
	return out
}
