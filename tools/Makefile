.PHONY: help deps gen-api-gin-api gen-api-gin-swagger gen-model-ddl gen-model-dsn gen-web-ts-api gen-web-ts-swagger clean

help:
	@echo "可用命令："
	@echo "  deps                  - 安装依赖"
	@echo "  gen-api-gin-api       - 使用.api文件生成 Gin 框架代码"
	@echo "  gen-api-gin-swagger   - 使用swagger.json文件生成 Gin 框架代码"
	@echo "  gen-model-ddl         - 从SQL文件生成数据库模型"
	@echo "  gen-model-dsn         - 从数据库连接生成数据库模型"
	@echo "  gen-web-ts-api        - 使用.api文件生成 TypeScript 代码"
	@echo "  gen-web-ts-swagger    - 使用swagger.json文件生成 TypeScript 代码"
	@echo "  clean                 - 清理生成的代码"

# 安装依赖
deps:
	go mod tidy

# 使用.api文件生成 Gin 框架代码
gen-api-gin-api:
	go run main.go api gin api \
		-f ../blog-gozero/service/api/blog/proto/blog.api \
		-t ./template/api/gin \
		-o ../blog-gin/api/blog \
		-c github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx \
		-n '%s.go'

# 使用swagger.json文件生成 Gin 框架代码
gen-api-gin-swagger:
	go run main.go api gin swagger \
		-f ../blog-gozero/service/api/blog/docs/blog.json \
		-t ./template/api/gin \
		-o ../blog-gin/api/blog \
		-c github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx \
		-n '%s.go'

# 从SQL文件生成数据库模型
gen-model-ddl:
	go run main.go model mysql ddl \
		-s ../blog-veweiyi-init.sql \
		-t ./template/model/model.tpl \
		-o ./runtime/model \
		-n '%v_model.go'

# 从数据库连接生成数据库模型
gen-model-dsn:
	go run main.go model mysql dsn \
		-u 'root:mysql7914@(127.0.0.1:3306)/blog-veweiyi?charset=utf8mb4&parseTime=True&loc=Local' \
		-t ./template/model/model.tpl \
		-o ./runtime/model \
		-n '%v_model.go'

# 使用.api文件生成 TypeScript 代码
gen-web-ts-api:
	go run main.go web ts api \
		-f ../blog-gozero/service/api/blog/proto/blog.api \
		-t ./template/web/ts \
		-o ./runtime/web/ts/blog/api \
		-n '%v.ts'

# 使用swagger.json文件生成 TypeScript 代码
gen-web-ts-swagger:
	go run main.go web ts swagger \
		-f ../blog-gozero/service/api/blog/docs/blog.json \
		-t ./template/web/ts \
		-o ./runtime/web/ts/blog/api \
		-n '%v.ts'

# 清理生成的代码
clean:
	rm -rf ./runtime
