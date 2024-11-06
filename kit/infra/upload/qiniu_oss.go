package upload

import (
	"context"
	"fmt"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// https://developer.qiniu.com/kodo/1238/go
type Qiniu struct {
	cfg           *UploadConfig
	storageConfig *storage.Config
}

func (s *Qiniu) UploadFile(prefix string, file *multipart.FileHeader) (url string, err error) {
	var filename string
	// 读取文件名
	filename = s.FileNameAsKey(file)

	// 本地文件目录
	dir := path.Join(s.cfg.BasePath, prefix)
	// 本地文件路径
	localPath := path.Join(dir, filename)
	// 访问文件路径
	accessPath := s.cfg.BucketUrl + "/" + localPath

	// 读取本地文件
	f, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("Qiniu.UploadFile file.Open() Failed, err:" + err.Error())
	}
	defer f.Close() // 创建文件 defer 关闭

	// 上传策略
	putPolicy := storage.PutPolicy{Scope: s.cfg.BucketName}
	mac := qbox.NewMac(s.cfg.AccessKeyId, s.cfg.AccessKeySecret)
	upToken := putPolicy.UploadToken(mac)
	resumeUploader := storage.NewResumeUploaderV2(s.storageConfig)
	// 上传文件
	ret := storage.PutRet{}
	putExtra := storage.RputV2Extra{}
	putErr := resumeUploader.Put(context.Background(), &ret, upToken, localPath, f, file.Size, &putExtra)
	if putErr != nil {
		return "", fmt.Errorf("Qiniu.UploadFile formUploader.Put() Filed, err:" + putErr.Error())
	}
	return accessPath, nil
}

func (s *Qiniu) DeleteFile(key string) error {
	mac := qbox.NewMac(s.cfg.AccessKeyId, s.cfg.AccessKeySecret)
	bucketManager := storage.NewBucketManager(mac, s.storageConfig)
	if err := bucketManager.Delete(s.cfg.BucketName, key); err != nil {
		return fmt.Errorf("Qiniu.UploadFile bucketManager.Delete() Filed, err:" + err.Error())
	}
	return nil
}

func (s *Qiniu) FileNameAsKey(file *multipart.FileHeader) string {
	// 读取文件后缀
	ext := path.Ext(file.Filename)
	// 读取文件名并加密
	name := strings.TrimSuffix(file.Filename, ext)
	// 拼接新文件名
	filename := fmt.Sprintf("%s-%s%s", name, time.Now().Format("20060102150405"), ext)

	return filename
}
func NewQiniu(conf *UploadConfig) *Qiniu {

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
