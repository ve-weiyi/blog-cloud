package menu

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/accessservice"
)

type QueryMenuListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取菜单列表
func NewQueryMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryMenuListLogic {
	return &QueryMenuListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryMenuListLogic) QueryMenuList(req *types.QueryMenuListReq) (resp *types.ListResult, err error) {
	out, err := l.svcCtx.AccessService.ListMenus(l.ctx, &accessservice.ListMenusRequest{})
	if err != nil {
		return nil, err
	}

	menus := convertMenuVOs(out.List)
	return &types.ListResult{
		Page:     1,
		PageSize: int64(len(menus)),
		Total:    int64(len(menus)),
		List:     menus,
	}, nil
}

func convertMenuVOs(list []*accessservice.Menu) []*types.MenuVO {
	if list == nil {
		return nil
	}
	result := make([]*types.MenuVO, 0, len(list))
	for _, v := range list {
		result = append(result, convertMenuVO(v))
	}
	return result
}

func convertMenuVO(v *accessservice.Menu) *types.MenuVO {
	if v == nil {
		return nil
	}
	m := &types.MenuVO{
		Id:        v.Id,
		ParentId:  v.ParentId,
		Path:      v.Path,
		Name:      v.Name,
		Component: v.Component,
		Redirect:  v.Redirect,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
		Children:  convertMenuVOs(v.Children),
	}
	if v.Meta != nil {
		m.MenuMeta = types.MenuMeta{
			Type:       v.Meta.Type,
			Title:      v.Meta.Title,
			Icon:       v.Meta.Icon,
			Rank:       v.Meta.Rank,
			Perm:       v.Meta.Perm,
			KeepAlive:  v.Meta.KeepAlive,
			AlwaysShow: v.Meta.AlwaysShow,
			Visible:    v.Meta.Visible,
			Status:     v.Meta.Status,
		}
	}
	return m
}
