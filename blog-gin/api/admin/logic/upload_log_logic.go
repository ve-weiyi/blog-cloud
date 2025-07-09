package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type UploadLogLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewUploadLogLogic(svcCtx *svctx.ServiceContext) *UploadLogLogic {
	return &UploadLogLogic{
		svcCtx: svcCtx,
	}
}

// 删除登录日志
func (s *UploadLogLogic) DeletesUploadLog(reqCtx *request.Context, in *types.IdsReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 查询登录日志
func (s *UploadLogLogic) FindUploadLogList(reqCtx *request.Context, in *types.UploadLogQuery) (out *types.PageResp, err error) {
	// todo

	return
}
