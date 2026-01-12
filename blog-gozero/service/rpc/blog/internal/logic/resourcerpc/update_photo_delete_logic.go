package resourcerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePhotoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePhotoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePhotoDeleteLogic {
	return &UpdatePhotoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新照片删除状态
func (l *UpdatePhotoDeleteLogic) UpdatePhotoDelete(in *resourcerpc.UpdatePhotoDeleteReq) (*resourcerpc.UpdatePhotoDeleteResp, error) {
	rows, err := l.svcCtx.TPhotoModel.Updates(l.ctx, map[string]interface{}{
		"is_delete": in.IsDelete,
	}, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &resourcerpc.UpdatePhotoDeleteResp{
		SuccessCount: rows,
	}, nil
}
