package jsonconv

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// 默认json
func AnyToAny(data any, obj any) (err error) {
	jb, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jb, obj)
	if err != nil {
		return err
	}

	return nil
}

func AnyToAnyNE[T any](data any) (out *T) {
	jb, err := json.Marshal(data)
	if err != nil {
		fmt.Println("AnyToAny:json convert fail:", err)
		return nil
	}

	err = json.Unmarshal(jb, &out)
	if err != nil {
		fmt.Println("AnyToAny:json convert fail:", err)
		return nil
	}

	return out
}

func AnyToJson(data any) (string, error) {
	jb, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jb), nil

}

func AnyToJsonNE(data any) string {
	jb, err := json.Marshal(data)
	if err != nil {
		fmt.Println("AnyToJsonNE:json convert fail:", err)
		return ""
	}

	return string(jb)
}

// 转换行结构json
func AnyToJsonIndent(data any) string {
	jb, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		fmt.Println("AnyToJsonIndent:json convert fail:", err)
		return ""
	}
	if string(jb) == "{}" {
		return fmt.Sprintf("%+v", data)
	}
	return string(jb)
}

// 调用 JsonToAny(jsonStr , &obj)
func JsonToAny(js string, obj any) error {
	if js == "" {
		return nil
	}
	err := json.Unmarshal([]byte(js), obj)
	if err != nil {
		return err
	}

	return nil
}

// 不支持根据返回值类型推断
func JsonToAnyNE[T any](js string) (out T) {
	if js == "" {
		return out
	}
	err := json.Unmarshal([]byte(js), &out)
	if err != nil {
		fmt.Println("JsonToAny:json convert fail:", err)
		return out
	}

	return out
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
