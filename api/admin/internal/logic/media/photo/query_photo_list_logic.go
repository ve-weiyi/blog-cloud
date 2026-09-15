package photo

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/mediaservice"
)

type QueryPhotoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取照片列表
func NewQueryPhotoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryPhotoListLogic {
	return &QueryPhotoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryPhotoListLogic) QueryPhotoList(req *types.QueryPhotoListReq) (resp *types.ListResult, err error) {
	out, err := l.svcCtx.MediaService.ListPhotos(l.ctx, &mediaservice.ListPhotosRequest{
		ListQuery: &mediaservice.ListQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		AlbumId:   *req.AlbumId,
		IsDelete:  req.IsDelete,
	})
	if err != nil {
		return nil, err
	}

	var list []*types.PhotoVO
	for _, v := range out.List {
		list = append(list, &types.PhotoVO{
			Id:        v.Id,
			AlbumId:   v.AlbumId,
			PhotoName: v.PhotoName,
			PhotoDesc: v.PhotoDesc,
			PhotoSrc:  v.PhotoSrc,
			IsDelete:  v.IsDelete,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return &types.ListResult{
		Page:     out.ListResult.Page,
		PageSize: out.ListResult.PageSize,
		Total:    out.ListResult.Total,
		List:     list,
	}, nil
}
