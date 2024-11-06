package photo

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/photorpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeletePhotoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除照片
func NewBatchDeletePhotoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeletePhotoLogic {
	return &BatchDeletePhotoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeletePhotoLogic) BatchDeletePhoto(req *types.IdsReq) (resp *types.BatchResp, err error) {
	in := &photorpc.IdsReq{
		Ids: req.Ids,
	}

	out, err := l.svcCtx.PhotoRpc.DeletePhoto(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
