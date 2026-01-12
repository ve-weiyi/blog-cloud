package album

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/resourcerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAlbumDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新相册删除状态
func NewUpdateAlbumDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAlbumDeleteLogic {
	return &UpdateAlbumDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAlbumDeleteLogic) UpdateAlbumDelete(req *types.UpdateAlbumDeleteReq) (resp *types.BatchResp, err error) {
	in := &resourcerpc.UpdateAlbumDeleteReq{
		Ids:      req.Ids,
		IsDelete: req.IsDelete,
	}
	out, err := l.svcCtx.ResourceRpc.UpdateAlbumDelete(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
