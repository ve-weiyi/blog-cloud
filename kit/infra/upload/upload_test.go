package upload

import (
	"log"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpload(t *testing.T) {
	// 创建一个虚拟的 multipart.FileHeader
	fileHeader := &multipart.FileHeader{
		Filename: "test.txt",
		Size:     1234,
	}

	open, err := fileHeader.Open()
	log.Println("open:", open, err)
	if err != nil {
		return
	}

	local := NewLocal(&UploadConfig{
		BucketUrl: "http://localhost:9999",
		BasePath:  "runtime/uploads",
	})
	//// 调用 UploadFile 方法进行测试
	url, err := local.UploadFile("test", fileHeader)

	// 验证结果是否符合预期
	assert.NoError(t, err)  // 验证错误是否为 nil
	assert.NotEmpty(t, url) // 验证返回的 url 不为空
}
