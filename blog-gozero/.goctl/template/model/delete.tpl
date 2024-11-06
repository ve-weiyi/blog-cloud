
// 删除{{.upperStartCamelObject}}记录
func (m *default{{.upperStartCamelObject}}Model) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	query := db.Delete(&{{.upperStartCamelObject}}{})
	return query.RowsAffected, query.Error
}

// 删除{{.upperStartCamelObject}}记录
func (m *default{{.upperStartCamelObject}}Model) DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&{{.upperStartCamelObject}}{})
	return result.RowsAffected, result.Error
}
