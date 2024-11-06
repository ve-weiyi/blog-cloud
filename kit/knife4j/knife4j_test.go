package knife4j

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestDir(t *testing.T) {
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

	prefix := "/admin_api/v1/swagger/"

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
	log.Println("WebDAV 服务器启动于 http://localhost:9091")
	if err := http.ListenAndServe(":9091", nil); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
