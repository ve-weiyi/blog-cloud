build-server:
	go run server/main.go

build-blog-rpc:
	go run zero/service/rpc/blog/blog.go

build-blog-api:
	go run zero/service/api/blog/blog.go

build-admin-api:
	go run zero/service/api/admin/admin.go