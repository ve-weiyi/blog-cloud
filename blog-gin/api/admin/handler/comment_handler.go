package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/response"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type CommentController struct {
	svcCtx *svctx.ServiceContext
}

func NewCommentController(svcCtx *svctx.ServiceContext) *CommentController {
	return &CommentController{
		svcCtx: svcCtx,
	}
}

// @Tags		Comment
// @Summary		"删除评论"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.IdsReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.BatchResp}	"返回信息"
// @Router		/admin-api/v1/comment/deletes_comment [DELETE]
func (s *CommentController) DeletesComment(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.IdsReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewCommentLogic(s.svcCtx).DeletesComment(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Comment
// @Summary		"查询评论列表(后台)"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.QueryCommentReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageResp}	"返回信息"
// @Router		/admin-api/v1/comment/find_comment_back_list [POST]
func (s *CommentController) FindCommentBackList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.QueryCommentReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewCommentLogic(s.svcCtx).FindCommentBackList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Comment
// @Summary		"更新评论状态"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.UpdateCommentStatusReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.BatchResp}	"返回信息"
// @Router		/admin-api/v1/comment/update_comment_status [PUT]
func (s *CommentController) UpdateCommentStatus(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.UpdateCommentStatusReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewCommentLogic(s.svcCtx).UpdateCommentStatus(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
