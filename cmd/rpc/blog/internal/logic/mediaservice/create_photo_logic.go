package mediaservicelogic

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/mediarpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/model"

	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
)

type CreatePhotoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePhotoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePhotoLogic {
	return &CreatePhotoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建照片
func (l *CreatePhotoLogic) CreatePhoto(in *mediarpc.CreatePhotoRequest) (*mediarpc.CreatePhotoResponse, error) {
	if _, err := l.svcCtx.TAlbumModel.FindOne(l.ctx, "", "`id` = ? AND `is_delete` = ?", in.AlbumId, 0); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerr.NewBizError(bizcode.CodeResourceNotFound, "相册不存在")
		}
		return nil, err
	}

	entity := &model.TPhoto{
		AlbumId:   in.AlbumId,
		PhotoName: in.PhotoName,
		PhotoDesc: in.PhotoDesc,
		PhotoSrc:  in.PhotoSrc,
	}

	_, err := l.svcCtx.TPhotoModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &mediarpc.CreatePhotoResponse{Id: entity.Id, Photo: convertPhotoOut(entity)}, nil
}
