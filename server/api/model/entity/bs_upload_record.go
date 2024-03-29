package entity

import "time"

// TableNameUploadRecord return the table name of <upload_record>
const TableNameUploadRecord = "upload_record"

// UploadRecord mapped from table <upload_record>
type UploadRecord struct {
	ID        int       `gorm:"column:id;type:int;primaryKey;autoIncrement:true;comment:id" json:"id"`                              // id
	UserID    int       `gorm:"column:user_id;type:int;not null;index:idx_uid,priority:1;comment:用户id" json:"user_id"`              // 用户id
	Label     string    `gorm:"column:label;type:varchar(128);not null;comment:标签" json:"label"`                                    // 标签
	FileName  string    `gorm:"column:file_name;type:varchar(64);not null;comment:文件名称" json:"file_name"`                           // 文件名称
	FileSize  int       `gorm:"column:file_size;type:int;not null;comment:文件大小" json:"file_size"`                                   // 文件大小
	FileMd5   string    `gorm:"column:file_md5;type:varchar(128);not null;comment:文件md5值" json:"file_md5"`                          // 文件md5值
	FileURL   string    `gorm:"column:file_url;type:varchar(256);not null;comment:上传路径" json:"file_url"`                            // 上传路径
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName UploadRecord's table name
func (*UploadRecord) TableName() string {
	return TableNameUploadRecord
}
