package permissionrpclogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
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
func (l *SyncMenuListLogic) SyncMenuList(in *permissionrpc.SyncMenuListReq) (*permissionrpc.BatchResp, error) {
	go l.syncMenuList(in)

	return &permissionrpc.BatchResp{
		SuccessCount: 0,
	}, nil
}

func (l *SyncMenuListLogic) syncMenuList(in *permissionrpc.SyncMenuListReq) {
	// 使用后台上下文，服务返回后仍然可以继续执行
	ctx := context.Background()

	var err error
	for _, item := range in.Menus {
		err = l.InsertMenu(ctx, item)
		if err != nil {
			l.Errorf("插入数据失败: %v", err)
			return
		}
	}

	l.Infof("成功同步菜单")
	return
}

func (l *SyncMenuListLogic) InsertMenu(ctx context.Context, item *permissionrpc.MenuNewReq) (err error) {
	// 已存在则跳过
	parent, _ := l.svcCtx.TMenuModel.FindOneByPath(ctx, item.Path)
	if parent == nil {
		// 插入数据
		parent = convertMenuIn(item)
		_, err = l.svcCtx.TMenuModel.Insert(ctx, parent)
		if err != nil {
			return err
		}
	}

	for _, child := range item.Children {
		child.ParentId = parent.Id
		err = l.InsertMenu(ctx, child)
		if err != nil {
			return err
		}
	}

	return nil
}
