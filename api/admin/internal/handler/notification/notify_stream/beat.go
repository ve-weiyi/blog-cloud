// Code scaffolded by goctl. Safe to edit.
package notify_stream

import (
	"fmt"
	"net/http"
	"time"
)

// HeartbeatInterval 心跳间隔。推送流可能长时间无数据，固定间隔写一条 SSE
// 注释行保活，防止中间代理按空闲超时切断连接。
//
// 前端据此判定链路是否已死（静默超时 = 本值的 3 倍，见 useSse.ts 的 staleTimeout）：
// **调小本值不会立即生效**——前端阈值没跟着改就会把正常连接误判为断开。
const HeartbeatInterval = 15 * time.Second

// beat 写一条 SSE 注释行。注释行不触发客户端事件分发，仅用于保活。
func beat(w http.ResponseWriter) error {
	if _, err := fmt.Fprint(w, ": ping\n\n"); err != nil {
		return err
	}
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
	return nil
}
