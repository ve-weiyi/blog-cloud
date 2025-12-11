package oss

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// https://developer.qiniu.com/kodo/1238/go
type Qiniu struct {
	cfg           *Config
	storageConfig *storage.Config
}

func (s *Qiniu) UploadFile(f io.Reader, prefix string, filename string) (filepath string, err error) {
	// 本地文件目录
	key := path.Join(prefix, filename)

	// 上传策略
	putPolicy := storage.PutPolicy{Scope: s.cfg.BucketName}
	mac := qbox.NewMac(s.cfg.AccessKeyId, s.cfg.AccessKeySecret)
	upToken := putPolicy.UploadToken(mac)
	resumeUploader := storage.NewResumeUploaderV2(s.storageConfig)
	// 上传文件
	ret := storage.PutRet{}
	putExtra := storage.RputV2Extra{}
	err = resumeUploader.PutWithoutSize(context.Background(), &ret, upToken, key, f, &putExtra)
	if err != nil {
		return "", fmt.Errorf("Qiniu.UploadHttpFile formUploader.Put() Filed, err: %v" + err.Error())
	}
	return s.cfg.BucketUrl + "/" + key, nil
}

func (s *Qiniu) DeleteFile(filepath string) error {
	mac := qbox.NewMac(s.cfg.AccessKeyId, s.cfg.AccessKeySecret)
	bucketManager := storage.NewBucketManager(mac, s.storageConfig)

	key := strings.TrimPrefix(filepath, s.cfg.BucketUrl+"/")

	if err := bucketManager.Delete(s.cfg.BucketName, key); err != nil {
		return fmt.Errorf("Qiniu.UploadHttpFile bucketManager.Delete() Filed, err: %v", err.Error())
	}
	return nil
}

func (s *Qiniu) ListFiles(prefix string, limit int) (files []*FileInfo, err error) {
	mac := qbox.NewMac(s.cfg.AccessKeyId, s.cfg.AccessKeySecret)
	bucketManager := storage.NewBucketManager(mac, s.storageConfig)

	// 参数设置
	delimiter := "/" // 模拟目录结构,目录分割标识
	marker := ""     // 分页标记

	entries, prefixes, _, _, err := bucketManager.ListFiles(s.cfg.BucketName, prefix, delimiter, marker, limit)
	if err != nil {
		return nil, fmt.Errorf("Qiniu.ListFiles bucketManager.ListFiles() Filed, err: %v" + err.Error())
	}

	for _, fix := range prefixes {
		f := &FileInfo{
			IsDir:    true,
			FilePath: fix,
			FileName: path.Base(fix),
			FileType: "",
			FileSize: 0,
			FileUrl:  s.cfg.BucketUrl + "/" + fix,
			UpTime:   0,
		}
		files = append(files, f)
	}

	for _, entry := range entries {
		if entry.Fsize == 0 {
			continue
		}

		f := &FileInfo{
			IsDir:    false,
			FilePath: entry.Key,
			FileName: path.Base(entry.Key),
			FileType: path.Ext(entry.Key),
			FileSize: entry.Fsize,
			FileUrl:  s.cfg.BucketUrl + "/" + entry.Key,
			UpTime:   entry.PutTime / 10000,
		}
		files = append(files, f)
	}

	return files, nil
}

func NewQiniu(conf *Config) *Qiniu {

	var region *storage.Region
	switch conf.Zone { // 根据配置文件进行初始化空间对应的机房
	case "ZoneHuadong":
		region = &storage.ZoneHuadong
	case "ZoneHuabei":
		region = &storage.ZoneHuabei
	case "ZoneHuanan":
		region = &storage.ZoneHuanan
	case "ZoneBeimei":
		region = &storage.ZoneBeimei
	case "ZoneXinjiapo":
		region = &storage.ZoneXinjiapo
	}
	//region, err := storage.GetRegion(conf.AccessKeyId, conf.BucketName)
	//if err != nil {
	//	return nil
	//}

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Region = region
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	return &Qiniu{
		cfg:           conf,
		storageConfig: &cfg,
	}
}
