package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type AlbumLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewAlbumLogic(svcCtx *svctx.ServiceContext) *AlbumLogic {
	return &AlbumLogic{
		svcCtx: svcCtx,
	}
}

// 获取相册列表
func (s *AlbumLogic) FindAlbumList(reqCtx *request.Context, in *types.AlbumQueryReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 获取相册下的照片列表
func (s *AlbumLogic) FindPhotoList(reqCtx *request.Context, in *types.PhotoQueryReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 获取相册
func (s *AlbumLogic) GetAlbum(reqCtx *request.Context, in *types.IdReq) (out *types.Album, err error) {
	// todo

	return
}
