package jsonconv

import (
	"fmt"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

/*
*
ConfigDefault(默认API行为)、
ConfigCompatibleWithStandardLibrary(支持标准库的行为，比如encoding/jjson)、
ConfigFast(通过忽略float类型数据的精度保证最高效的性能)
*/
var jjson jsoniter.API

func init() {
	//- `IndentionStep`：设置 JSON 缩进的空格数，默认为 0，表示不缩进。
	//- `MarshalFloatWith6Digits`：设置是否将浮点数序列化为 6 位小数，默认为 false。
	//- `EscapeHTML`：设置是否将 HTML 字符转义，默认为 true。
	//- `SortMapKeys`：设置是否按照键名排序序列化 map，默认为 false。
	//- `UseNumber`：设置是否将数字解码为 `json.Number` 类型，默认为 false。
	//- `DisallowUnknownFields`：设置是否在解码时禁止未知字段，默认为 false。
	//- `TagKey`：设置结构体 tag 的键名，默认为 `json`。
	//- `OnlyTaggedField`：设置是否只序列化带有 tag 的结构体字段，默认为 false。
	//- `ValidateJsonRawMessage`：设置是否验证 `json.RawMessage` 类型的字段，默认为 false。
	//- `ObjectFieldMustBeSimpleString`：设置是否要求对象类型的字段必须是简单字符串，默认为 false。
	//- `CaseSensitive`：设置是否大小写敏感，默认为 true。
	jjson = jsoniter.Config{
		CaseSensitive: false,
	}.Froze()
}

// 序列化成字符串
func TestMarshalToString(t *testing.T) {
	order := struct {
		Id       int
		OrderNum string
		Money    float32
		PayTime  time.Time
		Extend   map[string]string
	}{
		Id:       10,
		OrderNum: "100200300",
		Money:    99.99,
		PayTime:  time.Now(),
		Extend:   map[string]string{"name": "张三"},
	}
	// 定义
	var jsonNew = jsoniter.ConfigCompatibleWithStandardLibrary
	// 设置后，没有json标签的属性，会自动转成 xx_xx
	extra.SetNamingStrategy(extra.LowerCaseWithUnderscores)
	// 直接转成字符串
	jsonStr, _ := jsonNew.MarshalToString(order)
	fmt.Println("jsonStr:", jsonStr)

	jb, _ := jjson.MarshalIndent(order, "", " ")
	fmt.Println("jsonStr:", string(jb))
}
