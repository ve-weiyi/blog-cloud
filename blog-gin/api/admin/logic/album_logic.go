package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
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

// 创建相册
func (s *AlbumLogic) AddAlbum(reqCtx *request.Context, in *types.NewAlbumReq) (out *types.AlbumBackVO, err error) {
	// todo

	return
}

// 删除相册
func (s *AlbumLogic) DeletesAlbum(reqCtx *request.Context, in *types.IdsReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 分页获取相册列表
func (s *AlbumLogic) FindAlbumList(reqCtx *request.Context, in *types.QueryAlbumReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 查询相册
func (s *AlbumLogic) GetAlbum(reqCtx *request.Context, in *types.IdReq) (out *types.AlbumBackVO, err error) {
	// todo

	return
}

// 更新相册
func (s *AlbumLogic) UpdateAlbum(reqCtx *request.Context, in *types.NewAlbumReq) (out *types.AlbumBackVO, err error) {
	// todo

	return
}

// 更新相册删除状态
func (s *AlbumLogic) UpdateAlbumDelete(reqCtx *request.Context, in *types.UpdateAlbumDeleteReq) (out *types.BatchResp, err error) {
	// todo

	return
}
