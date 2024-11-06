package model

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/kit/tools/field"
	"github.com/ve-weiyi/ve-blog-golang/kit/tools/invent"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/files"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
	"github.com/ve-weiyi/ve-blog-golang/quickstart/cmd/model/helper"
)

type modelConfig struct {
	SqlFile string
	TplFile string
	OutPath string
	NameAs  string
}

func generateModel(tables []*helper.Table, conf modelConfig) error {
	t := files.ToAbs(conf.TplFile)
	o := conf.OutPath
	n := conf.NameAs

	var metas []invent.TemplateMeta

	tpl, err := os.ReadFile(t)
	if err != nil {
		return err
	}

	for _, table := range tables {
		data := convertTableToData(table)

		meta := invent.TemplateMeta{
			Key:            "",
			Mode:           invent.ModeCreateOrReplace,
			CodeOutPath:    path.Join(o, fmt.Sprintf(n, table.Name)),
			TemplateString: string(tpl),
			FunMap: map[string]any{
				"funcFieldsKey": func(fs []*field.Field) string {
					var name string
					for _, ff := range fs {
						name += ff.Name
					}
					return name
				},
				"funcFieldsKeyVar": func(fs []*field.Field) string {
					var name string
					for _, ff := range fs {
						v := jsonconv.Case2Snake(ff.Name)
						tp := ff.Type
						if name != "" {
							name += ", "
						}
						name += fmt.Sprintf("%s %s", v, tp)
					}
					return name
				},
				"funcFieldsKeyCond": func(fs []*field.Field) string {
					var name string
					for _, ff := range fs {
						v := jsonconv.Case2Snake(ff.Name)
						if name != "" {
							name += " and "
						}
						name += fmt.Sprintf("`%s` = ?", v)
					}
					return name
				},
				"funcFieldsKeyCondVar": func(fs []*field.Field) string {
					var name string
					for _, ff := range fs {
						v := jsonconv.Case2Snake(ff.Name)
						if name != "" {
							name += ", "
						}
						name += v
					}
					return name
				},
			},
			Data: data,
		}
		metas = append(metas, meta)
	}

	for _, m := range metas {
		err := m.Execute()
		if err != nil {
			return err
		}
	}

	return nil
}

func convertTableToData(table *helper.Table) any {

	var fs []*field.Field
	for _, e := range table.Fields {
		fs = append(fs, convertField(e))
	}

	var ufs [][]*field.Field
	for _, e := range table.UniqueIndex {
		var u []*field.Field
		for _, f := range e {
			u = append(u, convertField(f))
		}
		ufs = append(ufs, u)
	}

	data := map[string]any{
		"TableName":           table.Name,
		"UpperStartCamelName": jsonconv.Case2Camel(table.Name),
		"LowerStartCamelName": jsonconv.FirstLower(jsonconv.Case2Camel(table.Name)),
		"SnakeName":           jsonconv.Case2Snake(table.Name),
		"Fields":              fs,
		"UniqueFields":        ufs,
	}

	return data
}

func convertField(e *helper.Field) *field.Field {

	return &field.Field{
		Name:    jsonconv.Case2Camel(e.Name),
		Type:    strings.TrimPrefix(e.DataType, "u"),
		Comment: e.Comment,
		Tag: []field.Tag{
			{
				Name:  "json",
				Value: []string{e.Name},
			},
			{
				Name: "gorm",
				Value: []string{
					fmt.Sprintf("column:%s", e.Name),
				},
			},
		},
		Docs:     nil,
		IsInline: false,
	}
}
