package request

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/biz/bizerr"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/restx"
)

// 请求上下文,一般存放请求头参数
type Context struct {
	context.Context
	Token     string `json:"token" header:"token" example:""`
	Uid       int64  `json:"uid" header:"-" example:""`
	IpAddress string `json:"ip_address" header:"-" example:""`
	UserAgent string `json:"user_agent" header:"-" example:""`
}

// 获取请求上下文
func ParseRequestContext(c *gin.Context) (*Context, error) {
	reqCtx := &Context{}
	reqCtx.Token = c.GetHeader(restx.HeaderToken)
	reqCtx.Uid = cast.ToInt64(c.GetHeader(restx.HeaderUid))
	reqCtx.IpAddress = c.ClientIP()
	reqCtx.UserAgent = c.Request.UserAgent()
	reqCtx.Context = c.Request.Context()
	return reqCtx, nil
}

type IsValidChecker interface {
	IsValid() error
}

func ShouldBindJSON(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindJSON(&req); err != nil {
		return bizerr.NewBizError(bizerr.CodeInvalidParam, err.Error())
	}

	isValid, ok := req.(IsValidChecker)
	if !ok {
		return nil
	}

	if err := isValid.IsValid(); err != nil {
		return bizerr.NewBizError(bizerr.CodeInvalidParam, err.Error())
	}

	return nil
}

func ShouldBindQuery(c *gin.Context, req interface{}) error {
	// ShouldBindQuery使用tag "form"
	if err := c.ShouldBind(req); err != nil {
		return bizerr.NewBizError(bizerr.CodeInvalidParam, err.Error())
	}
	isValid, ok := req.(IsValidChecker)
	if !ok {
		return nil
	}

	if err := isValid.IsValid(); err != nil {
		return bizerr.NewBizError(bizerr.CodeInvalidParam, err.Error())
	}

	return nil
}

func ShouldBind(c *gin.Context, req interface{}) error {
	if c.Request.Method == "GET" {
		return ShouldBindQuery(c, req)
	}
	return ShouldBindJSON(c, req)
}
