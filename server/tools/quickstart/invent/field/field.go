package field

import (
	"strings"

	"gorm.io/gen/field"
)

type Field struct {
	Name             string `json:"name"`
	Type             string `json:"type"`           // Field数据类型
	ColumnName       string `json:"column_name"`    // 数据库字段
	ColumnComment    string `json:"column_comment"` // 数据库字段描述
	MultilineComment bool
	Tag              Tag
	GORMTag          GormTag
}

func (m *Field) Tags() string {
	if _, ok := m.Tag[field.TagKeyGorm]; ok {
		return m.Tag.Build()
	}

	if gormTag := strings.TrimSpace(m.GORMTag.Build()); gormTag != "" {
		m.Tag.Set(field.TagKeyGorm, gormTag)
	}
	return m.Tag.Build()
}
