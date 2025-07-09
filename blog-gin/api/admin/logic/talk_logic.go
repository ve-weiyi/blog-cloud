package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type TalkLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewTalkLogic(svcCtx *svctx.ServiceContext) *TalkLogic {
	return &TalkLogic{
		svcCtx: svcCtx,
	}
}

// 创建说说
func (s *TalkLogic) AddTalk(reqCtx *request.Context, in *types.TalkNewReq) (out *types.TalkBackVO, err error) {
	// todo

	return
}

// 删除说说
func (s *TalkLogic) DeleteTalk(reqCtx *request.Context, in *types.IdReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 分页获取说说列表
func (s *TalkLogic) FindTalkList(reqCtx *request.Context, in *types.TalkQuery) (out *types.PageResp, err error) {
	// todo

	return
}

// 查询说说
func (s *TalkLogic) GetTalk(reqCtx *request.Context, in *types.IdReq) (out *types.TalkBackVO, err error) {
	// todo

	return
}

// 更新说说
func (s *TalkLogic) UpdateTalk(reqCtx *request.Context, in *types.TalkNewReq) (out *types.TalkBackVO, err error) {
	// todo

	return
}
