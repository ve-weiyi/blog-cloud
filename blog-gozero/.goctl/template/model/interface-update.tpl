Update(ctx context.Context, data *{{.upperStartCamelObject}}) (*{{.upperStartCamelObject}}, error)
UpdateColumns(ctx context.Context, id int64, columns map[string]interface{}) (rows int64, err error)
Save(ctx context.Context, data *{{.upperStartCamelObject}}) (*{{.upperStartCamelObject}}, error)