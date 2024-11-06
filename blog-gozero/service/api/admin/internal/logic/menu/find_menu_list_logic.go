package menu

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/permissionrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindMenuListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页获取菜单列表
func NewFindMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindMenuListLogic {
	return &FindMenuListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindMenuListLogic) FindMenuList(req *types.MenuQuery) (resp *types.PageResp, err error) {
	in := &permissionrpc.FindMenuListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
		Title:    req.Title,
	}

	out, err := l.svcCtx.PermissionRpc.FindMenuList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var list []*types.MenuBackDTO
	for _, v := range out.List {
		m := ConvertMenuTypes(v)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = in.Page
	resp.PageSize = in.PageSize
	resp.Total = int64(len(list))
	resp.List = list
	return resp, nil
}
