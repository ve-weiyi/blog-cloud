package slicex

import (
	"fmt"
)

// Join ["a", "b", "c"], "," => "a,b,c"
func Join[S ~[]E, E any](s S, sep string) string {
	var out string

	for _, e := range s {
		if out == "" {
			out = out + fmt.Sprintf("%v", e)
		} else {
			out = out + sep
			out = out + fmt.Sprintf("%v", e)
		}
	}

	return out
}

// Reverse reverses a slice.
func Reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Exist 判断数组中是否存在指定元素
func Exist[S ~[]E, E comparable](s S, e E) bool {
	for _, item := range s {
		if item == e {
			return true
		}
	}
	return false
}

// Append 添加元素到数组中，返回一个新的数组，不改变原数组
func Append[S ~[]E, E any](s S, e E) S {
	return append(s, e)
}

// Remove 删除数组中指定元素，返回一个新的数组，不改变原数组
func Remove[S ~[]E, E comparable](s S, e E) (ret S) {
	ret = make(S, 0)
	for _, se := range s {
		if se != e {
			ret = append(ret, se)
		}
	}

	return
}

// RemoveDuplicates 删除数组中的重复元素，返回一个新的数组，不改变原数组
func RemoveDuplicates[S ~[]E, E comparable](s S) (ret S) {
	seen := make(map[E]struct{}) // 用于记录已出现的元素

	ret = make(S, 0)
	for _, e := range s {
		if _, exists := seen[e]; !exists {
			seen[e] = struct{}{} // 标记元素已出现
			ret = append(ret, e) // 仅添加未重复的元素
		}
	}
	return ret
}

// MapKeys 返回映射 m 中的所有键，按任意顺序返回一个键的切片
func MapKeys[M ~map[K]V, K comparable, V any](m M) []K {
	keys := make([]K, 0, len(m)) // 创建一个切片，容量为映射的键数量
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
