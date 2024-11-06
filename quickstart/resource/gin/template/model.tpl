package model

import (
	"context"

	"gorm.io/gorm"
)

var _ {{.UpperStartCamelName}}Model = (*default{{.UpperStartCamelName}}Model)(nil)

type (
	// 接口定义
	{{.UpperStartCamelName}}Model interface {
		TableName() string
		// 在事务中操作
		WithTransaction(tx *gorm.DB) (out {{.UpperStartCamelName}}Model)
		// 插入
		Insert(in *{{.UpperStartCamelName}}) (rows int64, err error)
		InsertBatch(in ...*{{.UpperStartCamelName}}) (rows int64, err error)
        // 删除
		Delete(id int64) (rows int64, err error)
		DeleteBatch(conditions string, args ...interface{}) (rows int64, err error)
        // 更新
        Update(in *{{.UpperStartCamelName}}) (rows int64, err error)
        Save(in *{{.UpperStartCamelName}}) (rows int64, err error)
		// 查询
		FindOne(id int64) (out *{{.UpperStartCamelName}}, err error)
		First(conditions string, args ...interface{}) (out *{{.UpperStartCamelName}}, err error)
		FindCount(conditions string, args ...interface{}) (count int64, err error)
		FindALL(conditions string, args ...interface{}) (list []*{{.UpperStartCamelName}}, err error)
		FindList(page int, size int, sorts string, conditions string, args ...interface{}) (list []*{{.UpperStartCamelName}}, err error)
	    // add extra method in here
        {{- range $key, $value := .UniqueFields}}
            FindOneBy{{ funcFieldsKey $value}}({{funcFieldsKeyVar $value}}) (out *{{$.UpperStartCamelName}},err error)
        {{- end}}
	}

	// 表字段定义
	{{.UpperStartCamelName}} struct {
	{{range .Fields}}
        {{.Group}} {{.Type}} `{{.Tags}}` {{if .Comment}}// {{.Comment}}{{end}}
	{{- end}}
	}

    // 接口实现
    default{{.UpperStartCamelName}}Model struct {
        DbEngin    *gorm.DB
        table       string
    }
)

func New{{.UpperStartCamelName}}Model(db *gorm.DB) {{.UpperStartCamelName}}Model {
	return &default{{.UpperStartCamelName}}Model{
		DbEngin:    db,
		table:      "`{{.SnakeName}}`",
	}
}

func (m *default{{.UpperStartCamelName}}Model) TableName() string {
	return m.table
}

// 在事务中操作
func (m *default{{.UpperStartCamelName}}Model) WithTransaction(tx *gorm.DB) (out {{.UpperStartCamelName}}Model) {
	return New{{.UpperStartCamelName}}Model(tx)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *default{{.UpperStartCamelName}}Model) Insert(in *{{.UpperStartCamelName}}) (rows int64, err error) {
	db := m.DbEngin.Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录
func (m *default{{.UpperStartCamelName}}Model) InsertBatch(in ...*{{.UpperStartCamelName}}) (rows int64, err error) {
	db := m.DbEngin.Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *default{{.UpperStartCamelName}}Model) Delete(id int64) (rows int64, err error) {
	db := m.DbEngin.Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&{{.UpperStartCamelName}}{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *default{{.UpperStartCamelName}}Model) DeleteBatch(conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&{{.UpperStartCamelName}}{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存记录（更新零值）
func (m *default{{.UpperStartCamelName}}Model) Save(in *{{.UpperStartCamelName}}) (rows int64, err error) {
	db := m.DbEngin.Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *default{{.UpperStartCamelName}}Model) Update(in *{{.UpperStartCamelName}}) (rows int64, err error) {
	db := m.DbEngin.Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *default{{.UpperStartCamelName}}Model) FindOne(id int64) (out *{{.UpperStartCamelName}}, err error) {
	db := m.DbEngin.Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *default{{.UpperStartCamelName}}Model) First(conditions string, args ...interface{}) (out *{{.UpperStartCamelName}}, err error) {
	db := m.DbEngin.Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	err = db.First(&out).Error
	if err != nil {
		return nil, err
	}
	return out, err
}

// 查询总数
func (m *default{{.UpperStartCamelName}}Model) FindCount(conditions string, args ...interface{}) (count int64, err error) {
	db := m.DbEngin.Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	err = db.Model(&{{.UpperStartCamelName}}{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 查询列表
func (m *default{{.UpperStartCamelName}}Model) FindALL(conditions string, args ...interface{}) (out []*{{.UpperStartCamelName}}, err error) {
	db := m.DbEngin.Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	err = db.Find(&out).Error
	if err != nil {
		return nil, err
	}
	return out, err
}

// 分页查询记录
func (m *default{{.UpperStartCamelName}}Model) FindList(page int, size int, sorts string, conditions string, args ...interface{}) (list []*{{.UpperStartCamelName}}, err error) {
	// 插入db
	db := m.DbEngin.Table(m.table)

	// 如果有搜索条件
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	// 如果有排序参数
	if len(sorts) != 0 {
		db = db.Order(sorts)
	}

	// 如果有分页参数
	if page > 0 && size > 0 {
		limit := size
		offset := (page - 1) * limit
		db = db.Limit(limit).Offset(offset)
	}

	// 查询数据
	err = db.Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

// add extra method in here
{{- range $key, $value := .UniqueFields}}
func (m *default{{$.UpperStartCamelName}}Model) FindOneBy{{ funcFieldsKey $value}}({{funcFieldsKeyVar $value}}) (out *{{$.UpperStartCamelName}},err error) {
	db := m.DbEngin.Table(m.table)

	err = db.Where("{{funcFieldsKeyCond $value}}", {{funcFieldsKeyCondVar $value}}).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
{{- end}}
