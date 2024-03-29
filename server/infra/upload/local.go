package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

type Local struct {
	cfg *UploadConfig
}

func (s *Local) UploadFile(prefix string, file *multipart.FileHeader) (url string, err error) {
	var filename string
	// 读取文件名
	if s.cfg.FileNameAsKey != nil {
		filename = s.cfg.FileNameAsKey(file)
	} else {
		filename = file.Filename
	}

	// 本地文件目录
	dir := path.Join(s.cfg.BasePath, prefix)
	// 本地文件路径
	localPath := path.Join(dir, filename)
	// 访问文件路径
	accessPath := s.cfg.BucketUrl + "/" + localPath

	// 尝试创建上传目录
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("Local.UploadFile os.MkdirAll() Filed, err:" + err.Error())
	}
	// 读取文件
	f, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("Local.UploadFile file.Open() Filed, err:" + err.Error())
	}
	defer f.Close()
	// 创建文件
	out, err := os.Create(localPath)
	if err != nil {
		return "", fmt.Errorf("Local.UploadFile os.Create() Filed, err:" + err.Error())
	}
	defer out.Close()
	// 传输（拷贝）文件内容
	_, copyErr := io.Copy(out, f)
	if copyErr != nil {
		return "", fmt.Errorf("Local.UploadFile io.Copy() Filed, err:" + copyErr.Error())
	}

	return accessPath, nil
}

func (s *Local) DeleteFile(key string) error {
	p := s.cfg.BucketUrl + "/" + key
	if strings.Contains(p, s.cfg.BasePath) {
		if err := os.Remove(p); err != nil {
			return fmt.Errorf("本地文件删除失败, err:" + err.Error())
		}
	}
	return nil
}

func NewLocal(cfg *UploadConfig) *Local {
	return &Local{
		cfg: cfg,
	}
}
