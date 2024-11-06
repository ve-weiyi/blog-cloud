package permissionrpclogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"
)

type SyncMenuListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncMenuListLogic {
	return &SyncMenuListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 同步菜单列表
func (l *SyncMenuListLogic) SyncMenuList(in *permissionrpc.SyncMenuReq) (*permissionrpc.BatchResp, error) {
	var err error
	var data int64
	for _, item := range in.Menus {
		// 已存在则跳过
		exist, _ := l.svcCtx.TMenuModel.First(l.ctx, "path = ?", item.Path)
		if exist == nil {

			// 插入数据
			exist = convertMenuIn(item)
			_, err = l.svcCtx.TMenuModel.Insert(l.ctx, exist)
			if err != nil {
				return nil, err
			}

			data++
		}

		for _, child := range item.Children {
			// 已存在则跳过
			menu, _ := l.svcCtx.TMenuModel.First(l.ctx, "path = ?", child.Path)
			if menu == nil {
				// 插入数据
				menu = convertMenuIn(child)
				menu.ParentId = exist.Id
				_, err = l.svcCtx.TMenuModel.Insert(l.ctx, menu)
				if err != nil {
					return nil, err
				}

				data++
			}
		}
	}

	return &permissionrpc.BatchResp{
		SuccessCount: data,
	}, nil
}
