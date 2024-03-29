package entity

import "time"

// TableNamePhoto return the table name of <photo>
const TableNamePhoto = "photo"

// Photo mapped from table <photo>
type Photo struct {
	ID        int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:主键" json:"id"`                    // 主键
	AlbumID   int       `gorm:"column:album_id;type:int;not null;comment:相册id" json:"album_id"`                                    // 相册id
	PhotoName string    `gorm:"column:photo_name;type:varchar(32);not null;comment:照片名" json:"photo_name"`                         // 照片名
	PhotoDesc string    `gorm:"column:photo_desc;type:varchar(64);not null;comment:照片描述" json:"photo_desc"`                        // 照片描述
	PhotoSrc  string    `gorm:"column:photo_src;type:varchar(255);not null;comment:照片地址" json:"photo_src"`                         // 照片地址
	IsDelete  int       `gorm:"column:is_delete;type:tinyint(1);not null;comment:是否删除" json:"is_delete"`                           // 是否删除
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName Photo's table name
func (*Photo) TableName() string {
	return TableNamePhoto
}
