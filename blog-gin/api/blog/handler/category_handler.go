package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/response"
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
// @Param		data	body		types.CategoryQueryReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageResp}	"返回信息"
// @Router		/blog-api/v1/category/find_category_list [POST]
func (s *CategoryController) FindCategoryList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.CategoryQueryReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewCategoryLogic(s.svcCtx).FindCategoryList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
