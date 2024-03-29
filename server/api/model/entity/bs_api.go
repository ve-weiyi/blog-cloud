package entity

import "time"

// TableNameApi return the table name of <api>
const TableNameApi = "api"

// Api mapped from table <api>
type Api struct {
	ID        int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:主键id" json:"id"`                             // 主键id
	Name      string    `gorm:"column:name;type:varchar(128);not null;uniqueIndex:idx_path_method,priority:3;comment:api名称" json:"name"`      // api名称
	Path      string    `gorm:"column:path;type:varchar(128);not null;uniqueIndex:idx_path_method,priority:1;comment:api路径" json:"path"`      // api路径
	Method    string    `gorm:"column:method;type:varchar(16);not null;uniqueIndex:idx_path_method,priority:2;comment:api请求方法" json:"method"` // api请求方法
	ParentID  int       `gorm:"column:parent_id;type:int;not null;comment:分组id" json:"parent_id"`                                             // 分组id
	Traceable int       `gorm:"column:traceable;type:tinyint;not null;comment:是否追溯操作记录 0需要，1是" json:"traceable"`                              // 是否追溯操作记录 0需要，1是
	Status    int       `gorm:"column:status;type:tinyint;not null;default:1;comment:状态 1开，2关" json:"status"`                                 // 状态 1开，2关
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`            // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`            // 更新时间
}

// TableName Api's table name
func (*Api) TableName() string {
	return TableNameApi
}
