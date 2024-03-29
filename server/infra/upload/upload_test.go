package upload

import (
	"log"
	"mime/multipart"
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/server/infra/initest"
)

func TestUpload(t *testing.T) {
	initest.Init()

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

	//local := NewLocal(&global.CONFIG.Upload.Local)
	//// 调用 UploadFile 方法进行测试
	//url, path, err := local.UploadFile(fileHeader)

	// 验证结果是否符合预期
	//assert.NoError(t, err)   // 验证错误是否为 nil
	//assert.NotEmpty(t, url)  // 验证返回的 URL 不为空
	//assert.NotEmpty(t, path) // 验证返回的路径不为空
	// 可根据实际情况进行更多的验证

	// 进一步验证返回的 URL 和路径是否符合预期要求
	// assert.Equal(t, expectedURL, url)
	// assert.Equal(t, expectedPath, path)
}
