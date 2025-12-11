package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/response"
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
// @Summary		"查询评论列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.CommentQueryReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageResp}	"返回信息"
// @Router		/blog-api/v1/comment/find_comment_list [POST]
func (s *CommentController) FindCommentList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.CommentQueryReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewCommentLogic(s.svcCtx).FindCommentList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Comment
// @Summary		"查询最新评论回复列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.CommentQueryReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageResp}	"返回信息"
// @Router		/blog-api/v1/comment/find_comment_recent_list [POST]
func (s *CommentController) FindCommentRecentList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.CommentQueryReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewCommentLogic(s.svcCtx).FindCommentRecentList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Comment
// @Summary		"查询评论回复列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.CommentQueryReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageResp}	"返回信息"
// @Router		/blog-api/v1/comment/find_comment_reply_list [POST]
func (s *CommentController) FindCommentReplyList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.CommentQueryReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewCommentLogic(s.svcCtx).FindCommentReplyList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Comment
// @Summary		"创建评论"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.CommentNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.Comment}	"返回信息"
// @Router		/blog-api/v1/comment/add_comment [POST]
func (s *CommentController) AddComment(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.CommentNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewCommentLogic(s.svcCtx).AddComment(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Comment
// @Summary		"点赞评论"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/blog-api/v1/comment/like_comment [POST]
func (s *CommentController) LikeComment(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.IdReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewCommentLogic(s.svcCtx).LikeComment(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Comment
// @Summary		"更新评论"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.UpdateCommentReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.Comment}	"返回信息"
// @Router		/blog-api/v1/comment/update_comment [POST]
func (s *CommentController) UpdateComment(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.UpdateCommentReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewCommentLogic(s.svcCtx).UpdateComment(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
