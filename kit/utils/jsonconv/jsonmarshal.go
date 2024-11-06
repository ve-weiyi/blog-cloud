package jsonconv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
)

// https://www.cnblogs.com/chenqionghe/p/13067596.html
/*************************************** 下划线json ***************************************/
type JsonSnakeCase struct {
	Value interface{}
}

// 转换为json，key全部变为驼峰
func (c JsonSnakeCase) MarshalJSON() ([]byte, error) {
	// Regexp definitions
	var keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)
	var wordBarrierRegex = regexp.MustCompile(`(\w)([A-Z])`)
	marshalled, err := json.Marshal(c.Value)
	converted := keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {
			return bytes.ToLower(wordBarrierRegex.ReplaceAll(
				match,
				[]byte(`${1}_${2}`),
			))
		},
	)
	return converted, err
}

/*************************************** 驼峰json ***************************************/
type JsonCamelCase struct {
	Value interface{}
}

// 转换为json，key全部变为下划线
func (c JsonCamelCase) MarshalJSON() ([]byte, error) {
	var keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)
	marshalled, err := json.Marshal(c.Value)
	converted := keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {
			matchStr := string(match)
			key := matchStr[1 : len(matchStr)-2]
			if key == "id" {
				return match // 不转换"id"字段
			}
			resKey := FirstLower(Case2Camel(key))
			return []byte(`"` + resKey + `":`)
		},
	)
	return converted, err
}

/*************************************** 其他方法 ***************************************/
// 转下划线json
func AnyToJsonSnake(data any) string {
	jb, err := json.Marshal(JsonSnakeCase{Value: data})
	if err != nil {
		fmt.Println("AnyToJsonSnake:json convert fail:", err)
		return ""
	}

	return string(jb)
}

// 转首字母驼峰json
func AnyToJsonCamel(data any) string {
	jb, err := json.Marshal(JsonCamelCase{Value: data})
	if err != nil {
		fmt.Println("AnyToJsonCamel:json convert fail:", err)
		return ""
	}

	return string(jb)
}
