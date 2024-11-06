package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin/render"
	"github.com/go-sql-driver/mysql"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
)

// mysql 错误码
const (
	SqlKeyDuplicateCode = 1062
	DataTooLongError    = 1406
)

// https://fromdual.com/mysql-error-codes-and-messages-1050-1099
func SqlErrorI18n(errorMessage *mysql.MySQLError) string {
	var message string
	switch int(errorMessage.Number) {
	case SqlKeyDuplicateCode:
		// Message: Duplicate entry '%s' for key %d
		reg := regexp.MustCompile(`Duplicate entry ('(.*?)')`)
		subText := reg.FindString(errorMessage.Message)
		start := len("Duplicate entry '")
		msg := []rune(subText)[start : utf8.RuneCountInString(subText)-1]
		message = fmt.Sprintf("%s已存在，请更改后提交", string(msg))
	case DataTooLongError:
		// Data too long for column 'value' at row 1
		reg := regexp.MustCompile(`Data too long for column ('(.*?)')`)
		subText := reg.FindString(errorMessage.Error())
		start := len("Data too long for column '")
		column := []rune(subText)[start : utf8.RuneCountInString(subText)-1]
		msg := "tableColumns." + string(column)
		message = fmt.Sprintf("%s 内容超过数据库字段限制，请重新输入", msg)
	default:
		message = errorMessage.Error()
	}
	return message
}

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
