package permissionrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

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
func (l *SyncApiListLogic) SyncApiList(in *permissionrpc.SyncApiReq) (*permissionrpc.BatchResp, error) {
	var err error
	var data int64
	for _, item := range in.Apis {
		// 已存在则跳过
		exist, _ := l.svcCtx.TApiModel.First(l.ctx, "path = ?", item.Path)
		if exist == nil {

			// 插入数据
			exist = convertApiIn(item)
			_, err = l.svcCtx.TApiModel.Insert(l.ctx, exist)
			if err != nil {
				return nil, err
			}

			data++
		}

		for _, child := range item.Children {
			// 已存在则跳过
			menu, _ := l.svcCtx.TApiModel.First(l.ctx, "path = ?", child.Path)
			if menu == nil {
				// 插入数据
				menu = convertApiIn(child)
				menu.ParentId = exist.Id
				_, err = l.svcCtx.TApiModel.Insert(l.ctx, menu)
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
