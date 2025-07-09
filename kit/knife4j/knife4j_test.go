package knife4j

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func Test_Knife4j(t *testing.T) {
	f := []byte(` 
{ 
  "swagger": "2.0",
  "info": {
    "title": "",
    "version": ""
  },
  "paths": {
    "/api/v1/ping": {
      "get": {
        "summary": "ping",
        "operationId": "Ping"
      }
    }
  }
}
`)

	prefix := "/api/v1/swagger/"

	http.HandleFunc(fmt.Sprintf("%s%s", prefix, "swagger-resources"), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[
			  {
				"name": "2.X版本",
				"url": "/v2/api-docs",
				"swaggerVersion": "2.0",
				"location": "v2"
			  }
			]`))
	})

	http.HandleFunc(
		fmt.Sprintf("%s%s", prefix, "v2/api-docs"),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write(f)
		},
	)

	// 注册 WebDAV 路由
	http.Handle(prefix, http.StripPrefix(prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Knife4jSwagHandler(prefix).ServeHTTP(w, r)
	})))

	// 启动 HTTP 服务器
	log.Println("接口文档地址:http://localhost:8088/api/v1/swagger/index.html")
	if err := http.ListenAndServe(":8088", nil); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

func Test_Embed(t *testing.T) {
	// 读取 static 目录中的 index.html 文件
	data, err := content.ReadFile("resource/doc.html")
	if err != nil {
		log.Fatalf("读取文件失败: %v", err)
	}

	log.Println(string(data))
}
