package websocket

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
)

type WebsocketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// WebSocket消息
func NewWebsocketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebsocketLogic {
	return &WebsocketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WebsocketLogic) Websocket(w http.ResponseWriter, r *http.Request) error {
	l.svcCtx.StompHubServer.HandleWebSocket(w, r)
	return nil
}
