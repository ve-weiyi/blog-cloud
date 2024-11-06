package slicex

import (
	"reflect"
	"testing"
)

// 测试Join函数
// 原函数：Join ["a", "b", "c"], "," => "a,b,c"
func TestJoin(t *testing.T) {
	tests := []struct {
		slice []string
		sep   string
		want  string
	}{
		{[]string{"a", "b", "c"}, ",", "a,b,c"},
		{[]string{"1", "2", "3"}, "-", "1-2-3"},
		{[]string{}, ",", ""},
		{[]string{"single"}, ",", "single"},
	}

	for _, tt := range tests {
		t.Run("Join", func(t *testing.T) {
			if got := Join(tt.slice, tt.sep); got != tt.want {
				t.Errorf("Join(%v, %q) = %q, want %q", tt.slice, tt.sep, got, tt.want)
			}
		})
	}
}

// 测试Reverse函数
// 原函数：Reverse reverses a slice.
func TestReverse(t *testing.T) {
	tests := []struct {
		slice []int
		want  []int
	}{
		{[]int{1, 2, 3}, []int{3, 2, 1}},
		{[]int{}, []int{}},
		{[]int{1}, []int{1}},
	}

	for _, tt := range tests {
		t.Run("Reverse", func(t *testing.T) {
			Reverse(tt.slice)
			if !reflect.DeepEqual(tt.slice, tt.want) {
				t.Errorf("Reverse() got %v, want %v", tt.slice, tt.want)
			}
		})
	}
}

// 测试Exist函数
// 原函数：Exist 判断数组中是否存在指定元素
func TestExist(t *testing.T) {
	tests := []struct {
		slice []int
		elem  int
		want  bool
	}{
		{[]int{1, 2, 3}, 2, true},
		{[]int{1, 2, 3}, 4, false},
		{[]int{}, 1, false},
		{[]int{1}, 1, true},
	}

	for _, tt := range tests {
		t.Run("Exist", func(t *testing.T) {
			if got := Exist(tt.slice, tt.elem); got != tt.want {
				t.Errorf("Exist(%v, %d) = %v, want %v", tt.slice, tt.elem, got, tt.want)
			}
		})
	}
}

// 测试Append函数
// 原函数：Append 添加元素到数组中，返回一个新的数组，不改变原数组
func TestAppend(t *testing.T) {
	tests := []struct {
		slice []int
		elem  int
		want  []int
	}{
		{[]int{1, 2, 3}, 4, []int{1, 2, 3, 4}},
		{[]int{}, 1, []int{1}},
		{[]int{1}, 2, []int{1, 2}},
	}

	for _, tt := range tests {
		t.Run("Append", func(t *testing.T) {
			if got := Append(tt.slice, tt.elem); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Append(%v, %d) = %v, want %v", tt.slice, tt.elem, got, tt.want)
			}
		})
	}
}

// 测试Remove函数
// 原函数：Remove 删除数组中指定元素，返回一个新的数组，不改变原数组
func TestRemove(t *testing.T) {
	tests := []struct {
		slice []int
		elem  int
		want  []int
	}{
		{[]int{1, 2, 3, 2}, 2, []int{1, 3}},
		{[]int{1, 2, 3}, 4, []int{1, 2, 3}},
		{[]int{}, 1, []int{}},
		{[]int{1}, 1, []int{}},
	}

	for _, tt := range tests {
		t.Run("Remove", func(t *testing.T) {
			if got := Remove(tt.slice, tt.elem); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Remove(%v, %d) = %v, want %v", tt.slice, tt.elem, got, tt.want)
			}
		})
	}
}

// 测试RemoveDuplicates函数
// 原函数：RemoveDuplicates 删除数组中的重复元素，返回一个新的数组，不改变原数组
func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		slice []int
		want  []int
	}{
		{[]int{1, 2, 2, 3}, []int{1, 2, 3}},
		{[]int{1, 1, 1}, []int{1}},
		{[]int{}, []int{}},
		{[]int{1}, []int{1}},
	}

	for _, tt := range tests {
		t.Run("RemoveDuplicates", func(t *testing.T) {
			if got := RemoveDuplicates(tt.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveDuplicates(%v) = %v, want %v", tt.slice, got, tt.want)
			}
		})
	}
}

// 测试MapKeys函数
// 原函数：MapKeys 返回映射 m 中的所有键，按任意顺序返回一个键的切片
func TestMapKeys(t *testing.T) {
	tests := []struct {
		m    map[int]string
		want []int
	}{
		{map[int]string{1: "a", 2: "b"}, []int{1, 2}},
		{map[int]string{}, []int{}},
		{map[int]string{1: "a"}, []int{1}},
	}

	for _, tt := range tests {
		t.Run("MapKeys", func(t *testing.T) {
			got := MapKeys(tt.m)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapKeys(%v) = %v, want %v", tt.m, got, tt.want)
			}
		})
	}
}
