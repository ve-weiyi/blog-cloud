package oss

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

type Local struct {
	cfg *Config
}

func (s *Local) UploadFile(f io.Reader, prefix string, filename string) (filepath string, err error) {
	// 本地文件目录
	key := path.Join(prefix, filename)

	// 尝试创建上传目录
	err = os.MkdirAll(prefix, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("Local.UploadHttpFile os.MkdirAll() Filed, err: %v" + err.Error())
	}

	// 创建目标文件
	out, err := os.Create(key)
	if err != nil {
		return "", fmt.Errorf("Local.UploadHttpFile os.Create() Filed, err: %v" + err.Error())
	}
	defer out.Close()

	// 传输（拷贝）文件内容
	_, copyErr := io.Copy(out, f)
	if copyErr != nil {
		return "", fmt.Errorf("Local.UploadHttpFile io.Copy() Filed, err:" + copyErr.Error())
	}

	return s.cfg.BucketUrl + "/" + key, nil
}

func (s *Local) DeleteFile(filepath string) error {
	p := s.cfg.BucketUrl + "/" + filepath

	if err := os.Remove(p); err != nil {
		return fmt.Errorf("本地文件删除失败, err: %v" + err.Error())
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
			FileUrl:  s.cfg.BucketUrl + "/" + filepath,
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

func NewLocal(cfg *Config) *Local {
	return &Local{
		cfg: cfg,
	}
}
