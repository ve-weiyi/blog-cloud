package accessservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type SyncMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncMenusLogic {
	return &SyncMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 同步菜单列表
func (l *SyncMenusLogic) SyncMenus(in *accessrpc.SyncMenusRequest) (*accessrpc.SyncMenusResponse, error) {
	go l.syncMenuList(in)

	return &accessrpc.SyncMenusResponse{
		SuccessCount: 0,
	}, nil
}

func (l *SyncMenusLogic) syncMenuList(in *accessrpc.SyncMenusRequest) {
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

func (l *SyncMenusLogic) InsertMenu(ctx context.Context, item *accessrpc.CreateMenuRequest) (err error) {
	// 已存在则跳过
	parent, _ := l.svcCtx.TMenuModel.FindOneByPathPerm(ctx, item.Path, item.Meta.Perm)
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
