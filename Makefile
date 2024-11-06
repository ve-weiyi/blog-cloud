build-server:
	go run blog-gin/main.go

build-blog-rpc:
	go run blog-gozero/service/rpc/blog/blog.go

build-blog-api:
	go run blog-gozero/service/api/blog/blog.go

build-admin-api:
	go run blog-gozero/service/api/admin/admin.go