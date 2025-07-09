package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type TagLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewTagLogic(svcCtx *svctx.ServiceContext) *TagLogic {
	return &TagLogic{
		svcCtx: svcCtx,
	}
}

// 创建标签
func (s *TagLogic) AddTag(reqCtx *request.Context, in *types.TagNewReq) (out *types.TagBackVO, err error) {
	// todo

	return
}

// 删除标签
func (s *TagLogic) DeletesTag(reqCtx *request.Context, in *types.IdsReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 分页获取标签列表
func (s *TagLogic) FindTagList(reqCtx *request.Context, in *types.TagQuery) (out *types.PageResp, err error) {
	// todo

	return
}

// 更新标签
func (s *TagLogic) UpdateTag(reqCtx *request.Context, in *types.TagNewReq) (out *types.TagBackVO, err error) {
	// todo

	return
}
