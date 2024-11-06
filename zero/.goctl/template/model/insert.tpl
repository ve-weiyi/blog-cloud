
// 切换事务操作
func (m *default{{.upperStartCamelObject}}Model) WithTransaction(tx *gorm.DB) (out {{.upperStartCamelObject}}Model) {
    return New{{.upperStartCamelObject}}Model(tx, m.CacheEngin)
}

// 插入{{.upperStartCamelObject}}记录
func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, in *{{.upperStartCamelObject}}) (out *{{.upperStartCamelObject}}, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Create(&in).Error
	if err != nil {
		return nil, err
	}
	return in, err
}

// 插入{{.upperStartCamelObject}}记录
func (m *default{{.upperStartCamelObject}}Model) InsertBatch(ctx context.Context, in ...*{{.upperStartCamelObject}}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.CreateInBatches(&in, len(in)).Error
	if err != nil {
		return 0, err
	}
	return rows, err
}
