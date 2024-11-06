package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type TalkService struct {
	svcCtx *svctx.ServiceContext
}

func NewTalkService(svcCtx *svctx.ServiceContext) *TalkService {
	return &TalkService{
		svcCtx: svcCtx,
	}
}

// 分页获取说说列表
func (s *TalkService) FindTalkList(reqCtx *request.Context, in *dto.TalkQueryReq) (out *dto.PageResp, err error) {
	// todo

	return
}

// 查询说说
func (s *TalkService) GetTalk(reqCtx *request.Context, in *dto.IdReq) (out *dto.Talk, err error) {
	// todo

	return
}

// 点赞说说
func (s *TalkService) LikeTalk(reqCtx *request.Context, in *dto.IdReq) (out *dto.EmptyResp, err error) {
	// todo

	return
}
