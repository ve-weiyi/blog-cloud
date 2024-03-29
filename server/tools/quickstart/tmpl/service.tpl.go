package tmpl

const AppService = `
package service

import (

)

type AppService struct {
	svcCtx *svc.ServiceContext //持有的repository层引用
}

func NewService(svcCtx *svc.ServiceContext) *AppService {
	return &AppService{
		svcCtx: svcCtx,
	}
}
`

const ServiceContext = `
package svc

import (

)

// 注册需要用到的gorm、redis、model
type ServiceContext struct {
	*repository.AppRepository
	Config *config.Config
	Log    *glog.Glogger
	//下面是引用的repository
}

func NewServiceContext(cfg *config.Config) *ServiceContext {
	ctx := svc.NewRepositoryContext(cfg)
	repo := repository.NewRepository(ctx)
	if repo == nil {
		panic("repository cannot be null")
	}

	return &ServiceContext{
		AppRepository: repo,
		Config:        cfg,
		Log:           global.LOG,
	}
}

`

const Service = `
package logic

import (
	{{range .ImportPkgPaths}}{{.}} ` + "\n" + `{{end}}
)

type {{.UpperStartCamelName}}Service struct {
	svcCtx *svc.ServiceContext
}

func New{{.UpperStartCamelName}}Service(svcCtx *svc.ServiceContext) *{{.UpperStartCamelName}}Service {
	return &{{.UpperStartCamelName}}Service{
		svcCtx: svcCtx,
	}
}

// 创建{{.UpperStartCamelName}}记录
func (s *{{.UpperStartCamelName}}Service) Create{{.UpperStartCamelName}}(reqCtx *request.Context, {{.LowerStartCamelName}} *entity.{{.UpperStartCamelName}}) (data *entity.{{.UpperStartCamelName}}, err error) {
	return s.svcCtx.{{.UpperStartCamelName}}Repository.Create(reqCtx, {{.LowerStartCamelName}})
}

// 更新{{.UpperStartCamelName}}记录
func (s *{{.UpperStartCamelName}}Service) Update{{.UpperStartCamelName}}(reqCtx *request.Context, {{.LowerStartCamelName}} *entity.{{.UpperStartCamelName}}) (data *entity.{{.UpperStartCamelName}}, err error) {
	return s.svcCtx.{{.UpperStartCamelName}}Repository.Update(reqCtx, {{.LowerStartCamelName}})
}

// 删除{{.UpperStartCamelName}}记录
func (s *{{.UpperStartCamelName}}Service) Delete{{.UpperStartCamelName}}(reqCtx *request.Context, id int) (rows int64, err error) {
	return s.svcCtx.{{.UpperStartCamelName}}Repository.Delete(reqCtx, "id = ?", id)
}

// 查询{{.UpperStartCamelName}}记录
func (s *{{.UpperStartCamelName}}Service) Find{{.UpperStartCamelName}}(reqCtx *request.Context, id int) (data *entity.{{.UpperStartCamelName}}, err error) {
	return s.svcCtx.{{.UpperStartCamelName}}Repository.First(reqCtx, "id = ?", id)
}

// 批量删除{{.UpperStartCamelName}}记录
func (s *{{.UpperStartCamelName}}Service) Delete{{.UpperStartCamelName}}ByIds(reqCtx *request.Context, ids []int) (rows int64, err error) {
	return s.svcCtx.{{.UpperStartCamelName}}Repository.Delete(reqCtx, "id in (?)", ids)
}

// 分页获取{{.UpperStartCamelName}}记录
func (s *{{.UpperStartCamelName}}Service) Find{{.UpperStartCamelName}}List(reqCtx *request.Context, page *request.PageQuery) (list []*entity.{{.UpperStartCamelName}}, total int64, err error) {
	cond, args := page.ConditionClause()
	order := page.OrderClause()

	list, err = s.svcCtx.{{.UpperStartCamelName}}Repository.FindList(reqCtx, page.Page, page.PageSize, order, cond, args...)
	if err != nil {
		return nil, 0, err
	}
	total, err = s.svcCtx.{{.UpperStartCamelName}}Repository.Count(reqCtx, cond, args...)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
`

const CommonService = `
package logic

import (

)

type {{.UpperStartCamelName}}Service struct {
	svcCtx *svc.ServiceContext
}

func New{{.UpperStartCamelName}}Service(svcCtx *svc.ServiceContext) *{{.UpperStartCamelName}}Service {
	return &{{.UpperStartCamelName}}Service{
		svcCtx: svcCtx,
	}
}

`
