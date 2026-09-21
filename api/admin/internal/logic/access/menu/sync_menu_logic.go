package menu

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/accessservice"
)

type SyncMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 同步菜单列表
func NewSyncMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncMenuLogic {
	return &SyncMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncMenuLogic) SyncMenu(req *types.SyncMenuReq) (resp *types.SyncMenuResp, err error) {
	menus := convertCreateMenuReqs(req.Menus)
	out, err := l.svcCtx.AccessService.SyncMenus(l.ctx, &accessservice.SyncMenusRequest{
		Menus: menus,
	})
	if err != nil {
		return nil, err
	}

	return &types.SyncMenuResp{
		SuccessCount: out.SuccessCount,
	}, nil
}

func convertCreateMenuReqs(items []*types.CreateMenuReq) []*accessservice.CreateMenuRequest {
	if items == nil {
		return nil
	}
	result := make([]*accessservice.CreateMenuRequest, 0, len(items))
	for _, item := range items {
		result = append(result, convertCreateMenuReq(item))
	}
	return result
}

func convertCreateMenuReq(in *types.CreateMenuReq) *accessservice.CreateMenuRequest {
	if in == nil {
		return nil
	}

	paramsJSON := "[]"
	if len(in.Params) > 0 {
		if b, err := json.Marshal(in.Params); err == nil {
			paramsJSON = string(b)
		}
	}

	return &accessservice.CreateMenuRequest{
		ParentId:  in.ParentId,
		Path:      in.Path,
		Name:      in.Name,
		Component: in.Component,
		Redirect:  in.Redirect,
		Meta: &accessservice.MenuMeta{
			Type:       in.Type,
			Title:      in.Title,
			Icon:       in.Icon,
			Rank:       in.Rank,
			Perm:       in.Perm,
			Params:     paramsJSON,
			KeepAlive:  in.KeepAlive,
			AlwaysShow: in.AlwaysShow,
			Visible:    in.Visible,
			Status:     in.Status,
		},
		Children: convertCreateMenuReqs(in.Children),
	}
}
