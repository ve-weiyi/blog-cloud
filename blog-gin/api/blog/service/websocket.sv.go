package service

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type WebsocketService struct {
	svcCtx *svctx.ServiceContext
}

func NewWebsocketService(svcCtx *svctx.ServiceContext) *WebsocketService {
	return &WebsocketService{
		svcCtx: svcCtx,
	}
}

// WebSocket消息
func (s *WebsocketService) WebSocket(reqCtx *request.Context) (err error) {
	// todo

	return
}
