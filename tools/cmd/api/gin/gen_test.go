package gin

import (
	"fmt"
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/tools/parserx/apiparser/aspec"
)

func TestName(t *testing.T) {
	tp := aspec.DefineStruct{
		RawName: "UserInfoResp",
		Members: []aspec.Member{
			{
				Name:    "Id",
				Type:    aspec.PrimitiveType{RawName: "int64"},
				Tag:     "`json:\"id\"`",
				Comment: "用户ID",
			},
			{
				Name:    "Username",
				Type:    aspec.PrimitiveType{RawName: "string"},
				Tag:     "`json:\"username\"`",
				Comment: "用户名",
			},
			{
				Name:    "Email",
				Type:    aspec.PrimitiveType{RawName: "string"},
				Tag:     "`json:\"email\"`",
				Comment: "邮箱",
			},
		},
		Docs: nil,
	}
	val, err := buildTypes(tp)
	if err != nil {
		return
	}

	fmt.Println("---", val)
}
