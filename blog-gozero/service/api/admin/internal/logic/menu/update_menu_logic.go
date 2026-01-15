package menu

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/jsonconv"
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

func (l *UpdateMenuLogic) UpdateMenu(req *types.NewMenuReq) (resp *types.MenuBackVO, err error) {
	in := &permissionrpc.UpdateMenuReq{
		Id:        req.Id,
		ParentId:  req.ParentId,
		Path:      req.Path,
		Name:      req.Name,
		Component: req.Component,
		Redirect:  req.Redirect,
		Meta: &permissionrpc.MenuMeta{
			Type:       req.Type,
			Title:      req.Title,
			Icon:       req.Icon,
			Rank:       req.Rank,
			Perm:       req.Perm,
			Params:     jsonconv.AnyToJsonNE(req.Params),
			KeepAlive:  req.KeepAlive,
			AlwaysShow: req.AlwaysShow,
			Visible:    req.Visible,
			Status:     req.Status,
		},
	}
	out, err := l.svcCtx.PermissionRpc.UpdateMenu(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = convertMenuTypes(out.Menu)
	return resp, nil
}
