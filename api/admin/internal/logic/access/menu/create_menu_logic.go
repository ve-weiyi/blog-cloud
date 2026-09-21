package menu

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/accessservice"
)

type CreateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建菜单
func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMenuLogic) CreateMenu(req *types.CreateMenuReq) (resp *types.MenuVO, err error) {
	out, err := l.svcCtx.AccessService.CreateMenu(l.ctx, &accessservice.CreateMenuRequest{
		ParentId:  req.ParentId,
		Path:      req.Path,
		Name:      req.Name,
		Component: req.Component,
		Redirect:  req.Redirect,
		Meta: &accessservice.MenuMeta{
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

	return &types.MenuVO{
		Id:        out.Id,
		ParentId:  req.ParentId,
		Path:      req.Path,
		Name:      req.Name,
		Component: req.Component,
		Redirect:  req.Redirect,
		MenuMeta:  req.MenuMeta,
	}, nil
}
