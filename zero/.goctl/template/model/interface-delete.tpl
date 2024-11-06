Delete(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) (rows int64, err error)
DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
// add extra method in here
