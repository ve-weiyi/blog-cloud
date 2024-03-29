package test

import (
	"go/ast"

	"github.com/dave/dst"
	jsoniter "github.com/json-iterator/go"
)

// 注册需要用到的rpc
type Context struct {
	visitor *ast.Ident   //
	jjson   jsoniter.API //注释
	Alias   *dst.Visitor ////元素别名
}

func NewApiContext(cfg string) *Context {

	_ = 1 //
	return &Context{
		visitor: ast.NewIdent("hello"),
	}
}
