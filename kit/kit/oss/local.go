package oss

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

type Local struct {
	dir string
}

func (s *Local) UploadFile(f io.Reader, prefix string, filename string) (filepath string, err error) {
	// 本地文件目录
	key := path.Join(prefix, filename)

	// 尝试创建上传目录
	dir := path.Dir(key)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("Local.UploadFile MkdirAll() Failed, err: %v", err)
	}

	// 创建目标文件
	out, err := os.Create(key)
	if err != nil {
		return "", fmt.Errorf("Local.UploadFile Create() Failed, err: %v", err)
	}
	defer out.Close()

	// 传输（拷贝）文件内容
	_, copyErr := io.Copy(out, f)
	if copyErr != nil {
		return "", fmt.Errorf("Local.UploadFile Copy() Failed, err: %v", copyErr)
	}

	return s.dir + "/" + key, nil
}

func (s *Local) DeleteFile(filepath string) error {
	p := s.dir + "/" + filepath

	if err := os.Remove(p); err != nil {
		return fmt.Errorf("Local.DeleteFile Remove() Failed, err: %v", err)
	}
	return nil
}

func (s *Local) ListFiles(prefix string, limit int) (files []*FileInfo, err error) {
	// 获取指定目录下的所有文件

	err = filepath.Walk(prefix, func(filepath string, info os.FileInfo, err error) error {
		if err != nil {
			return err // 如果遇到错误，返回
		}

		f := &FileInfo{
			IsDir:    info.IsDir(),
			FilePath: filepath,
			FileName: info.Name(),
			FileType: path.Ext(info.Name()),
			FileSize: info.Size(),
			FileUrl:  s.dir + "/" + filepath,
			UpTime:   info.ModTime().UnixMilli(),
		}
		files = append(files, f)

		// 如果已经达到限制数量，提前返回
		if len(files) >= limit {
			return fmt.Errorf("limit reached") // 触发 Walk 停止
		}

		return nil
	})

	if err != nil && err.Error() != "limit reached" {
		return nil, err // 如果是其他错误，返回
	}

	// 返回符合条件的文件列表
	return files, nil
}

func NewLocal(dir string) *Local {
	return &Local{
		dir: dir,
	}
}
