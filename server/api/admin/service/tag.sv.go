package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type TagService struct {
	svcCtx *svctx.ServiceContext
}

func NewTagService(svcCtx *svctx.ServiceContext) *TagService {
	return &TagService{
		svcCtx: svcCtx,
	}
}

// 分页获取标签列表
func (s *TagService) FindTagList(reqCtx *request.Context, in *dto.TagQuery) (out *dto.PageResp, err error) {
	// todo

	return
}

// 创建标签
func (s *TagService) AddTag(reqCtx *request.Context, in *dto.TagNewReq) (out *dto.TagBackDTO, err error) {
	// todo

	return
}

// 批量删除标签
func (s *TagService) BatchDeleteTag(reqCtx *request.Context, in *dto.IdsReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 删除标签
func (s *TagService) DeleteTag(reqCtx *request.Context, in *dto.IdReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 更新标签
func (s *TagService) UpdateTag(reqCtx *request.Context, in *dto.TagNewReq) (out *dto.TagBackDTO, err error) {
	// todo

	return
}
