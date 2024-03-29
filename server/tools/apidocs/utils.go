package apidocs

import (
	"strings"
)

func convertGoTypeToTsType(name string) string {
	if strings.HasPrefix(name, "*") {
		return convertGoTypeToTsType(name[1:]) // 指针
	}
	if strings.HasPrefix(name, "[]") {
		return convertGoTypeToTsType(name[2:]) + "[]" // 数组
	}
	if strings.LastIndex(name, ".") > 0 {
		return convertGoTypeToTsType(name[strings.LastIndex(name, ".")+1:]) // 去掉包名
	}
	switch name {
	case "int", "int32", "int64", "uint", "uint32", "uint64", "float32", "float64":
		return "number"
	case "string":
		return "string"
	case "bool":
		return "boolean"
	case "file":
		return "File"
	case "Time":
		return "string"
	case "FileHeader":
		return "File"
	case "interface{}", "object":
		return "any"
	default:
		return name
	}
}

// [a,b,c] -> a, b, c
func joinArray(arr []string) string {
	var result string
	for i, v := range arr {
		result += v
		if i < len(arr)-1 {
			result += ", "
		}
	}
	return result
}

// 数组 [a, b, c] 转换为字符串 a<b<c>>
func joinBracket(params []string) string {
	var result string
	for i, val := range params {
		if i > 0 {
			result += "<"
		}
		result += val
	}
	for i := 0; i < len(params)-1; i++ {
		result += ">"
	}

	return result
}

// 删除数组中的重复元素
func removeDuplicates(input []string) []string {
	encountered := map[string]bool{}
	var result []string

	for _, v := range input {
		if encountered[v] == false {
			encountered[v] = true
			result = append(result, v)
		}
	}

	return result
}
