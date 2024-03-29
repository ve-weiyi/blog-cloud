package initialize

import (
	"fmt"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"github.com/ve-weiyi/ve-blog-golang/server/global"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/upload"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/copyutil"
)

func Upload() {

	var cfg upload.UploadConfig
	var up upload.Uploader
	switch global.CONFIG.Upload.Mode {
	case "local":
		err := copyutil.DeepCopyByJson(global.CONFIG.Upload.Local, &cfg)
		if err != nil {
			global.LOG.Errorf("上传组件初始化失败！%v", err)
		}
		cfg.FileNameAsKey = fileNameAsKey
		up = upload.NewLocal(&cfg)
	case "qiniu":
		err := copyutil.DeepCopyByJson(global.CONFIG.Upload.Qiniu, &cfg)
		if err != nil {
			global.LOG.Errorf("上传组件初始化失败！%v", err)
		}
		cfg.FileNameAsKey = fileNameAsKey
		up = upload.NewQiniu(&cfg)
	case "aliyun":
		err := copyutil.DeepCopyByJson(global.CONFIG.Upload.Aliyun, &cfg)
		global.LOG.JsonIndent("cfg", cfg)
		if err != nil {
			global.LOG.Errorf("上传组件初始化失败！%v", err)
		}
		cfg.FileNameAsKey = fileNameAsKey
		up = upload.NewAliyunOSS(&cfg)
	default:
		up = upload.NewLocal(&cfg)
	}

	global.Uploader = up
	if global.Uploader != nil {
		global.LOG.Infof("上传组件初始化成功！%v", global.CONFIG.Upload.Mode)
	}
}

func fileNameAsKey(file *multipart.FileHeader) string {
	// 读取文件后缀
	ext := path.Ext(file.Filename)
	// 读取文件名并加密
	name := strings.TrimSuffix(file.Filename, ext)
	// 拼接新文件名
	filename := fmt.Sprintf("%s/%s%s", time.Now().Format("20060102150405"), name, ext)

	return filename
}
