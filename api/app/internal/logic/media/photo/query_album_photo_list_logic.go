package photo

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/mediaservice"
)

type QueryAlbumPhotoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取相册下的照片列表
func NewQueryAlbumPhotoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryAlbumPhotoListLogic {
	return &QueryAlbumPhotoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryAlbumPhotoListLogic) QueryAlbumPhotoList(req *types.QueryAlbumPhotoListReq) (resp *types.ListResult, err error) {
	isDelete := int64(0)

	in := &mediaservice.ListPhotosRequest{
		ListQuery: &mediaservice.ListQuery{
			Page:     1,
			PageSize: 100,
		},
		AlbumId:  req.AlbumId,
		IsDelete: &isDelete,
	}

	out, err := l.svcCtx.MediaService.ListPhotos(l.ctx, in)
	if err != nil {
		return nil, err
	}

	list := make([]*types.Photo, 0)
	for _, v := range out.List {
		list = append(list, &types.Photo{
			Id:       v.Id,
			PhotoUrl: v.PhotoSrc,
		})
	}

	resp = &types.ListResult{
		Page:     out.ListResult.Page,
		PageSize: out.ListResult.PageSize,
		Total:    out.ListResult.Total,
		List:     list,
	}
	return
}
