// 查询{{.upperStartCamelObject}}记录
func (m *default{{.upperStartCamelObject}}Model) First(ctx context.Context, conditions string, args ...interface{}) (out *{{.upperStartCamelObject}}, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	out = new({{.upperStartCamelObject}})
	err = db.First(&out).Error
	if err != nil {
		return nil, err
	}
	return out, err
}

// 查询{{.upperStartCamelObject}}总数
func (m *default{{.upperStartCamelObject}}Model) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	err = db.Model(&{{.upperStartCamelObject}}{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 查询{{.upperStartCamelObject}}列表
func (m *default{{.upperStartCamelObject}}Model) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*{{.upperStartCamelObject}}, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

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

// 分页查询{{.upperStartCamelObject}}记录
func (m *default{{.upperStartCamelObject}}Model) FindList(ctx context.Context, limit int, offset int, sorts string, conditions string, args ...interface{}) (list []*{{.upperStartCamelObject}}, err error) {
	// 创建db
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有搜索条件
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	// 如果有排序参数
	if len(sorts) != 0 {
		db = db.Order(sorts)
	}

	// 如果有分页参数
	if limit > 0 || offset > 0 {
		db = db.Limit(limit).Offset(offset)
	}

	// 查询数据
	err = db.Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

// 查询{{.upperStartCamelObject}}记录
func (m *default{{.upperStartCamelObject}}Model) FindOne(ctx context.Context, id int64) (out *{{.upperStartCamelObject}}, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}
