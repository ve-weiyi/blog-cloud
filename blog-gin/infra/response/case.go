package response

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin/render"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
)

// json解析忽略大小写
type caseInsensitiveDecoder struct {
	*json.Decoder
}

func (d *caseInsensitiveDecoder) Token() (json.Token, error) {
	t, err := d.Decoder.Token()
	if err != nil {
		return nil, err
	}
	if s, ok := t.(string); ok {
		return strings.ToLower(s), nil
	}
	return t, nil
}

type underscoreJSONRender struct {
	render.JSON
}

// 转换下划线的返回值
func (r underscoreJSONRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	// 获取原始的 JSON 数据
	data := r.JSON.Data

	// 转换 JSON 数据中的字段名为下划线格式
	convertedJSON := jsonconv.AnyToJsonSnake(data)
	// 将转换后的 JSON 数据写入响应
	_, err := w.Write([]byte(convertedJSON))
	return err
}

type camelJSONRender struct {
	render.JSON
}

// 转换下划线的返回值
func (r camelJSONRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	// 获取原始的 JSON 数据
	data := r.JSON.Data

	// 转换 JSON 数据中的字段名为首字母小写的驼峰格式
	convertedJSON := jsonconv.AnyToJsonCamel(data)
	// 将转换后的 JSON 数据写入响应
	_, err := w.Write([]byte(convertedJSON))
	return err
}
