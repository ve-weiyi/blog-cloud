// Code scaffolded by goctl. Safe to edit.
package notify_stream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/threading"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/logic/notification/notify_stream"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
)

// 通知消息推送流（SSE）
func NotifyStreamHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// SSE 长连接经 nginx 反代时防止响应被缓冲（没有这个头，客户端会等到缓冲填满才收到帧）
		w.Header().Set("X-Accel-Buffering", "no")

		// Buffer size of 16 is chosen as a reasonable default to balance throughput and memory usage.
		// You can change this based on your application's needs.
		// if your go-zero version less than 1.8.1, you need to add 3 lines below.
		// w.Header().Set("Content-Type", "text/event-stream")
		// w.Header().Set("Cache-Control", "no-cache")
		// w.Header().Set("Connection", "keep-alive")
		client := make(chan *types.NotifyStreamEvent, 16)

		l := notify_stream.NewNotifyStreamLogic(r.Context(), svcCtx)
		threading.GoSafeCtx(r.Context(), func() {
			defer close(client)
			err := l.NotifyStream(client)
			if err != nil {
				logc.Errorw(r.Context(), "NotifyStreamHandler", logc.Field("error", err))
				return
			}
		})

		heartbeat := time.NewTicker(HeartbeatInterval)
		defer heartbeat.Stop()

		for {
			select {
			case data, ok := <-client:
				if !ok {
					return
				}
				output, err := json.Marshal(data)
				if err != nil {
					logc.Errorw(r.Context(), "NotifyStreamHandler", logc.Field("error", err))
					continue
				}

				if _, err := fmt.Fprintf(w, "data: %s\n\n", string(output)); err != nil {
					logc.Errorw(r.Context(), "NotifyStreamHandler", logc.Field("error", err))
					return
				}
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
			case <-heartbeat.C:
				if err := beat(w); err != nil {
					return
				}
			case <-r.Context().Done():
				return
			}
		}
	}
}
