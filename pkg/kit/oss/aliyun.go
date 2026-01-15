package oss

import (
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Aliyun struct {
	cfg    *Config
	bucket *oss.Bucket
}

func (s *Aliyun) UploadFile(f io.Reader, prefix string, filename string) (filepath string, err error) {
	// 本地文件目录
	key := path.Join(prefix, filename)

	// 上传策略
	bucket := s.bucket
	// 上传文件流。
	err = bucket.PutObject(key, f)
	if err != nil {
		return "", fmt.Errorf("Aliyun.UploadFile PutObject() Failed, err: %v", err)
	}

	return s.cfg.BucketUrl + "/" + key, nil
}

func (s *Aliyun) DeleteFile(filepath string) (err error) {
	bucket := s.bucket

	key := strings.TrimPrefix(filepath, s.cfg.BucketUrl+"/")
	// 删除单个文件。objectName表示删除OSS文件时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	// 如需删除文件夹，请将objectName设置为对应的文件夹名称。如果文件夹非空，则需要将文件夹下的所有object删除后才能删除该文件夹。
	err = bucket.DeleteObject(key)
	if err != nil {
		return fmt.Errorf("Aliyun.DeleteFile DeleteObject() Failed, err: %v", err)
	}

	return nil
}

func (s *Aliyun) ListFiles(prefix string, limit int) (files []*FileInfo, err error) {
	bucket := s.bucket

	result, err := bucket.ListObjectsV2(oss.Prefix(prefix), oss.MaxKeys(limit))
	if err != nil {
		return nil, fmt.Errorf("Aliyun.ListFiles ListObjectsV2() Failed, err: %v", err)
	}

	for _, entry := range result.Objects {
		f := &FileInfo{
			IsDir:    entry.Type == "",
			FilePath: entry.Key,
			FileName: path.Base(entry.Key),
			FileType: path.Ext(entry.Key),
			FileSize: entry.Size,
			FileUrl:  s.cfg.BucketUrl + "/" + entry.Key,
			UpTime:   entry.LastModified.UnixMilli(),
		}
		files = append(files, f)
	}

	return files, nil
}

func NewAliyunOSS(cfg *Config) *Aliyun {

	// 创建OSSClient实例。
	client, err := oss.New(cfg.Endpoint, cfg.AccessKeyId, cfg.AccessKeySecret)
	if err != nil {
		panic(err)
	}

	// 获取存储空间。
	bucket, err := client.Bucket(cfg.BucketName)
	if err != nil {
		panic(err)
	}

	return &Aliyun{
		cfg:    cfg,
		bucket: bucket,
	}
}
