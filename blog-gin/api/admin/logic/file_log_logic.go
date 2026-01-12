package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type FileLogLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewFileLogLogic(svcCtx *svctx.ServiceContext) *FileLogLogic {
	return &FileLogLogic{
		svcCtx: svcCtx,
	}
}

// 删除文件日志
func (s *FileLogLogic) DeletesFileLog(reqCtx *request.Context, in *types.IdsReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 查询文件日志
func (s *FileLogLogic) FindFileLogList(reqCtx *request.Context, in *types.QueryFileLogReq) (out *types.PageResp, err error) {
	// todo

	return
}
