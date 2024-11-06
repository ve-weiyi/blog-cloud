package middlewarex

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/apierr"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/apierr/codex"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/constant"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/ipx"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"

	"github.com/ve-weiyi/ve-blog-golang/zero/internal/rbacx"
	"github.com/ve-weiyi/ve-blog-golang/zero/internal/responsex"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/syslogrpc"
)

type OperationMiddleware struct {
	holder *rbacx.RbacHolder

	sr syslogrpc.SyslogRpc
}

func NewOperationMiddleware(holder *rbacx.RbacHolder, sr syslogrpc.SyslogRpc) *OperationMiddleware {
	return &OperationMiddleware{
		holder: holder,
		sr:     sr,
	}
}

func (m *OperationMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pl := m.holder.GetPolicy(r.URL.Path, r.Method)
		logx.Infof("OperationMiddleware Handle path: %v,pl: %v", r.URL.Path, pl)
		if pl != nil && pl.Disable() {
			responsex.Response(r, w, nil, apierr.NewApiError(codex.CodeUserNotExist, "资源已禁用"))
			return
		}

		roles := r.Context().Value("roles")
		if pl != nil && !pl.HasPermission(strings.Split(cast.ToString(roles), ",")) {
			responsex.Response(r, w, nil, apierr.NewApiError(codex.CodeUserNotExist, "无权限访问"))
			return
		}

		// 记录处理开始时间
		start := time.Now()

		// 创建一个响应记录器
		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		// 调用下一层的处理
		next.ServeHTTP(rec, r)

		if pl != nil && pl.Traceable() {
			ip := r.RemoteAddr
			is, err := ipx.GetIpSourceByBaidu(ip)
			if err != nil {
				logx.Errorf("OperationMiddleware Handle GetIpInfoByBaidu err: %v", err)
			}

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

			header := make(map[string][]string)
			if len(r.Header) > 0 {
				for k, v := range r.Header {
					if len(k) == 0 {
						continue
					}

					keyLowercase := strings.ToLower(k)
					for _, key := range constant.HeaderFields {
						if key == keyLowercase {
							header[key] = v
						}
					}
				}
			}

			// 计算请求响应的耗时
			cost := time.Since(start)

			op := &syslogrpc.OperationLogNewReq{
				UserId:         cast.ToString(r.Header[constant.HeaderUid]),
				Nickname:       "",
				IpAddress:      ip,
				IpSource:       is,
				OptModule:      pl.Module(),
				OptDesc:        pl.Desc(),
				RequestUrl:     r.URL.Path,
				RequestMethod:  r.Method,
				RequestHeader:  jsonconv.AnyToJsonNE(header),
				RequestData:    req,
				ResponseData:   resp,
				ResponseStatus: int64(rec.statusCode),
				Cost:           fmt.Sprintf("%v", cost),
			}

			_, err = m.sr.AddOperationLog(context.Background(), op)
			if err != nil {
				logx.Errorf("OperationMiddleware Handle AddOperationLog err: %v", err)
			}
		}
	}
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
