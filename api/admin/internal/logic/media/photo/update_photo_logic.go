package photo

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/mediaservice"
)

type UpdatePhotoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新照片
func NewUpdatePhotoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePhotoLogic {
	return &UpdatePhotoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePhotoLogic) UpdatePhoto(req *types.UpdatePhotoReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.MediaService.UpdatePhoto(l.ctx, &mediaservice.UpdatePhotoRequest{
		Id:        req.Id,
		AlbumId:   req.AlbumId,
		PhotoName: req.PhotoName,
		PhotoSrc:  req.PhotoSrc,
		PhotoDesc: req.PhotoDesc,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
