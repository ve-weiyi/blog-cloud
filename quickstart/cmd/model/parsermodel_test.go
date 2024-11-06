package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
	"github.com/ve-weiyi/ve-blog-golang/quickstart/cmd/model/helper"
)

func Test_ParseTableFormDsn(t *testing.T) {
	const dsn = "root:mysql7914@(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	fromDsn, err := helper.ParseTableFromDsn(dsn)
	t.Log(err)
	// t.Log(jsonconv.AnyToJsonIndent(fromDsn))

	const sql = "test.sql"
	fromSql, err := helper.ParseTableFromSql(sql)
	t.Log(err)
	// t.Log(jsonconv.AnyToJsonIndent(fromSql))

	assert.Equal(t, jsonconv.AnyToJsonIndent(fromDsn), jsonconv.AnyToJsonIndent(fromSql))
}

func Test_ParseTableFromSql(t *testing.T) {
	const sql = "test.sql"
	fromSql, err := helper.ParseTableFromSql(sql)
	t.Log(err)
	t.Log(jsonconv.AnyToJsonIndent(fromSql))
}
