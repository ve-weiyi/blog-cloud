package quickstart

import (
	"fmt"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/server/tools/quickstart/inject"
	"github.com/ve-weiyi/ve-blog-golang/server/tools/quickstart/invent"
	"github.com/ve-weiyi/ve-blog-golang/server/tools/quickstart/invent/field"
	"github.com/ve-weiyi/ve-blog-golang/server/tools/quickstart/invent/model"
	"github.com/ve-weiyi/ve-blog-golang/server/tools/quickstart/tmpl"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/jsonconv"
)

// AutoCodeModel 初始版本自动化代码工具
type AutoCodeModel struct {
	DbName              string // 数据库名
	TableName           string // 表名 				auto_code
	UpperStartCamelName string // Struct名称 		AutoCode 大写驼峰命名
	LowerStartCamelName string // Struct变量名 		autoCode 小写驼峰命名
	SnakeName           string // StructJson名		auto_code api路径前缀
	CommentName         string // Struct中文名称 	「代码」	创建api的描述和注释
	Fields              []*field.Field
	ImportPkgPaths      []string
}

type TableParser struct {
	Config
}

func NewTableParser(config Config) *TableParser {
	return &TableParser{
		Config: config,
	}
}

func (t *TableParser) ParseModelFromSchema() []*AutoCodeModel {
	var models []*AutoCodeModel
	//dbName := t.DbEngin.Migrator().CurrentDatabase()
	tables, err := t.DbEngin.Migrator().GetTables()
	if err != nil {
		return nil
	}
	for _, table := range tables {
		m := t.ParseModelFromTable(table)
		if err != nil {
			return nil
		}
		models = append(models, m)
	}
	return models
}

func (t *TableParser) ParseModelFromTable(tableName string) *AutoCodeModel {
	dbName := t.DbEngin.Migrator().CurrentDatabase()
	tableInfos, err := t.DbEngin.Migrator().TableType(tableName)
	tableComment, _ := tableInfos.Comment()

	table, err := model.GetTable(t.DbEngin, tableName)
	if err != nil {
		return nil
	}

	out := &AutoCodeModel{
		DbName:              dbName,
		TableName:           tableName,
		UpperStartCamelName: jsonconv.Case2Camel(tableName),
		LowerStartCamelName: jsonconv.Case2CamelLowerStart(tableName),
		SnakeName:           jsonconv.Camel2Case(tableName),
		CommentName:         tableComment,
		Fields:              t.ConvertField(table.Columns),
		ImportPkgPaths:      []string{
			//"github.com/ve-weiyi/ve-blog-golang/server/api/blog/controller/svc",
			//"github.com/ve-weiyi/ve-blog-golang/server/api/blog/model/entity",
			//"github.com/ve-weiyi/ve-blog-golang/server/api/blog/model/request",
			//"github.com/ve-weiyi/ve-blog-golang/server/api/common/controller",
			//"github.com/ve-weiyi/ve-blog-golang/server/api/common/model/response",
		},
	}

	return out
}

func (t *TableParser) ConvertField(columns []*model.Column) []*field.Field {
	var out []*field.Field
	cfg := t.FieldConfig
	for _, c := range columns {
		comment, _ := c.Comment()

		f := &field.Field{
			Name:             jsonconv.Case2Camel(c.Name()),
			Type:             c.FiledType(&cfg),
			ColumnName:       c.Name(),
			ColumnComment:    comment,
			MultilineComment: strings.Contains(comment, "\n"),
			Tag:              map[string]string{field.TagKeyJson: cfg.FieldJSONTagNS(c.Name())},
			GORMTag:          c.BuildGormTag(),
		}

		out = append(out, f)
	}

	return out
}

func (t *TableParser) GenerateInjectMetas(models ...*AutoCodeModel) []*inject.AstInjectMeta {
	var injectMetas []*inject.AstInjectMeta

	for _, data := range models {
		temporaryRoot := t.OutPath

		injectMetas = append(injectMetas, &inject.AstInjectMeta{
			Key:      tmpl.KeyRepository,
			FilePath: fmt.Sprintf("%v/repository/repository.go", temporaryRoot),
			StructMetas: []*inject.StructMeta{
				inject.NewStructMete("AppRepository", fmt.Sprintf(`%vRepository *logic.%vRepository //%v`, data.UpperStartCamelName, data.UpperStartCamelName, data.CommentName)),
			},
			FuncMetas: []*inject.FuncMeta{
				inject.NewFuncMete("NewRepository", fmt.Sprintf(`return &AppRepository{
			%vRepository: logic.New%vRepository(svcCtx),
			}`, data.UpperStartCamelName, data.UpperStartCamelName)),
			},
		})

		injectMetas = append(injectMetas, &inject.AstInjectMeta{
			Key:      tmpl.KeyService,
			FilePath: fmt.Sprintf("%v/service/service.go", temporaryRoot),
			StructMetas: []*inject.StructMeta{
				inject.NewStructMete("AppService", fmt.Sprintf(`%vService *logic.%vService //%v`, data.UpperStartCamelName, data.UpperStartCamelName, data.CommentName)),
			},
			FuncMetas: []*inject.FuncMeta{
				inject.NewFuncMete("NewService", fmt.Sprintf(`return &AppService{
			%vService: logic.New%vService(svcCtx),
			}`, data.UpperStartCamelName, data.UpperStartCamelName)),
			},
		})

		injectMetas = append(injectMetas, &inject.AstInjectMeta{
			Key:      tmpl.KeyController,
			FilePath: fmt.Sprintf("%v/controller/controller.go", temporaryRoot),
			StructMetas: []*inject.StructMeta{
				inject.NewStructMete("AppController", fmt.Sprintf(`%vController *logic.%vController //%v`, data.UpperStartCamelName, data.UpperStartCamelName, data.CommentName)),
			},
			FuncMetas: []*inject.FuncMeta{
				inject.NewFuncMete("NewController", fmt.Sprintf(`return &AppController{
			%vController: logic.New%vController(svcCtx),
			}`, data.UpperStartCamelName, data.UpperStartCamelName)),
			},
		})

		injectMetas = append(injectMetas, &inject.AstInjectMeta{
			Key:      tmpl.KeyRouter,
			FilePath: fmt.Sprintf("%v/router/router.go", temporaryRoot),
			StructMetas: []*inject.StructMeta{
				inject.NewStructMete("AppRouter", fmt.Sprintf(`%vRouter *logic.%vRouter //%v`, data.UpperStartCamelName, data.UpperStartCamelName, data.CommentName)),
			},
			FuncMetas: []*inject.FuncMeta{
				inject.NewFuncMete("NewRouter", fmt.Sprintf(`return &AppRouter{
			%vRouter: logic.New%vRouter(svcCtx),
			}`, data.UpperStartCamelName, data.UpperStartCamelName)),
			},
		})

		injectMetas = append(injectMetas, &inject.AstInjectMeta{
			Key:      tmpl.KeyRouter,
			FilePath: fmt.Sprintf("%v/router/logic/register.rt.go", temporaryRoot),
			DeclMetas: []*inject.DeclMeta{inject.NewDeclMeta(fmt.Sprintf(`
	// 初始化 %s 路由信息
	// publicRouter 公开路由，不登录就可以访问
	// loginRouter  登录路由，登录后才可以访问
	func (s *%sRouter) Init%sRouter(publicRouter *gin.RouterGroup, loginRouter *gin.RouterGroup) {
		s.Init%sBasicRouter(publicRouter, loginRouter)
	}
`, data.UpperStartCamelName, data.UpperStartCamelName, data.UpperStartCamelName, data.UpperStartCamelName))},
		})
	}
	return injectMetas
}

func (t *TableParser) GenerateInventMetas(models ...*AutoCodeModel) []*invent.TemplateMeta {
	var metas []*invent.TemplateMeta

	for _, data := range models {
		temporaryRoot := t.OutPath
		fileName := t.OutFileNS(data.TableName)
		mode := t.ReplaceMode
		metas = append(metas, &invent.TemplateMeta{
			Key:            tmpl.KeyModel,
			TemplateString: tmpl.Model,
			CodeOutPath:    fmt.Sprintf("%v/model/entity/%s.go", temporaryRoot, fileName),
			Data:           data,
			Mode:           mode,
		})
		metas = append(metas, &invent.TemplateMeta{
			Key:            tmpl.KeyRepository,
			TemplateString: tmpl.Repository,
			CodeOutPath:    fmt.Sprintf("%v/repository/logic/%s.rp.go", temporaryRoot, fileName),
			Data:           data,
			Mode:           mode,
		})
		metas = append(metas, &invent.TemplateMeta{
			Key:            tmpl.KeyService,
			TemplateString: tmpl.Service,
			CodeOutPath:    fmt.Sprintf("%v/service/logic/%s.sv.go", temporaryRoot, fileName),
			Data:           data,
			Mode:           mode,
		})
		metas = append(metas, &invent.TemplateMeta{
			Key:            tmpl.KeyController,
			TemplateString: tmpl.Controller,
			CodeOutPath:    fmt.Sprintf("%v/controller/logic/%s.ctl.go", temporaryRoot, fileName),
			Data:           data,
			Mode:           mode,
		})
		metas = append(metas, &invent.TemplateMeta{
			Key:            tmpl.KeyRouter,
			TemplateString: tmpl.Router,
			CodeOutPath:    fmt.Sprintf("%v/router/logic/%s.rt.go", temporaryRoot, fileName),
			Data:           data,
			Mode:           mode,
		})
		//metas = append(metas, &invent.TemplateMeta{
		//	Key:            tmpl.KeyApi,
		//	TemplateString: tmpl.Api,
		//	CodeOutPath:    fmt.Sprintf("%v/api/%s.ts", temporaryRoot, fileName),
		//	Data:           data,
		//	Mode:           mode,
		//})
	}

	return metas
}

func (t *TableParser) GenerateCommonInventMetas(models ...*AutoCodeModel) []*invent.TemplateMeta {
	var metas []*invent.TemplateMeta

	for _, data := range models {
		temporaryRoot := t.OutPath
		fileName := t.OutFileNS(data.TableName)
		mode := t.ReplaceMode

		metas = append(metas, &invent.TemplateMeta{
			TemplateString: tmpl.CommonRepository,
			CodeOutPath:    fmt.Sprintf("%v/repository/logic/%s.rp.go", temporaryRoot, fileName),
			Data:           data,
			Mode:           mode,
		})

		metas = append(metas, &invent.TemplateMeta{
			TemplateString: tmpl.CommonService,
			CodeOutPath:    fmt.Sprintf("%v/service/logic/%s.sv.go", temporaryRoot, fileName),
			Data:           data,
			Mode:           mode,
		})

		metas = append(metas, &invent.TemplateMeta{
			TemplateString: tmpl.CommonController,
			CodeOutPath:    fmt.Sprintf("%v/controller/logic/%s.ctl.go", temporaryRoot, fileName),
			Data:           data,
			Mode:           mode,
		})
		metas = append(metas, &invent.TemplateMeta{
			TemplateString: tmpl.CommonRouter,
			CodeOutPath:    fmt.Sprintf("%v/router/logic/%s.rt.go", temporaryRoot, fileName),
			Data:           data,
			Mode:           mode,
		})
	}
	return metas
}

func (t *TableParser) GeneratePkgMetas(models ...*AutoCodeModel) ([]*invent.TemplateMeta, []*inject.AstInjectMeta) {
	var metas []*invent.TemplateMeta

	for _, data := range models {
		temporaryRoot := t.OutPath
		fileName := "context"
		mode := t.ReplaceMode

		//context 是下层引用
		metas = append(metas, &invent.TemplateMeta{
			TemplateString: tmpl.RouterContext,
			CodeOutPath:    fmt.Sprintf("%v/router/svc/%s.rt.go", temporaryRoot, fileName),
			Data:           data,
			Mode:           mode,
		})
		metas = append(metas, &invent.TemplateMeta{
			TemplateString: tmpl.ControllerContext,
			CodeOutPath:    fmt.Sprintf("%v/controller/svc/%s.ctl.go", temporaryRoot, fileName),
			Data:           data,
			Mode:           mode,
		})
		metas = append(metas, &invent.TemplateMeta{
			TemplateString: tmpl.ServiceContext,
			CodeOutPath:    fmt.Sprintf("%v/service/svc/%s.sv.go", temporaryRoot, fileName),
			Data:           data,
			Mode:           mode,
		})
		metas = append(metas, &invent.TemplateMeta{
			TemplateString: tmpl.RepositoryContext,
			CodeOutPath:    fmt.Sprintf("%v/repository/svc/%s.rp.go", temporaryRoot, fileName),
			Data:           data,
			Mode:           mode,
		})

		//入口文件
		metas = append(metas, &invent.TemplateMeta{
			TemplateString: tmpl.AppRouter,
			CodeOutPath:    fmt.Sprintf("%v/router/router.go", temporaryRoot),
			Data:           data,
			Mode:           mode,
		})
		metas = append(metas, &invent.TemplateMeta{
			TemplateString: tmpl.AppController,
			CodeOutPath:    fmt.Sprintf("%v/controller/controller.go", temporaryRoot),
			Data:           data,
			Mode:           mode,
		})
		metas = append(metas, &invent.TemplateMeta{
			TemplateString: tmpl.AppService,
			CodeOutPath:    fmt.Sprintf("%v/service/service.go", temporaryRoot),
			Data:           data,
			Mode:           mode,
		})
		metas = append(metas, &invent.TemplateMeta{
			TemplateString: tmpl.AppRepository,
			CodeOutPath:    fmt.Sprintf("%v/repository/repository.go", temporaryRoot),
			Data:           data,
			Mode:           mode,
		})
	}

	return metas, nil
}
