package resourceservicelogic

import (
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/x/jsonconv"
)

func convertAlbumOut(in *model.TAlbum, photoCountMap map[int64]int) *resourcerpc.Album {
	var count int
	if photoCountMap != nil {
		count = photoCountMap[in.Id]
	}
	return &resourcerpc.Album{
		Id:         in.Id,
		AlbumName:  in.AlbumName,
		AlbumDesc:  in.AlbumDesc,
		AlbumCover: in.AlbumCover,
		IsDelete:   in.IsDelete,
		Status:     in.Status,
		CreatedAt:  in.CreatedAt.UnixMilli(),
		UpdatedAt:  in.UpdatedAt.UnixMilli(),
		PhotoCount: int64(count),
	}
}

func convertPhotoOut(in *model.TPhoto) *resourcerpc.Photo {
	return &resourcerpc.Photo{
		Id:        in.Id,
		AlbumId:   in.AlbumId,
		PhotoName: in.PhotoName,
		PhotoDesc: in.PhotoDesc,
		PhotoSrc:  in.PhotoSrc,
		IsDelete:  in.IsDelete,
		CreatedAt: in.CreatedAt.UnixMilli(),
		UpdatedAt: in.UpdatedAt.UnixMilli(),
	}
}

func convertPageOut(in *model.TPage) *resourcerpc.Page {
	out := &resourcerpc.Page{
		Id:         in.Id,
		PageName:   in.PageName,
		PageLabel:  in.PageLabel,
		PageCover:  in.PageCover,
		IsCarousel: in.IsCarousel,
		CreatedAt:  in.CreatedAt.UnixMilli(),
		UpdatedAt:  in.UpdatedAt.UnixMilli(),
	}
	jsonconv.JsonToAny(in.CarouselCovers, &out.CarouselCovers)
	return out
}
