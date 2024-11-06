package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type PhotoService struct {
	svcCtx *svctx.ServiceContext
}

func NewPhotoService(svcCtx *svctx.ServiceContext) *PhotoService {
	return &PhotoService{
		svcCtx: svcCtx,
	}
}

// 分页获取照片列表
func (s *PhotoService) FindPhotoList(reqCtx *request.Context, in *dto.PhotoQuery) (out *dto.PageResp, err error) {
	// todo

	return
}

// 批量删除照片
func (s *PhotoService) BatchDeletePhoto(reqCtx *request.Context, in *dto.IdsReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 创建照片
func (s *PhotoService) AddPhoto(reqCtx *request.Context, in *dto.PhotoNewReq) (out *dto.PhotoBackDTO, err error) {
	// todo

	return
}

// 删除照片
func (s *PhotoService) DeletePhoto(reqCtx *request.Context, in *dto.IdReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 更新照片
func (s *PhotoService) UpdatePhoto(reqCtx *request.Context, in *dto.PhotoNewReq) (out *dto.PhotoBackDTO, err error) {
	// todo

	return
}
