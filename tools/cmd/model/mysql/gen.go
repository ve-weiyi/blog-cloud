package mysql

import (
	"fmt"
	"os"
	"path"

	"github.com/ve-weiyi/ve-blog-golang/kit/quickstart/invent"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/filex"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
)

type (
	ModelData struct {
		TableName           string
		UpperStartCamelName string
		LowerStartCamelName string
		SnakeName           string
		Fields              []*ModelField
		UniqueFields        [][]*ModelField
	}

	ModelField struct {
		Name    string // 属性名称  Name
		Type    string // 属性类型  string、int、bool、float、{UpperStartCamelName}
		Tag     string // json tag
		Comment string // 行尾注释
	}
)

func generateModel(models []*ModelData, tplPath string, outPath string, nameAs string) error {
	var metas []invent.TemplateMeta

	tpl, err := os.ReadFile(filex.ToAbs(tplPath))
	if err != nil {
		return err
	}

	for _, model := range models {
		meta := invent.TemplateMeta{
			Mode:           invent.ModeCreateOrReplace,
			CodeOutPath:    path.Join(outPath, fmt.Sprintf(nameAs, model.TableName)),
			TemplateString: string(tpl),
			FunMap: map[string]any{
				"funcFieldsKey": func(fs []*ModelField) string {
					var name string
					for _, ff := range fs {
						name += ff.Name
					}
					return name
				},
				"funcFieldsKeyVar": func(fs []*ModelField) string {
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
				"funcFieldsKeyCond": func(fs []*ModelField) string {
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
				"funcFieldsKeyCondVar": func(fs []*ModelField) string {
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
			Data: model,
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
