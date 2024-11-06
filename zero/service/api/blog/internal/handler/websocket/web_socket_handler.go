package websocket

import (
	"net/http"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/logic/websocket"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
)

// WebSocket消息
func WebSocketHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 接收消息
		l := websocket.NewWebSocketLogic(r.Context(), svcCtx)
		l.WebSocket(w, r)
	}
}
