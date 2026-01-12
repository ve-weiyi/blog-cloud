package jsonconv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
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

// json解析忽略大小写
func JsonToAnyIgnoreCase(data string, obj interface{}) error {
	// 使用反射将匿名结构体中的字段值赋值给obj对象
	v := reflect.ValueOf(obj)
	t := v.Elem().Type()
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		panic("SetCamelCaseJsonTag only accepts a pointer to a struct")
	}

	toLowLetters := func(key string) string {
		return strings.ToLower(ExtractLetters(key))
	}

	tmp := make(map[string]interface{})
	// 解析JSON数据到匿名结构体中
	if err := json.Unmarshal([]byte(data), &tmp); err != nil {
		return err
	}

	lmp := make(map[string]interface{})
	for key, value := range tmp {
		lmp[toLowLetters(key)] = value
	}

	for i := 0; i < v.Elem().NumField(); i++ {
		//目标字段名
		fieldName := t.Field(i).Name
		fieldType := v.Elem().Field(i)

		value, ok := lmp[toLowLetters(fieldName)]
		if ok {
			setFieldValue(fieldType, value)
			continue
		}

	}
	return nil
}

func setFieldValue(fieldType reflect.Value, value interface{}) {
	jsonValue := reflect.ValueOf(value)
	//log.Println("--", fieldType.Kind(), jsonValue.Kind())

	switch fieldType.Kind() {
	//case reflect.String:
	//	fieldType.SetString(jsonValue.String())
	//case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	//	fieldType.SetInt(int64(jsonValue.Float()))
	//case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	//	fieldType.SetUint(jsonValue.Uint())
	//case reflect.Float32, reflect.Float64:
	//	fieldType.SetFloat(jsonValue.Float())
	//case reflect.Bool:
	//	fieldType.SetBool(jsonValue.Bool())
	case reflect.Slice:
		// 创建一个新的切片
		slice := reflect.MakeSlice(fieldType.Type(), 0, 0)

		// 迭代JSON数组中的每个元素
		for i := 0; i < jsonValue.Len(); i++ {
			// 获取JSON数组中的每个元素
			elementValue := jsonValue.Index(i).Interface()

			// 创建一个新的切片元素并设置其值
			newElement := reflect.New(fieldType.Type().Elem()).Elem()
			setFieldValue(newElement, elementValue)

			// 将新的切片元素添加到切片中
			slice = reflect.Append(slice, newElement)
		}

		// 将切片设置为字段的值
		fieldType.Set(slice)

	default:
		fieldType.Set(jsonValue.Convert(fieldType.Type()))
	}
}
