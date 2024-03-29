package entity

import "time"

// TableNameMenu return the table name of <menu>
const TableNameMenu = "menu"

// Menu mapped from table <menu>
type Menu struct {
	ID        int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:主键" json:"id"`                             // 主键
	ParentID  int       `gorm:"column:parent_id;type:int;not null;comment:父id" json:"parent_id"`                                            // 父id
	Title     string    `gorm:"column:title;type:varchar(32);not null;comment:菜单标题" json:"title"`                                           // 菜单标题
	Path      string    `gorm:"column:path;type:varchar(64);not null;uniqueIndex:uk_path,priority:1;comment:路由路径" json:"path"`              // 路由路径
	Name      string    `gorm:"column:name;type:varchar(32);not null;comment:路由名称" json:"name"`                                             // 路由名称
	Component string    `gorm:"column:component;type:varchar(64);not null;comment:路由组件" json:"component"`                                   // 路由组件
	Redirect  string    `gorm:"column:redirect;type:varchar(64);not null;comment:路由重定向" json:"redirect"`                                    // 路由重定向
	Type      int       `gorm:"column:type;type:tinyint unsigned;not null;default:0;comment:菜单类型（0代表菜单、1代表iframe、2代表外链、3代表按钮）" json:"type"` // 菜单类型（0代表菜单、1代表iframe、2代表外链、3代表按钮）
	Meta      string    `gorm:"column:meta;type:varchar(64);not null;comment:菜单元数据" json:"meta"`                                            // 菜单元数据
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`          // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`          // 更新时间
}

// TableName Menu's table name
func (*Menu) TableName() string {
	return TableNameMenu
}
