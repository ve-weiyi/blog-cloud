package websocket

import (
	"net/http"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/logic/common/websocket"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/responsex"
)

// WebSocket消息
func WebsocketHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := websocket.NewWebsocketLogic(r.Context(), svcCtx)
		err := l.Websocket(w, r)
		responsex.Response(r, w, nil, err)
	}
}
