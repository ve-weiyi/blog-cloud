package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type PhotoLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewPhotoLogic(svcCtx *svctx.ServiceContext) *PhotoLogic {
	return &PhotoLogic{
		svcCtx: svcCtx,
	}
}

// 创建照片
func (s *PhotoLogic) AddPhoto(reqCtx *request.Context, in *types.PhotoNewReq) (out *types.PhotoBackVO, err error) {
	// todo

	return
}

// 删除照片
func (s *PhotoLogic) DeletesPhoto(reqCtx *request.Context, in *types.IdsReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 分页获取照片列表
func (s *PhotoLogic) FindPhotoList(reqCtx *request.Context, in *types.PhotoQuery) (out *types.PageResp, err error) {
	// todo

	return
}

// 预删除照片
func (s *PhotoLogic) PreDeletePhoto(reqCtx *request.Context, in *types.PreDeletePhotoReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 更新照片
func (s *PhotoLogic) UpdatePhoto(reqCtx *request.Context, in *types.PhotoNewReq) (out *types.PhotoBackVO, err error) {
	// todo

	return
}
