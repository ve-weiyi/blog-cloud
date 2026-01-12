package oss

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Local(t *testing.T) {
	uploader := NewLocal(&Config{
		BucketUrl: "http://localhost:9999",
	})

	files, err := uploader.ListFiles("./", 10)
	assert.NoError(t, err) // 验证错误是否为 nil
	for _, file := range files {
		t.Log(file)
	}
}

func Test_Qiniu(t *testing.T) {
	uploader := NewQiniu(&Config{
		Zone:            "ZoneHuanan",
		Endpoint:        "s3.cn-south-1.qiniucs.com",
		AccessKeyId:     "",
		AccessKeySecret: "",
		BucketName:      "veweiyi",
		BucketUrl:       "https://static.veweiyi.cn",
	})

	_, err := uploader.UploadFile(bytes.NewReader([]byte("")), "blog/test/", "")
	assert.NoError(t, err) // 验证错误是否为 nil

	files, err := uploader.ListFiles("blog/", 100)
	assert.NoError(t, err) // 验证错误是否为 nil
	for _, file := range files {
		t.Log(file)
	}

	//err = uploader.DeleteFile("blog/test/")
	//assert.NoError(t, err) // 验证错误是否为 nil
}

func Test_Aliyun(t *testing.T) {
	uploader := NewAliyunOSS(&Config{
		Zone:            "",
		Endpoint:        "oss-cn-shenzhen.aliyuncs.com",
		AccessKeyId:     "",
		AccessKeySecret: "",
		BucketName:      "ve-blog",
		BucketUrl:       "http://ve-blog.oss-cn-shenzhen.aliyuncs.com",
	})

	files, err := uploader.ListFiles("", 10)
	assert.NoError(t, err) // 验证错误是否为 nil
	for _, file := range files {
		t.Log(file)
	}
}
