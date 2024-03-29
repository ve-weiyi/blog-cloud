package entity

// TableNameUserRole return the table name of <user_role>
const TableNameUserRole = "user_role"

// UserRole mapped from table <user_role>
type UserRole struct {
	ID     int `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:主键id" json:"id"` // 主键id
	UserID int `gorm:"column:user_id;type:int;not null;comment:用户id" json:"user_id"`                     // 用户id
	RoleID int `gorm:"column:role_id;type:int;not null;comment:角色id" json:"role_id"`                     // 角色id
}

// TableName UserRole's table name
func (*UserRole) TableName() string {
	return TableNameUserRole
}
