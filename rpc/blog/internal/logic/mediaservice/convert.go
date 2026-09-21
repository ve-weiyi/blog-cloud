package mediaservicelogic

import (
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/mediarpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

func convertAlbumOut(in *model.TAlbum, photoCountMap map[int64]int) *mediarpc.Album {
	var count int
	if photoCountMap != nil {
		count = photoCountMap[in.Id]
	}
	return &mediarpc.Album{
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

func convertPhotoOut(in *model.TPhoto) *mediarpc.Photo {
	return &mediarpc.Photo{
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
