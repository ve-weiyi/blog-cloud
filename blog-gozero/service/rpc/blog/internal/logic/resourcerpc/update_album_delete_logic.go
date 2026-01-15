package resourcerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAlbumDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAlbumDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAlbumDeleteLogic {
	return &UpdateAlbumDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新相册删除状态
func (l *UpdateAlbumDeleteLogic) UpdateAlbumDelete(in *resourcerpc.UpdateAlbumDeleteReq) (*resourcerpc.UpdateAlbumDeleteResp, error) {
	rows, err := l.svcCtx.TAlbumModel.Updates(l.ctx, map[string]interface{}{
		"is_delete": in.IsDelete,
	}, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.TPhotoModel.Updates(l.ctx, map[string]interface{}{
		"is_delete": in.IsDelete,
	}, "album_id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &resourcerpc.UpdateAlbumDeleteResp{
		SuccessCount: rows,
	}, nil
}
