package middleware

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-openapi/spec"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/middleware/tracelogx"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/syslogservice"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
	"github.com/ve-weiyi/blog-cloud/infra/http/routex"
	"github.com/ve-weiyi/vkit/x/jsonv"
	"github.com/ve-weiyi/vkit/x/redactx"
	"github.com/ve-weiyi/vkit/x/strx"
)

const (
	maxLogRequest  = 64 << 10
	maxLogResponse = 64 << 10
)

type (
	OperationLogMiddleware struct {
		specs  map[string][]apiSpecEntry
		holder tracelogx.Enforcer
		sr     syslogservice.SyslogService
	}

	apiSpecEntry struct {
		path string
		op   *spec.Operation
	}
)

func NewOperationLogMiddleware(sp *spec.Swagger, holder tracelogx.Enforcer, sr syslogservice.SyslogService) *OperationLogMiddleware {
	m := &OperationLogMiddleware{
		specs:  make(map[string][]apiSpecEntry),
		holder: holder,
		sr:     sr,
	}
	if sp == nil {
		return m
	}

	for path, item := range sp.Paths.Paths {
		normalized := routex.Normalize(path)
		for method, op := range map[string]*spec.Operation{
			http.MethodGet:     item.Get,
			http.MethodPost:    item.Post,
			http.MethodPut:     item.Put,
			http.MethodDelete:  item.Delete,
			http.MethodPatch:   item.Patch,
			http.MethodOptions: item.Options,
			http.MethodHead:    item.Head,
		} {
			if op != nil {
				m.specs[method] = append(m.specs[method], apiSpecEntry{path: normalized, op: op})
			}
		}
	}

	return m
}

func (m *OperationLogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok, err := m.holder.IsTraceLog(r.URL.Path, r.Method)
		if err != nil || !ok {
			logx.Debugf("OperationLogMiddleware skip path=%s err=%v", r.URL.Path, err)
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		reqBody := &cappedBuffer{limit: maxLogRequest}
		captureBody := strings.Contains(r.Header.Get("Content-Type"), "application/json")
		if captureBody && r.Body != nil {
			r.Body = io.NopCloser(io.TeeReader(r.Body, reqBody))
		}

		rec := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			body:           &cappedBuffer{limit: maxLogResponse},
		}
		next.ServeHTTP(rec, r)

		requestData := reqBody.String()
		if r.MultipartForm != nil && len(r.MultipartForm.Value) > 0 {
			requestData = jsonv.AnyToJsonNE(r.MultipartForm.Value)
		}
		requestData = redactx.JSON(requestData)
		responseData := redactx.JSON(rec.body.String())

		// 数据太长时按字节上限截断；截断必须落在字符边界上，
		// 否则审计日志里会存进非法 UTF-8
		requestData = strx.TruncateUTF8(requestData, maxLogRequest)
		responseData = strx.TruncateUTF8(responseData, maxLogResponse)

		var module, desc string
		if op := m.getApiSpec(r.URL.Path, r.Method); op != nil {
			if len(op.Tags) > 0 {
				module = op.Tags[0]
			}
			desc = op.OperationProps.Summary
		}

		// 身份取自已认证的上下文，客户端传入的头不可信。
		uid, _ := metax.GetApiUserIdFromCtx(r.Context())
		deviceID, _ := metax.GetApiDeviceIdFromCtx(r.Context())
		op := &syslogservice.CreateOperationLogRequest{
			UserId:         uid,
			DeviceId:       deviceID,
			Module:         module,
			Description:    desc,
			RequestUri:     r.URL.Path,
			RequestMethod:  r.Method,
			RequestData:    requestData,
			ResponseData:   responseData,
			ResponseStatus: int64(rec.statusCode),
			Cost:           time.Since(start).Milliseconds(),
		}

		if _, err = m.sr.CreateOperationLog(r.Context(), op); err != nil {
			logx.Errorf("OperationLogMiddleware CreateOperationLog err: %v", err)
		}
	}
}

func (m *OperationLogMiddleware) getApiSpec(path, method string) *spec.Operation {
	path = routex.Normalize(path)
	for _, entry := range m.specs[strings.ToUpper(method)] {
		if routex.Match(entry.path, path) {
			return entry.op
		}
	}
	return nil
}

// cappedBuffer 只保留前 limit 字节，超出部分丢弃但继续接收写入。
type cappedBuffer struct {
	buf   bytes.Buffer
	limit int
}

func (b *cappedBuffer) Write(p []byte) (int, error) {
	length := len(p)
	if b.limit <= 0 {
		return length, nil
	}
	if remaining := b.limit - b.buf.Len(); remaining > 0 {
		if len(p) > remaining {
			p = p[:remaining]
		}
		b.buf.Write(p)
	}
	return length, nil
}

func (b *cappedBuffer) String() string {
	return b.buf.String()
}

// responseRecorder 记录响应体和状态码。
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *cappedBuffer
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	rr.body.Write(b)
	return rr.ResponseWriter.Write(b)
}

func (rr *responseRecorder) Flush() {
	if flusher, ok := rr.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (rr *responseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rr.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return hijacker.Hijack()
}
