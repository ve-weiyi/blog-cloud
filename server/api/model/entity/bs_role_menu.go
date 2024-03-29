package entity

// TableNameRoleMenu return the table name of <role_menu>
const TableNameRoleMenu = "role_menu"

// RoleMenu mapped from table <role_menu>
type RoleMenu struct {
	ID     int `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:主键id" json:"id"` // 主键id
	RoleID int `gorm:"column:role_id;type:int;not null;comment:角色id" json:"role_id"`                     // 角色id
	MenuID int `gorm:"column:menu_id;type:int;not null;comment:菜单id" json:"menu_id"`                     // 菜单id
}

// TableName RoleMenu's table name
func (*RoleMenu) TableName() string {
	return TableNameRoleMenu
}
