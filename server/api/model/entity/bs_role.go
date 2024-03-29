package entity

import "time"

// TableNameRole return the table name of <role>
const TableNameRole = "role"

// Role mapped from table <role>
type Role struct {
	ID          int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:主键id" json:"id"`                  // 主键id
	RolePID     int       `gorm:"column:role_pid;type:int;not null;comment:父角色id" json:"role_pid"`                                   // 父角色id
	RoleDomain  string    `gorm:"column:role_domain;type:varchar(64);not null;default:0;comment:角色域" json:"role_domain"`             // 角色域
	RoleName    string    `gorm:"column:role_name;type:varchar(64);not null;comment:角色名" json:"role_name"`                           // 角色名
	RoleComment string    `gorm:"column:role_comment;type:varchar(64);not null;comment:角色备注" json:"role_comment"`                    // 角色备注
	IsDisable   int       `gorm:"column:is_disable;type:tinyint(1);not null;comment:是否禁用  0否 1是" json:"is_disable"`                  // 是否禁用  0否 1是
	IsDefault   int       `gorm:"column:is_default;type:tinyint(1);not null;comment:是否默认角色 0否 1是" json:"is_default"`                 // 是否默认角色 0否 1是
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName Role's table name
func (*Role) TableName() string {
	return TableNameRole
}
