
// 更新{{.upperStartCamelObject}}记录
func (m *default{{.upperStartCamelObject}}Model) Update(ctx context.Context, in *{{.upperStartCamelObject}}) (out *{{.upperStartCamelObject}}, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Updates(&in).Error
	if err != nil {
		return nil, err
	}
	return in, err
}

// 更新{{.upperStartCamelObject}}记录
func (m *default{{.upperStartCamelObject}}Model) UpdateColumns(ctx context.Context, id int64, columns map[string]interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where("`id` = ?", id).UpdateColumns(&columns)
	if result.Error != nil {
		return 0, err
	}

	return result.RowsAffected, err
}

// 更新{{.upperStartCamelObject}}记录
func (m *default{{.upperStartCamelObject}}Model) Save(ctx context.Context, in *{{.upperStartCamelObject}}) (out *{{.upperStartCamelObject}}, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Save(&in).Error
	if err != nil {
		return nil, err
	}
	return in, err
}