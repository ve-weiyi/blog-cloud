package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type RemarkLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewRemarkLogic(svcCtx *svctx.ServiceContext) *RemarkLogic {
	return &RemarkLogic{
		svcCtx: svcCtx,
	}
}

// 分页获取留言列表
func (s *RemarkLogic) FindRemarkList(reqCtx *request.Context, in *types.QueryRemarkReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 创建留言
func (s *RemarkLogic) AddRemark(reqCtx *request.Context, in *types.NewRemarkReq) (out *types.EmptyResp, err error) {
	// todo

	return
}
