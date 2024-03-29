package entity

// TableNameRoleApi return the table name of <role_api>
const TableNameRoleApi = "role_api"

// RoleApi mapped from table <role_api>
type RoleApi struct {
	ID     int `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:主键id" json:"id"` // 主键id
	RoleID int `gorm:"column:role_id;type:int;not null;comment:角色id" json:"role_id"`                     // 角色id
	ApiID  int `gorm:"column:api_id;type:int;not null;comment:接口id" json:"api_id"`                       // 接口id
}

// TableName RoleApi's table name
func (*RoleApi) TableName() string {
	return TableNameRoleApi
}
