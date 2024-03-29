package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/response"
	"github.com/ve-weiyi/ve-blog-golang/server/global"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/apierr/codex"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/jsonconv"
)

// 操作日志
func OperationRecord() gin.HandlerFunc {
	permissionHolder := global.Permission

	return func(c *gin.Context) {
		// 检测接口是否需要操作记录
		permission, err := permissionHolder.FindApiPermission(c.Request.URL.Path, c.Request.Method)
		if err != nil {
			global.LOG.Error(err)
		}
		// 未加载接口权限信息，或接口未开放，或接口不需要记录操作日志
		if permission == nil {
			c.Next()
			return
		}
		if permission.Traceable == 0 {
			c.Next()
			return
		}

		start := time.Now()
		var reqData interface{}

		contentType := c.Request.Header.Get("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") {
			// 如果请求为 multipart/form-data 格式，解析并保存请求参数
			form, _ := c.MultipartForm()
			reqData = form
		} else {
			// 否则，读取请求体，并保存为 JSON 或字符串
			body, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			reqJson := make(map[string]interface{})
			if err := jsoniter.Unmarshal(body, &reqJson); err == nil {
				reqData = reqJson
			} else {
				reqData = string(body)
			}
		}

		// 替换原始的 ResponseWriter，以便在处理响应时捕获响应体内容
		respBody := bytes.NewBufferString("")
		c.Writer = &responseBodyWriter{body: respBody, ResponseWriter: c.Writer}

		// 挂起当前中间件，执行下一个中间件
		c.Next()

		// 计算请求响应的耗时
		cost := time.Since(start)
		clientIP := c.ClientIP()
		if clientIP == "" {
			clientIP = c.RemoteIP()
		}
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		if query != "" {
			path = path + "?" + query
		}

		var respData interface{}
		respParam := make(map[string]interface{})
		// 尝试将响应体解析为 JSON，并保存为 map[string]interface{} 或字符串
		if err := jsoniter.Unmarshal(respBody.Bytes(), &respParam); err == nil {
			respData = jsonconv.ObjectToJson(respParam)
		} else {
			respData = respBody.String()
		}

		var req, resp string
		req = jsonconv.ObjectToJson(reqData)
		resp = jsonconv.ObjectToJson(respData)

		// 数据太长时，需要截取
		if len(req) > 4000 {
			req = jsonconv.ObjectToJsonIndent(&response.Response{})
		}
		if len(resp) > 4000 {
			resp = jsonconv.ObjectToJsonIndent(&response.Response{})
		}

		op := entity.OperationLog{
			ID:            0,
			UserID:        cast.ToInt(c.GetString("uid")),
			Nickname:      c.GetString("username"),
			IpAddress:     c.GetString("ip_address"),
			IpSource:      c.GetString("ip_source"),
			OptModule:     permission.Group,
			OptDesc:       permission.Name,
			RequestURL:    c.Request.URL.String(),
			RequestMethod: c.Request.Method,
			// 请求头携带token，数据太多
			//RequestHeader: jsonconv.ObjectToJson(c.Request.Header),
			RequestData:    req,
			ResponseData:   resp,
			ResponseStatus: c.Writer.Status(),
			Cost:           fmt.Sprintf("%v", cost),
			CreatedAt:      time.Now(),
		}
		err = global.DB.Create(&op).Error
		if err != nil {
			global.LOG.Error(err)
			c.JSON(http.StatusOK, response.Response{
				Code:    codex.CodeInternalServerError,
				Message: "日志记录错误",
				Data:    nil,
			})
			c.Abort()
			return
		}
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// 获取 handler 处理函数的名称
func getHandlerFunc(handler gin.HandlerFunc) *runtime.Func {
	handlerFunc := reflect.ValueOf(handler).Pointer()
	return runtime.FuncForPC(handlerFunc)
}
