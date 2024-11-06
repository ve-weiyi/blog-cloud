package album

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/photorpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAlbumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除相册
func NewDeleteAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAlbumLogic {
	return &DeleteAlbumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAlbumLogic) DeleteAlbum(req *types.IdReq) (resp *types.BatchResp, err error) {
	in := &photorpc.IdsReq{
		Ids: []int64{req.Id},
	}

	out, err := l.svcCtx.PhotoRpc.DeleteAlbum(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
