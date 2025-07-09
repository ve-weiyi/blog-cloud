build-gin-server:
	go run blog-gin/main.go

build-gozero-blog-rpc:
	go run blog-gozero/service/rpc/blog/blog.go

build-gozero-blog-api:
	go run blog-gozero/service/api/blog/blog.go

build-gozero-admin-api:
	go run blog-gozero/service/api/admin/admin.go