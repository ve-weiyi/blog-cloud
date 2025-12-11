package websocket

import (
	"net/http"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/responsex"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/logic/websocket"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
)

// WebSocket消息
func WebsocketHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 接收消息
		l := websocket.NewWebsocketLogic(r.Context(), svcCtx)
		err := l.Websocket(w, r)
		responsex.Response(r, w, nil, err)
	}
}
