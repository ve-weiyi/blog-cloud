package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type UploadLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewUploadLogic(svcCtx *svctx.ServiceContext) *UploadLogic {
	return &UploadLogic{
		svcCtx: svcCtx,
	}
}

// 删除文件列表
func (s *UploadLogic) DeletesUploadFile(reqCtx *request.Context, in *types.DeletesUploadFileReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 获取文件列表
func (s *UploadLogic) ListUploadFile(reqCtx *request.Context, in *types.ListUploadFileReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 上传文件列表
func (s *UploadLogic) MultiUploadFile(reqCtx *request.Context, in *types.MultiUploadFileReq) (out []*types.FileInfoVO, err error) {
	// todo

	return
}

// 上传文件
func (s *UploadLogic) UploadFile(reqCtx *request.Context, in *types.UploadFileReq) (out *types.FileInfoVO, err error) {
	// todo

	return
}
