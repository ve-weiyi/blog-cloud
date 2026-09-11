package photo

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/resourceservice"
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

func (l *QueryPhotoListLogic) QueryPhotoList(req *types.QueryPhotoListReq) (resp *types.PageResult, err error) {
	out, err := l.svcCtx.ResourceService.ListPhotos(l.ctx, &resourceservice.ListPhotosRequest{
		PageQuery: &resourceservice.PageQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
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

	return &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     list,
	}, nil
}
