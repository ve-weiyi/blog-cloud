package mediaservicelogic

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/mediarpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
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
func (l *UpdatePhotoLogic) UpdatePhoto(in *mediarpc.UpdatePhotoRequest) (*mediarpc.UpdatePhotoResponse, error) {
	if _, err := l.svcCtx.TAlbumModel.FindOne(l.ctx, "", "`id` = ? AND `is_delete` = ?", in.AlbumId, 0); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerr.NewBizError(bizcode.CodeResourceNotFound, "相册不存在")
		}
		return nil, err
	}

	entity, err := l.svcCtx.TPhotoModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	entity.AlbumId = in.AlbumId
	entity.PhotoName = in.PhotoName
	entity.PhotoDesc = in.PhotoDesc
	entity.PhotoSrc = in.PhotoSrc

	_, err = l.svcCtx.TPhotoModel.Update(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &mediarpc.UpdatePhotoResponse{
		Photo: convertPhotoOut(entity),
	}, nil
}
