package entity

import "time"

// TableNamePhotoAlbum return the table name of <photo_album>
const TableNamePhotoAlbum = "photo_album"

// PhotoAlbum mapped from table <photo_album>
type PhotoAlbum struct {
	ID         int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:主键" json:"id"`                    // 主键
	AlbumName  string    `gorm:"column:album_name;type:varchar(32);not null;comment:相册名" json:"album_name"`                         // 相册名
	AlbumDesc  string    `gorm:"column:album_desc;type:varchar(64);not null;comment:相册描述" json:"album_desc"`                        // 相册描述
	AlbumCover string    `gorm:"column:album_cover;type:varchar(255);not null;comment:相册封面" json:"album_cover"`                     // 相册封面
	IsDelete   int       `gorm:"column:is_delete;type:tinyint(1);not null;comment:是否删除" json:"is_delete"`                           // 是否删除
	Status     int       `gorm:"column:status;type:tinyint;not null;default:1;comment:状态值 1公开 2私密" json:"status"`                   // 状态值 1公开 2私密
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName PhotoAlbum's table name
func (*PhotoAlbum) TableName() string {
	return TableNamePhotoAlbum
}
