package resource

import (
	"runtime"
	"strings"
)

func GetRuntimeRoot() string {
	//获得当前方法运行文件名
	_, filename, _, _ := runtime.Caller(1)
	//找到 resource/language目录
	index := strings.LastIndex(filename, "/")
	root := filename[:index]
	return root
}

func GetTestdataRoot() string {
	root := GetRuntimeRoot()
	return root + "/testdata"
}

func GetTemplateRoot() string {
	root := GetRuntimeRoot()
	return root + "/template"
}
