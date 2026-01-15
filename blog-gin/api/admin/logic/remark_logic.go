package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
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

// 删除留言
func (s *RemarkLogic) DeletesRemark(reqCtx *request.Context, in *types.IdsReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 分页获取留言列表
func (s *RemarkLogic) FindRemarkList(reqCtx *request.Context, in *types.QueryRemarkReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 更新留言状态
func (s *RemarkLogic) UpdateRemarkStatus(reqCtx *request.Context, in *types.UpdateRemarkStatusReq) (out *types.BatchResp, err error) {
	// todo

	return
}
