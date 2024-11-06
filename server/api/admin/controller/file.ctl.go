package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/service"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/response"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type FileController struct {
	svcCtx *svctx.ServiceContext
}

func NewFileController(svcCtx *svctx.ServiceContext) *FileController {
	return &FileController{
		svcCtx: svcCtx,
	}
}

// @Tags		File
// @Summary		"分页获取文件列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.FileQuery		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/admin_api/v1/file/find_file_list [POST]
func (s *FileController) FindFileList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.FileQuery
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewFileService(s.svcCtx).FindFileList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		File
// @Summary		"创建文件目录"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.FileFolderNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.FileBackDTO}	"返回信息"
// @Router		/admin_api/v1/file/add_file_folder [POST]
func (s *FileController) AddFileFolder(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.FileFolderNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewFileService(s.svcCtx).AddFileFolder(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		File
// @Summary		"删除文件列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.IdsReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.BatchResp}	"返回信息"
// @Router		/admin_api/v1/file/deletes_file [DELETE]
func (s *FileController) DeletesFile(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.IdsReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewFileService(s.svcCtx).DeletesFile(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		File
// @Summary		"上传文件列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.MultiUploadFileReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=[]dto.FileBackDTO}	"返回信息"
// @Router		/admin_api/v1/file/multi_upload_file [POST]
func (s *FileController) MultiUploadFile(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.MultiUploadFileReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewFileService(s.svcCtx).MultiUploadFile(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		File
// @Summary		"上传文件"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.UploadFileReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.FileBackDTO}	"返回信息"
// @Router		/admin_api/v1/file/upload_file [POST]
func (s *FileController) UploadFile(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.UploadFileReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewFileService(s.svcCtx).UploadFile(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}
