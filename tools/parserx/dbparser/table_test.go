package dbparser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/jsonconv"
)

func Test_ParseTableFormDsn(t *testing.T) {
	const dsn = "root:mysql7914@(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	fromDsn, err := ParseTableFromDsn(dsn)
	t.Log(err)
	t.Log(jsonconv.AnyToJsonIndent(fromDsn))
}

func Test_ParseTableFromSql(t *testing.T) {
	const sql = "../../testdata/test.sql"
	fromSql, err := ParseTableFromSql(sql)
	t.Log(err)
	t.Log(jsonconv.AnyToJsonIndent(fromSql))
}

func Test_Equal(t *testing.T) {
	const dsn = "root:mysql7914@(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	fromDsn, err := ParseTableFromDsn(dsn)
	t.Log(err)

	const sql = "test.sql"
	fromSql, err := ParseTableFromSql(sql)
	t.Log(err)

	assert.Equal(t, jsonconv.AnyToJsonIndent(fromDsn), jsonconv.AnyToJsonIndent(fromSql))
}
