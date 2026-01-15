package resourcerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePhotoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePhotoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePhotoLogic {
	return &UpdatePhotoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新照片
func (l *UpdatePhotoLogic) UpdatePhoto(in *resourcerpc.UpdatePhotoReq) (*resourcerpc.UpdatePhotoResp, error) {
	entity, err := l.svcCtx.TPhotoModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	entity.AlbumId = in.AlbumId
	entity.PhotoName = in.PhotoName
	entity.PhotoDesc = in.PhotoDesc
	entity.PhotoSrc = in.PhotoSrc
	entity.IsDelete = in.IsDelete

	_, err = l.svcCtx.TPhotoModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &resourcerpc.UpdatePhotoResp{
		Photo: convertPhotoOut(entity),
	}, nil
}
