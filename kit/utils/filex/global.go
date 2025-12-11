package filex

import (
	"os"
	"path"
	"runtime"
	"strings"
)

func GetRuntimeRoot() string {
	// 获得当前方法运行文件名
	_, filename, _, _ := runtime.Caller(1)
	// 找到 resource/language目录
	index := strings.LastIndex(filename, "/")
	root := filename[:index]
	return root
}

func ToAbs(filename string) string {
	if path.IsAbs(filename) {
		return filename
	}

	dir, err := os.Getwd()
	if err != nil {
		return filename
	}

	return path.Join(dir, filename)
}
