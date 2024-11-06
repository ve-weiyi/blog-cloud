package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type FileService struct {
	svcCtx *svctx.ServiceContext
}

func NewFileService(svcCtx *svctx.ServiceContext) *FileService {
	return &FileService{
		svcCtx: svcCtx,
	}
}

// 分页获取文件列表
func (s *FileService) FindFileList(reqCtx *request.Context, in *dto.FileQuery) (out *dto.PageResp, err error) {
	// todo

	return
}

// 创建文件目录
func (s *FileService) AddFileFolder(reqCtx *request.Context, in *dto.FileFolderNewReq) (out *dto.FileBackDTO, err error) {
	// todo

	return
}

// 删除文件列表
func (s *FileService) DeletesFile(reqCtx *request.Context, in *dto.IdsReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 上传文件列表
func (s *FileService) MultiUploadFile(reqCtx *request.Context, in *dto.MultiUploadFileReq) (out []*dto.FileBackDTO, err error) {
	// todo

	return
}

// 上传文件
func (s *FileService) UploadFile(reqCtx *request.Context, in *dto.UploadFileReq) (out *dto.FileBackDTO, err error) {
	// todo

	return
}
