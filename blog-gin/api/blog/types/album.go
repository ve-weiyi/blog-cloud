package types

type Album struct {
	Id         int64  `json:"id"`          // 主键
	AlbumName  string `json:"album_name"`  // 相册名
	AlbumDesc  string `json:"album_desc"`  // 相册描述
	AlbumCover string `json:"album_cover"` // 相册封面
}

type QueryAlbumReq struct {
	PageQuery
}

type QueryPhotoReq struct {
	AlbumId int64 `json:"album_id"` // 相册ID
}
