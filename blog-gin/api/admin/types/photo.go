package types

type NewPhotoReq struct {
	Id        int64  `json:"id,optional"` // 主键
	AlbumId   int64  `json:"album_id"`    // 相册id
	PhotoName string `json:"photo_name"`  // 照片名
	PhotoDesc string `json:"photo_desc"`  // 照片描述
	PhotoSrc  string `json:"photo_src"`   // 照片地址
	IsDelete  int64  `json:"is_delete"`   // 是否删除
}

type PhotoBackVO struct {
	Id        int64  `json:"id,optional"` // 主键
	AlbumId   int64  `json:"album_id"`    // 相册id
	PhotoName string `json:"photo_name"`  // 照片名
	PhotoDesc string `json:"photo_desc"`  // 照片描述
	PhotoSrc  string `json:"photo_src"`   // 照片地址
	IsDelete  int64  `json:"is_delete"`   // 是否删除
	CreatedAt int64  `json:"created_at"`  // 创建时间
	UpdatedAt int64  `json:"updated_at"`  // 更新时间
}

type QueryPhotoReq struct {
	PageQuery
	AlbumId  int64 `json:"album_id,optional"`  // 相册id
	IsDelete int64 `json:"is_delete,optional"` // 是否删除
}

type UpdatePhotoDeleteReq struct {
	Ids      []int64 `json:"ids"`       // 主键
	IsDelete int64   `json:"is_delete"` // 是否删除
}
