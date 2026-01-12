package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
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

// 分页获取说说列表
func (s *TalkLogic) FindTalkList(reqCtx *request.Context, in *types.QueryTalkReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 查询说说
func (s *TalkLogic) GetTalk(reqCtx *request.Context, in *types.IdReq) (out *types.Talk, err error) {
	// todo

	return
}

// 点赞说说
func (s *TalkLogic) LikeTalk(reqCtx *request.Context, in *types.IdReq) (out *types.EmptyResp, err error) {
	// todo

	return
}
