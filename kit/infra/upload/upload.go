package upload

import (
	"mime/multipart"
)

type UploadConfig struct {
	Zone            string `json:"zone"`
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"access-key-id"`
	AccessKeySecret string `json:"access-key-secret"`
	BucketName      string `json:"bucket-name"`
	BucketUrl       string `json:"bucket-url"`
	BasePath        string `json:"base-path"`
}

// Uploader 对象存储接口
type Uploader interface {
	UploadFile(prefix string, file *multipart.FileHeader) (string, error)
	DeleteFile(key string) error
	FileNameAsKey(file *multipart.FileHeader) string
}
