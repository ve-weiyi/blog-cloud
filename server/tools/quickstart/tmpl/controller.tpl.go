package tmpl

const AppController = `
package controller

import (

)

type AppController struct {
	svcCtx *svc.ControllerContext //持有的service引用
}

func NewController(svcCtx *svc.ControllerContext) *AppController {
	return &AppController{
		svcCtx: svcCtx,
	}
}
`

const ControllerContext = `
package svc

import (

)

// 注册需要用到的rpc
type ControllerContext struct {
	*service.AppService
}

func NewControllerContext(cfg *config.Config) *ControllerContext {
	ctx := svc.NewServiceContext(cfg)
	sv := service.NewService(ctx)
	if sv == nil {
		panic("sv cannot be null")
	}

	return &ControllerContext{
		AppService: sv,
	}
}

`

const Controller = `
package logic

import (
	"github.com/gin-gonic/gin"
	{{range .ImportPkgPaths}}{{.}} ` + "\n" + `{{end}}
)

type {{.UpperStartCamelName}}Controller struct {
	controller.BaseController
	svcCtx *svc.ControllerContext
}

func New{{.UpperStartCamelName}}Controller(svcCtx *svc.ControllerContext) *{{.UpperStartCamelName}}Controller {
	return &{{.UpperStartCamelName}}Controller{
		svcCtx: svcCtx,
		BaseController: controller.NewBaseController(svcCtx),
	}
}

// @Tags		{{.UpperStartCamelName}}
// @Summary		创建{{.CommentName}}
// @Accept		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param		data	body		entity.{{.UpperStartCamelName}}		true	"请求参数"
// @Success		200		{object}	response.Response{data=entity.{{.UpperStartCamelName}}}	"返回信息"
// @Router		/{{.SnakeName}} [post]
func (s *{{.UpperStartCamelName}}Controller) Create{{.UpperStartCamelName}}(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var {{.LowerStartCamelName}} entity.{{.UpperStartCamelName}}
	err = s.ShouldBind(c, &{{.LowerStartCamelName}})
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.{{.UpperStartCamelName}}Service.Create{{.UpperStartCamelName}}(reqCtx, &{{.LowerStartCamelName}});
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags 	 	{{.UpperStartCamelName}}
// @Summary		更新{{.CommentName}}
// @Accept 		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param 	 	data	body 	 	entity.{{.UpperStartCamelName}}		true	"请求参数"
// @Success		200		{object}	response.Response{data=entity.{{.UpperStartCamelName}}}	"返回信息"
// @Router 		/{{.SnakeName}} [put]
func (s *{{.UpperStartCamelName}}Controller) Update{{.UpperStartCamelName}}(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var {{.LowerStartCamelName}} entity.{{.UpperStartCamelName}}
	err = s.ShouldBind(c, &{{.LowerStartCamelName}})
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.{{.UpperStartCamelName}}Service.Update{{.UpperStartCamelName}}(reqCtx, &{{.LowerStartCamelName}});
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags		{{.UpperStartCamelName}}
// @Summary		删除{{.CommentName}}
// @Accept		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param 	 	id		path		int							true	"{{.UpperStartCamelName}}.id"
// @Success		200		{object}	response.Response{data=any}			"返回信息"
// @Router		/{{.SnakeName}}/{id} [delete]
func (s *{{.UpperStartCamelName}}Controller) Delete{{.UpperStartCamelName}}(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var id int
	id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.{{.UpperStartCamelName}}Service.Delete{{.UpperStartCamelName}}(reqCtx, id);
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags 	 	{{.UpperStartCamelName}}
// @Summary		查询{{.CommentName}}
// @Accept 		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param 	 	id		path		int							true	"{{.UpperStartCamelName}}.id"
// @Success		200		{object}	response.Response{data=entity.{{.UpperStartCamelName}}}	"返回信息"
// @Router 		/{{.SnakeName}}/{id} [get]
func (s *{{.UpperStartCamelName}}Controller) Find{{.UpperStartCamelName}}(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var id int
	id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.{{.UpperStartCamelName}}Service.Find{{.UpperStartCamelName}}(reqCtx, id);
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags 	 	{{.UpperStartCamelName}}
// @Summary		批量删除{{.CommentName}}
// @Accept 	 	application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param		data 	body		[]int 						true 	"删除id列表"
// @Success		200		{object}	response.Response{data=response.BatchResult}	"返回信息"
// @Router		/{{.SnakeName}}/batch_delete [delete]
func (s *{{.UpperStartCamelName}}Controller) Delete{{.UpperStartCamelName}}ByIds(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var ids []int
	err = s.ShouldBind(c, &ids)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.{{.UpperStartCamelName}}Service.Delete{{.UpperStartCamelName}}ByIds(reqCtx, ids);
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, response.BatchResult{
		SuccessCount: data,
	})
}

// @Tags 	 	{{.UpperStartCamelName}}
// @Summary		分页获取{{.CommentName}}列表
// @Accept 		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param 	 	page 	body		request.PageQuery 			true 	"分页参数"
// @Success		200		{object}	response.Response{data=response.PageResult{list=[]entity.{{.UpperStartCamelName}}}}	"返回信息"
// @Router		/{{.SnakeName}}/list [post]
func (s *{{.UpperStartCamelName}}Controller) Find{{.UpperStartCamelName}}List(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var page request.PageQuery
	err = s.ShouldBind(c, &page)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	list, total, err := s.svcCtx.{{.UpperStartCamelName}}Service.Find{{.UpperStartCamelName}}List(reqCtx, &page); 
	if err != nil {
		s.ResponseError(c, err)
		return
	}
	
	s.ResponseOk(c, response.PageResult{
		List:     list,
		Total:    total,
		Page:     page.Page,
		PageSize: page.PageSize,
	})
}
`

const CommonController = `
package logic

import (

)

type {{.UpperStartCamelName}}Controller struct {
	controller.BaseController
	svcCtx *svc.ControllerContext
}

func New{{.UpperStartCamelName}}Controller(svcCtx *svc.ControllerContext) *{{.UpperStartCamelName}}Controller {
	return &{{.UpperStartCamelName}}Controller{
		svcCtx: svcCtx,
		BaseController: controller.NewBaseController(svcCtx),
	}
}
`
