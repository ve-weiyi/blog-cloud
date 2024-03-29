package entity

// TableNameCasbinRule return the table name of <casbin_rule>
const TableNameCasbinRule = "casbin_rule"

// CasbinRule mapped from table <casbin_rule>
type CasbinRule struct {
	ID    int    `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true" json:"id"`
	Ptype string `gorm:"column:ptype;type:varchar(100);not null;uniqueIndex:idx_casbin_rule,priority:1" json:"ptype"`
	V0    string `gorm:"column:v0;type:varchar(100);not null;uniqueIndex:idx_casbin_rule,priority:2" json:"v0"`
	V1    string `gorm:"column:v1;type:varchar(100);not null;uniqueIndex:idx_casbin_rule,priority:3" json:"v1"`
	V2    string `gorm:"column:v2;type:varchar(100);not null;uniqueIndex:idx_casbin_rule,priority:4" json:"v2"`
	V3    string `gorm:"column:v3;type:varchar(100);not null;uniqueIndex:idx_casbin_rule,priority:5" json:"v3"`
	V4    string `gorm:"column:v4;type:varchar(100);not null;uniqueIndex:idx_casbin_rule,priority:6" json:"v4"`
	V5    string `gorm:"column:v5;type:varchar(100);not null;uniqueIndex:idx_casbin_rule,priority:7" json:"v5"`
}

// TableName CasbinRule's table name
func (*CasbinRule) TableName() string {
	return TableNameCasbinRule
}
