package photo

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/resourceservice"
)

type UpdatePhotoDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量更新照片删除状态
func NewUpdatePhotoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePhotoDeleteLogic {
	return &UpdatePhotoDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePhotoDeleteLogic) UpdatePhotoDelete(req *types.UpdatePhotoDeleteReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.ResourceService.PatchPhoto(l.ctx, &resourceservice.PatchPhotoRequest{
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
