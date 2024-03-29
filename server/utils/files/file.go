package files

import (
	"fmt"
	"os"
	"path/filepath"
)

// 检查文件夹或者文件是否存在
func IsExist(src string) bool {
	_, err := os.Stat(src)
	return !os.IsNotExist(err)
}

// 检查文件夹是否存在
func IsDirExist(path string) bool {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return true
		}
	}

	return !os.IsNotExist(err)
}

// 检查文件是否存在
func IsFileExist(path string) bool {
	fi, err := os.Stat(path)
	if err == nil {
		if !fi.IsDir() {
			return true
		}
	}

	return !os.IsNotExist(err)
}

// IsNotExistMkDir 如果不存在则新建文件夹
// `os.Mkdir` 用于创建一个目录。
// `os.MkdirAll` 用于创建一个目录以及它的所有父目录（如果它们不存在的话）
func MkDir(src string) error {
	// 已存在则返回
	if IsDirExist(src) {
		return nil
	}

	// 不存在则创建
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create directory '%s' error: %v", src, err)
	}

	return nil
}

// CheckPermission 检查文件权限
func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

// 打开文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// 打开文件 不存在则创建
func OpenExistFile(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("获取目录: %v", err)
	}

	src := dir + "/" + filePath
	perm := CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("无此权限: %s", src)
	}

	err = MkDir(src)
	if err != nil {
		return nil, fmt.Errorf("创建文件夹错误: %s, err: %v", src, err)
	}

	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败 :%v", err)
	}

	return f, nil
}

// 删除文件
func DeLFile(filePath string) error {
	return os.RemoveAll(filePath)
}

// @description: 文件移动供外部调用
// @param: src string, dst string(src: 源位置,绝对路径or相对路径, dst: 目标位置,绝对路径or相对路径,必须为文件夹)
func MoveFile(src string, dst string) (err error) {
	if dst == "" {
		return nil
	}
	src, err = filepath.Abs(src)
	if err != nil {
		return err
	}
	dst, err = filepath.Abs(dst)
	if err != nil {
		return err
	}
	dir := filepath.Dir(dst)
	err = MkDir(dir)
	if err != nil {
		return err
	}
	return os.Rename(src, dst)
}

// 深度遍历目录下的所有文件，包括目录和文件
func VisitFile(root string, visitFile func(path string, f os.FileInfo, err error) error) error {
	return filepath.Walk(root, visitFile)
}

// 向文件中写入内容
func WriteFile(filename string, content string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	l, err := f.WriteString(content)
	if err != nil {
		return err
	}

	fmt.Println(l, "bytes written successfully")
	return nil
}
