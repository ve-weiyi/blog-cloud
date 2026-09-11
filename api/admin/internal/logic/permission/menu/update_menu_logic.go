package menu

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/permissionservice"
)

type UpdateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新菜单
func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMenuLogic) UpdateMenu(req *types.UpdateMenuReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.PermissionService.UpdateMenu(l.ctx, &permissionservice.UpdateMenuRequest{
		Id:        req.Id,
		ParentId:  req.ParentId,
		Path:      req.Path,
		Name:      req.Name,
		Component: req.Component,
		Redirect:  req.Redirect,
		Meta: &permissionservice.MenuMeta{
			Type:       req.Type,
			Title:      req.Title,
			Icon:       req.Icon,
			Rank:       req.Rank,
			Perm:       req.Perm,
			KeepAlive:  req.KeepAlive,
			AlwaysShow: req.AlwaysShow,
			Visible:    req.Visible,
			Status:     req.Status,
		},
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
