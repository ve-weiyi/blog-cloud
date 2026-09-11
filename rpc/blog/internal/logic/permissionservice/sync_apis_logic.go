package permissionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type SyncApisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncApisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncApisLogic {
	return &SyncApisLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 同步接口列表
func (l *SyncApisLogic) SyncApis(in *permissionrpc.SyncApisRequest) (*permissionrpc.SyncApisResponse, error) {
	go l.syncApiList(in)

	return &permissionrpc.SyncApisResponse{
		SuccessCount: 0,
	}, nil
}

func (l *SyncApisLogic) syncApiList(in *permissionrpc.SyncApisRequest) {
	// 使用后台上下文，服务返回后仍然可以继续执行
	ctx := context.Background()

	var err error
	for _, item := range in.Apis {
		err = l.InsertApi(ctx, item)
		if err != nil {
			l.Errorf("插入数据失败: %v", err)
			return
		}
	}

	l.Infof("成功同步接口")
	return
}

func (l *SyncApisLogic) InsertApi(ctx context.Context, item *permissionrpc.CreateApiRequest) (err error) {
	// 已存在则跳过
	parent, _ := l.svcCtx.TApiModel.FindOneByPathMethodName(ctx, item.Path, item.Method, item.Name)
	if parent == nil {
		// 插入数据
		parent = convertApiIn(item)
		_, err = l.svcCtx.TApiModel.Insert(ctx, parent)
		if err != nil {
			return err
		}
	}

	for _, child := range item.Children {
		child.ParentId = parent.Id
		err = l.InsertApi(ctx, child)
		if err != nil {
			return err
		}
	}

	return nil
}
