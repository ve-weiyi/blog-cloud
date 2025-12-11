package jsonconv

import (
	"encoding/json"
	"fmt"
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
