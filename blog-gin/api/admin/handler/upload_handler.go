package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/response"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type UploadController struct {
	svcCtx *svctx.ServiceContext
}

func NewUploadController(svcCtx *svctx.ServiceContext) *UploadController {
	return &UploadController{
		svcCtx: svcCtx,
	}
}

// @Tags		Upload
// @Summary		"删除文件列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.DeletesUploadFileReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.BatchResp}	"返回信息"
// @Router		/admin-api/v1/upload/deletes_upload_file [DELETE]
func (s *UploadController) DeletesUploadFile(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.DeletesUploadFileReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewUploadLogic(s.svcCtx).DeletesUploadFile(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Upload
// @Summary		"获取文件列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.ListUploadFileReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageResp}	"返回信息"
// @Router		/admin-api/v1/upload/list_upload_file [POST]
func (s *UploadController) ListUploadFile(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.ListUploadFileReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewUploadLogic(s.svcCtx).ListUploadFile(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Upload
// @Summary		"上传文件列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.MultiUploadFileReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=[]types.FileInfoVO}	"返回信息"
// @Router		/admin-api/v1/upload/multi_upload_file [POST]
func (s *UploadController) MultiUploadFile(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.MultiUploadFileReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewUploadLogic(s.svcCtx).MultiUploadFile(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Upload
// @Summary		"上传文件"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.UploadFileReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.FileInfoVO}	"返回信息"
// @Router		/admin-api/v1/upload/upload_file [POST]
func (s *UploadController) UploadFile(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.UploadFileReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewUploadLogic(s.svcCtx).UploadFile(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
