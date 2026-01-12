package permissionrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncApiListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncApiListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncApiListLogic {
	return &SyncApiListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 同步接口列表
func (l *SyncApiListLogic) SyncApiList(in *permissionrpc.SyncApiListReq) (*permissionrpc.SyncApiListResp, error) {
	go l.syncApiList(in)

	return &permissionrpc.SyncApiListResp{
		SuccessCount: 0,
	}, nil
}

func (l *SyncApiListLogic) syncApiList(in *permissionrpc.SyncApiListReq) {
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

func (l *SyncApiListLogic) InsertApi(ctx context.Context, item *permissionrpc.AddApiReq) (err error) {
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
