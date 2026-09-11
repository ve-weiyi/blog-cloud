package resourceservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type PatchPhotoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPatchPhotoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PatchPhotoLogic {
	return &PatchPhotoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 部分更新照片
func (l *PatchPhotoLogic) PatchPhoto(in *resourcerpc.PatchPhotoRequest) (*resourcerpc.PatchPhotoResponse, error) {
	rows, err := l.svcCtx.TPhotoModel.UpdateFields(l.ctx, map[string]interface{}{
		"is_delete": in.IsDelete,
	}, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &resourcerpc.PatchPhotoResponse{SuccessCount: rows}, nil
}
