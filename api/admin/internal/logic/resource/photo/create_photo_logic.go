package photo

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/resourceservice"
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
	out, err := l.svcCtx.ResourceService.CreatePhoto(l.ctx, &resourceservice.CreatePhotoRequest{
		AlbumId:   req.AlbumId,
		PhotoName: req.PhotoName,
		PhotoDesc: req.PhotoDesc,
		PhotoSrc:  req.PhotoSrc,
	})
	if err != nil {
		return nil, err
	}

	return &types.PhotoVO{
		Id:        out.Id,
		AlbumId:   req.AlbumId,
		PhotoName: req.PhotoName,
		PhotoDesc: req.PhotoDesc,
		PhotoSrc:  req.PhotoSrc,
	}, nil
}
