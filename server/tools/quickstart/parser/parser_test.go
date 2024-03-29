package parser

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zeromicro/go-zero/tools/goctl/model/sql/model"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/ve-weiyi/ve-blog-golang/server/global"
	"github.com/ve-weiyi/ve-blog-golang/server/tools/quickstart"
	"github.com/ve-weiyi/ve-blog-golang/server/tools/quickstart/invent"
	"github.com/ve-weiyi/ve-blog-golang/server/tools/quickstart/invent/field"
	"github.com/ve-weiyi/ve-blog-golang/server/tools/quickstart/parser/gz_model"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/jsonconv"
)

func TestParsePlainText(t *testing.T) {
	sqlFile := filepath.Join(pathx.MustTempDir(), "tmp.sql")
	err := os.WriteFile(sqlFile, []byte("plain text"), 0o777)
	assert.Nil(t, err)

	_, err = Parse(sqlFile, "go_zero", false)
	assert.NotNil(t, err)
}

func TestParseSelect(t *testing.T) {
	sqlFile := filepath.Join(pathx.MustTempDir(), "tmp.sql")
	err := os.WriteFile(sqlFile, []byte("select * from user"), 0o777)
	assert.Nil(t, err)

	tables, err := Parse(sqlFile, "go_zero", false)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(tables))
}

//go:embed testdata/user.sql
var user string

func TestParseCreateTable(t *testing.T) {
	sqlFile := filepath.Join(pathx.MustTempDir(), "tmp.sql")
	err := os.WriteFile(sqlFile, []byte(user), 0o777)
	assert.Nil(t, err)

	tables, err := Parse(sqlFile, "go_zero", false)
	assert.Equal(t, 1, len(tables))
	table := tables[0]
	assert.Nil(t, err)
	assert.Equal(t, "test_user", table.Name.Source())
	assert.Equal(t, "id", table.PrimaryKey.Name.Source())
	assert.Equal(t, true, table.ContainsTime())
	assert.Equal(t, 2, len(table.UniqueIndex))

	log.Printf("table %+v", jsonconv.ObjectToJsonIndent(table))

	var fs []*field.Field

	assert.True(t, func() bool {
		for _, e := range table.Fields {
			if e.Comment != util.TrimNewLine(e.Comment) {
				return false
			}

			log.Printf("%+v", jsonconv.ObjectToJsonIndent(e))

			fs = append(fs, &field.Field{
				Name:          jsonconv.Case2Camel(e.Name.Source()),
				Type:          e.DataType,
				ColumnName:    e.Name.Source(),
				ColumnComment: e.Comment,
				Tag:           map[string]string{field.TagKeyJson: e.Name.Source()},
			})
		}

		return true
	}())

	meta := invent.TemplateMeta{
		Key:            "",
		Mode:           invent.ModeCreateOrReplace,
		CodeOutPath:    "./test/test.go",
		TemplateString: gz_model.GoZeroModel,
		FunMap:         nil,
		Data: quickstart.AutoCodeModel{
			DbName:              table.Db.Source(),
			TableName:           table.Name.Source(),
			UpperStartCamelName: jsonconv.Case2Camel(table.Name.Source()),
			LowerStartCamelName: jsonconv.Case2CamelLowerStart(table.Name.Source()),
			SnakeName:           jsonconv.Camel2Case(table.Name.Source()),
			CommentName:         table.Name.Source(),
			Fields:              fs,
			ImportPkgPaths:      nil,
		},
	}

	err = meta.CreateTempFile()
	if err != nil {
		log.Println(err)
	}

}

func TestParseDatabase(t *testing.T) {
	tables, err := Parse(filepath.Join(global.GetRuntimeRoot(), "server/test.sql"), "go_zero", false)
	assert.Nil(t, err)

	log.Printf("table %+v", jsonconv.ObjectToJsonIndent(tables))

	for _, table := range tables {
		var fs []*field.Field
		for _, e := range table.Fields {

			log.Printf("%+v", jsonconv.ObjectToJsonIndent(e))

			fs = append(fs, &field.Field{
				Name:          jsonconv.Case2Camel(e.Name.Source()),
				Type:          e.DataType,
				ColumnName:    e.Name.Source(),
				ColumnComment: e.Comment,
				Tag:           map[string]string{field.TagKeyJson: e.Name.Source()},
			})
		}

		meta := invent.TemplateMeta{
			Key:            "",
			Mode:           invent.ModeCreateOrReplace,
			CodeOutPath:    fmt.Sprintf("./test/%v_model.go", table.Name.Source()),
			TemplateString: gz_model.GoZeroModel,
			FunMap:         nil,
			Data: quickstart.AutoCodeModel{
				DbName:              table.Db.Source(),
				TableName:           table.Name.Source(),
				UpperStartCamelName: jsonconv.Case2Camel(table.Name.Source()),
				LowerStartCamelName: jsonconv.Case2CamelLowerStart(table.Name.Source()),
				SnakeName:           jsonconv.Camel2Case(table.Name.Source()),
				CommentName:         table.Name.Source(),
				Fields:              fs,
				ImportPkgPaths:      nil,
			},
		}

		err = meta.CreateTempFile()
		if err != nil {
			log.Println(err)
		}
	}
}

func TestConvertColumn(t *testing.T) {
	t.Run("missingPrimaryKey", func(t *testing.T) {
		columnData := model.ColumnData{
			Db:    "user",
			Table: "user",
			Columns: []*model.Column{
				{
					DbColumn: &model.DbColumn{
						Name:     "id",
						DataType: "bigint",
					},
				},
			},
		}
		_, err := columnData.Convert()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "missing primary key")
	})

	t.Run("jointPrimaryKey", func(t *testing.T) {
		columnData := model.ColumnData{
			Db:    "user",
			Table: "user",
			Columns: []*model.Column{
				{
					DbColumn: &model.DbColumn{
						Name:     "id",
						DataType: "bigint",
					},
					Index: &model.DbIndex{
						IndexName: "PRIMARY",
					},
				},
				{
					DbColumn: &model.DbColumn{
						Name:     "mobile",
						DataType: "varchar",
						Comment:  "手机号",
					},
					Index: &model.DbIndex{
						IndexName: "PRIMARY",
					},
				},
			},
		}
		_, err := columnData.Convert()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "joint primary key is not supported")
	})

	t.Run("normal", func(t *testing.T) {
		columnData := model.ColumnData{
			Db:    "user",
			Table: "user",
			Columns: []*model.Column{
				{
					DbColumn: &model.DbColumn{
						Name:     "id",
						DataType: "bigint",
						Extra:    "auto_increment",
					},
					Index: &model.DbIndex{
						IndexName:  "PRIMARY",
						SeqInIndex: 1,
					},
				},
				{
					DbColumn: &model.DbColumn{
						Name:     "mobile",
						DataType: "varchar",
						Comment:  "手机号",
					},
					Index: &model.DbIndex{
						IndexName:  "mobile_unique",
						SeqInIndex: 1,
					},
				},
			},
		}

		table, err := columnData.Convert()
		assert.Nil(t, err)
		assert.True(t, table.PrimaryKey.Index.IndexName == "PRIMARY" && table.PrimaryKey.Name == "id")
		for _, item := range table.Columns {
			if item.Name == "mobile" {
				assert.True(t, item.Index.NonUnique == 0)
				break
			}
		}
	})
}
