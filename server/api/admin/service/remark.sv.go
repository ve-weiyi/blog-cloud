package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type RemarkService struct {
	svcCtx *svctx.ServiceContext
}

func NewRemarkService(svcCtx *svctx.ServiceContext) *RemarkService {
	return &RemarkService{
		svcCtx: svcCtx,
	}
}

// 分页获取留言列表
func (s *RemarkService) FindRemarkList(reqCtx *request.Context, in *dto.RemarkQuery) (out *dto.PageResp, err error) {
	// todo

	return
}

// 批量删除留言
func (s *RemarkService) BatchDeleteRemark(reqCtx *request.Context, in *dto.IdsReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 删除留言
func (s *RemarkService) DeleteRemark(reqCtx *request.Context, in *dto.IdReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 更新留言
func (s *RemarkService) UpdateRemark(reqCtx *request.Context, in *dto.RemarkNewReq) (out *dto.RemarkBackDTO, err error) {
	// todo

	return
}
