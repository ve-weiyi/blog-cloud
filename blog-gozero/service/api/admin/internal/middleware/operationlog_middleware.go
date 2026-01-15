package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-openapi/spec"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizheader"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/jsonconv"
)

type OperationLogMiddleware struct {
	sp *spec.Swagger

	sr syslogrpc.SyslogRpc

	pr permissionrpc.PermissionRpc

	cache *collection.Cache
}

func NewOperationLogMiddleware(sp *spec.Swagger, sr syslogrpc.SyslogRpc, pr permissionrpc.PermissionRpc) *OperationLogMiddleware {
	cache, err := collection.NewCache(60 * 30)
	if err != nil {
		panic(err)
	}

	return &OperationLogMiddleware{
		sp:    sp,
		sr:    sr,
		pr:    pr,
		cache: cache,
	}
}

// 记录操作记录
func (m *OperationLogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Debugf("OperationLogMiddleware Handle path: %v", r.URL.Path)

		ok, err := m.isTraceable(r.URL.Path)
		if err != nil || !ok {
			logx.Debugf("OperationLogMiddleware Handle isTraceable err: %v", err)
			next.ServeHTTP(w, r)
			return
		}

		// 记录处理开始时间
		start := time.Now()

		// 创建一个响应记录器
		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		// 调用下一层的处理
		next.ServeHTTP(rec, r)

		var req, resp string

		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") {
			// 如果请求为 multipart/form-data 格式，解析并保存请求参数
			form := r.MultipartForm
			req = jsonconv.AnyToJsonNE(form)
		} else {
			// 否则，读取请求体，并保存为 JSON 或字符串
			body, _ := io.ReadAll(r.Body)
			// 重新设置请求体，以便后续处理中可以读取
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			req = string(body)
		}

		resp = rec.body.String()

		maxLen := 1000
		// 数据太长时，需要截取
		if len(req) > maxLen {
			req = jsonconv.AnyToJsonIndent(&req)
			req = req[:maxLen]
		}
		if len(resp) > maxLen {
			resp = jsonconv.AnyToJsonIndent(&resp)
			resp = resp[:maxLen]
		}

		//header := make(map[string][]string)
		//if len(r.Header) > 0 {
		//	for k, v := range r.Header {
		//		if len(k) == 0 {
		//			continue
		//		}
		//
		//		keyLowercase := strings.ToLower(k)
		//		for _, key := range restx.HeaderFields {
		//			if key == keyLowercase {
		//				header[key] = v
		//			}
		//		}
		//	}
		//}

		// 计算请求响应的耗时
		cost := time.Since(start)

		var module, desc string
		api := m.getApiSpec(r.URL.Path, r.Method)
		if api != nil {
			if len(api.Tags) > 0 {
				module = api.Tags[0]
			}
			desc = api.OperationProps.Summary
		}

		op := &syslogrpc.AddOperationLogReq{
			UserId:         r.Header.Get(bizheader.HeaderUid),
			TerminalId:     r.Header.Get(bizheader.HeaderXTerminalId),
			OptModule:      module,
			OptDesc:        desc,
			RequestUri:     r.URL.Path,
			RequestMethod:  r.Method,
			RequestData:    req,
			ResponseData:   resp,
			ResponseStatus: int64(rec.statusCode),
			Cost:           fmt.Sprintf("%v", cost),
		}

		_, err = m.sr.AddOperationLog(context.Background(), op)
		if err != nil {
			logx.Errorf("OperationLogMiddleware Handle AddOperationLog err: %v", err)
		}
	}
}

func (m *OperationLogMiddleware) getApiSpec(path string, method string) *spec.Operation {
	sp := m.sp
	for k, v := range sp.Paths.Paths {
		if k == path {
			switch method {
			case http.MethodGet:
				return v.Get
			case http.MethodPost:
				return v.Post
			case http.MethodPut:
				return v.Put
			case http.MethodDelete:
				return v.Delete
			case http.MethodPatch:
				return v.Patch
			case http.MethodOptions:
				return v.Options
			case http.MethodHead:
				return v.Head
			}
		}
	}

	return nil
}

func (m *OperationLogMiddleware) isTraceable(path string) (bool, error) {
	if v, ok := m.cache.Get(path); ok {
		return v.(bool), nil
	}

	// 收集资源
	apiList, err := m.pr.FindAllApi(context.Background(), &permissionrpc.FindAllApiReq{})
	if err != nil {
		return false, err
	}
	var res bool
	for _, v := range apiList.List {
		if v.Path == path {
			res = v.Traceable == 1
		}
		m.cache.Set(v.Path, v.Traceable == 1)
	}

	return res, nil
}

// responseRecorder 记录响应体和状态码
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

// WriteHeader 记录状态码
func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

// Write 记录响应体
func (rr *responseRecorder) Write(b []byte) (int, error) {
	rr.body.Write(b)
	return rr.ResponseWriter.Write(b)
}
