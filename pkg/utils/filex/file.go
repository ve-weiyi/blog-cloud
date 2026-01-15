package filex

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// IsExist 检查文件或目录是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsDirExist 检查路径是否存在且为目录
func IsDirExist(path string) bool {
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil && fi.IsDir()
}

// IsFileExist 检查路径是否存在且为普通文件
func IsFileExist(path string) bool {
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil && fi.Mode().IsRegular()
}

// MkdirAll 创建目录及所有必要的父目录
func MkdirAll(path string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) && IsDirExist(path) {
		return nil
	}
	return os.MkdirAll(path, 0755)
}

// HasPermission 检查是否有权限访问路径
func HasPermission(path string) bool {
	_, err := os.Stat(path)
	return !os.IsPermission(err)
}

// CreateFile 创建文件及所有必要的父目录
func CreateFile(filename string) (*os.File, error) {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, fmt.Errorf("获取绝对路径失败: %w", err)
	}

	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录 %s 失败: %w", dir, err)
	}

	return os.OpenFile(absPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
}

// OpenOrCreate 打开文件用于追加，如果不存在则创建
func OpenOrCreate(filename string) (*os.File, error) {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, fmt.Errorf("获取绝对路径失败: %w", err)
	}

	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录 %s 失败: %w", dir, err)
	}

	return os.OpenFile(absPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
}

// Remove 删除文件或目录
func Remove(path string) error {
	return os.RemoveAll(path)
}

// Move 移动文件或目录从 src 到 dst
func Move(src, dst string) error {
	if dst == "" {
		return errors.New("目标路径不能为空")
	}

	srcAbs, err := filepath.Abs(src)
	if err != nil {
		return fmt.Errorf("获取源文件绝对路径失败: %w", err)
	}

	dstAbs, err := filepath.Abs(dst)
	if err != nil {
		return fmt.Errorf("获取目标文件绝对路径失败: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(dstAbs), 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	return os.Rename(srcAbs, dstAbs)
}

// Copy 复制文件从 src 到 dst
func Copy(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer srcFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// Walk 遍历以 root 为根的文件树，为每个文件或目录调用 fn
func Walk(root string, fn filepath.WalkFunc) error {
	return filepath.Walk(root, fn)
}

// WalkDir 遍历以 root 为根的文件树，为每个文件或目录调用 fn (Go 1.16+)
func WalkDir(root string, fn func(path string, d os.DirEntry, err error) error) error {
	return filepath.WalkDir(root, fn)
}

// WriteFile 写入数据到文件，如需要则创建父目录
func WriteFile(filename string, data []byte) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}
	return os.WriteFile(filename, data, 0644)
}

// WriteString 写入字符串到文件，如需要则创建父目录
func WriteString(filename, content string) error {
	return WriteFile(filename, []byte(content))
}

// ReadFile 读取整个文件并返回其内容
func ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// ReadString 读取整个文件并返回字符串内容
func ReadString(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// AppendFile 追加数据到文件，如果不存在则创建
func AppendFile(filename string, data []byte) error {
	f, err := OpenOrCreate(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

// AppendString 追加字符串到文件，如果不存在则创建
func AppendString(filename, content string) error {
	return AppendFile(filename, []byte(content))
}

// Size 返回文件大小（字节）
func Size(filename string) (int64, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}
