package request

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/apierr"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/apierr/codex"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/constant"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
)

// 请求上下文,一般存放请求头参数
type Context struct {
	context.Context `json:"-" header:"-"`
	Token           string `json:"token" header:"token" example:""`
	Uid             int64  `json:"uid" header:"-" example:""`
	IpAddress       string `json:"ip_address" header:"-" example:""`
	UserAgent       string `json:"user_agent" header:"-" example:""`
}

func (s *Context) GetContext() context.Context {
	return s.Context
}

// 获取请求上下文
func ParseRequestContext(c *gin.Context) (*Context, error) {
	reqCtx := &Context{}
	reqCtx.Token = c.GetHeader(constant.HeaderToken)
	reqCtx.Uid = cast.ToInt64(c.GetHeader(constant.HeaderUid))
	reqCtx.IpAddress = c.ClientIP()
	reqCtx.UserAgent = c.Request.UserAgent()
	reqCtx.Context = c.Request.Context()
	return reqCtx, nil
}

type IdReq struct {
	Id int64 `json:"id" form:"id" binding:"required"`
}

type IdsReq struct {
	Ids []int64 `json:"ids" form:"ids" binding:"required"`
}

type EmptyReq struct{}

// IP限流

type IsValidChecker interface {
	IsValid() error
}

func ShouldBindJSON(c *gin.Context, req interface{}) error {
	//value := reflect.ValueOf(req)
	//if value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct {
	//	if err := BindJSONIgnoreCase(c, req); err != nil {
	//		return apierror.NewApiError(apierror.CodeMissingParameter, "参数错误").WrapMessage(err.Error())
	//	}
	//}

	if err := c.ShouldBindJSON(&req); err != nil {
		return apierr.NewApiError(codex.CodeInvalidParam, err.Error())
	}

	isValid, ok := req.(IsValidChecker)
	if !ok {
		return nil
	}

	if err := isValid.IsValid(); err != nil {
		return apierr.NewApiError(codex.CodeInvalidParam, err.Error())
	}

	return nil
}

// 把请求参数转换为小写
func BindJSONIgnoreCase(c *gin.Context, req interface{}) (err error) {

	var tmp map[string]interface{}
	err = c.ShouldBindJSON(&tmp)
	if err != nil {
		return err
	}
	//如果obj已经是指针，则此处不需要指针
	js := jsonconv.AnyToJsonSnake(tmp)
	err = json.Unmarshal([]byte(js), req)
	//Log.Logger(js)
	//Log.JsonIndent(req)
	if err != nil {
		return err
	}
	return nil
}

func ShouldBindQuery(c *gin.Context, req interface{}) error {
	// ShouldBindQuery使用tag "form"
	if err := c.ShouldBind(req); err != nil {
		return apierr.NewApiError(codex.CodeInvalidParam, err.Error())
	}
	isValid, ok := req.(IsValidChecker)
	if !ok {
		return nil
	}
	return isValid.IsValid()
}

func ShouldBind(c *gin.Context, req interface{}) error {
	if c.Request.Method == "GET" {
		return ShouldBindQuery(c, req)
	}
	return ShouldBindJSON(c, req)
}
