package photo

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/mediaservice"
)

type CreatePhotoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建照片
func NewCreatePhotoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePhotoLogic {
	return &CreatePhotoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePhotoLogic) CreatePhoto(req *types.CreatePhotoReq) (resp *types.PhotoVO, err error) {
	out, err := l.svcCtx.MediaService.CreatePhoto(l.ctx, &mediaservice.CreatePhotoRequest{
		AlbumId:   req.AlbumId,
		PhotoName: req.PhotoName,
		PhotoDesc: req.PhotoDesc,
		PhotoSrc:  req.PhotoSrc,
	})
	if err != nil {
		return nil, err
	}

	return &types.PhotoVO{
		Id:        out.Photo.Id,
		AlbumId:   out.Photo.AlbumId,
		PhotoName: out.Photo.PhotoName,
		PhotoDesc: out.Photo.PhotoDesc,
		PhotoSrc:  out.Photo.PhotoSrc,
		CreatedAt: out.Photo.CreatedAt,
		UpdatedAt: out.Photo.UpdatedAt,
	}, nil
}
