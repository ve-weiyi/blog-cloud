package album

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/resourceservice"
)

type UpdateAlbumDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量更新相册删除状态
func NewUpdateAlbumDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAlbumDeleteLogic {
	return &UpdateAlbumDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAlbumDeleteLogic) UpdateAlbumDelete(req *types.UpdateAlbumDeleteReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.ResourceService.PatchAlbum(l.ctx, &resourceservice.PatchAlbumRequest{
		Ids:      req.Ids,
		IsDelete: req.IsDelete,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
