package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type AlbumService struct {
	svcCtx *svctx.ServiceContext
}

func NewAlbumService(svcCtx *svctx.ServiceContext) *AlbumService {
	return &AlbumService{
		svcCtx: svcCtx,
	}
}

// 获取相册列表
func (s *AlbumService) FindAlbumList(reqCtx *request.Context, in *dto.AlbumQueryReq) (out *dto.PageResp, err error) {
	// todo

	return
}

// 获取相册下的照片列表
func (s *AlbumService) FindPhotoList(reqCtx *request.Context, in *dto.PhotoQueryReq) (out *dto.PageResp, err error) {
	// todo

	return
}

// 获取相册
func (s *AlbumService) GetAlbum(reqCtx *request.Context, in *dto.IdReq) (out *dto.Album, err error) {
	// todo

	return
}
