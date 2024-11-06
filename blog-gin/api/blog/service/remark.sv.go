package service

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/dto"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
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
func (s *RemarkService) FindRemarkList(reqCtx *request.Context, in *dto.RemarkQueryReq) (out *dto.PageResp, err error) {
	// todo

	return
}

// 创建留言
func (s *RemarkService) AddRemark(reqCtx *request.Context, in *dto.RemarkNewReq) (out *dto.Remark, err error) {
	// todo

	return
}
