package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type WebsocketLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewWebsocketLogic(svcCtx *svctx.ServiceContext) *WebsocketLogic {
	return &WebsocketLogic{
		svcCtx: svcCtx,
	}
}

// WebSocket消息
func (s *WebsocketLogic) Websocket(reqCtx *request.Context) (err error) {
	// todo

	return
}
