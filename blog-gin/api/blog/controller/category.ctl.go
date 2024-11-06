package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/dto"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/service"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/base/response"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type CategoryController struct {
	svcCtx *svctx.ServiceContext
}

func NewCategoryController(svcCtx *svctx.ServiceContext) *CategoryController {
	return &CategoryController{
		svcCtx: svcCtx,
	}
}

// @Tags		Category
// @Summary		"分页获取文章分类列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.CategoryQueryReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/api/v1/category/find_category_list [POST]
func (s *CategoryController) FindCategoryList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.CategoryQueryReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewCategoryService(s.svcCtx).FindCategoryList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}
