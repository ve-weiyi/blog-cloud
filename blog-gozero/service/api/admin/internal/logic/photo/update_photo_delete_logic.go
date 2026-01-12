package photo

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/resourcerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePhotoDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新照片删除状态
func NewUpdatePhotoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePhotoDeleteLogic {
	return &UpdatePhotoDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePhotoDeleteLogic) UpdatePhotoDelete(req *types.UpdatePhotoDeleteReq) (resp *types.BatchResp, err error) {
	in := &resourcerpc.UpdatePhotoDeleteReq{
		Ids:      req.Ids,
		IsDelete: req.IsDelete,
	}
	out, err := l.svcCtx.ResourceRpc.UpdatePhotoDelete(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{SuccessCount: out.SuccessCount}, nil
}
