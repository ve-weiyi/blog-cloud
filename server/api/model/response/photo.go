package response

import "github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"

type PhotoAlbumDetailsDTO struct {
	*entity.PhotoAlbum
	PhotoCount int64 `json:"photo_count"`
}
