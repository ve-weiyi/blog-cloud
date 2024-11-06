
func (m *default{{.upperStartCamelObject}}Model) FindOneBy{{.upperField}}(ctx context.Context, {{.in}}) (out *{{.upperStartCamelObject}},err error) {
	db := m.DbEngin.WithContext(ctx)

	err = db.Where("{{.originalField}}", {{.lowerStartCamelField}}).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
