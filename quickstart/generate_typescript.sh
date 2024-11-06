
# 生成api.ts
go run main.go web typescript \
-n='%v.ts' \
-t='./resource/web/template' \
-o='./runtime/blog/api'  \
-m='api' \
-i='IApiResponse' \
-f='/Users/weiyi/Github/ve-blog-golang/zero/service/api/blog/proto/blog.api'

# 生成api.ts
go run main.go web typescript \
-n='%v.ts' \
-t='./resource/web/template' \
-o='./runtime/admin/api'  \
-m='api' \
-i='IApiResponse' \
-f='/Users/weiyi/Github/ve-blog-golang/zero/service/api/admin/proto/admin.api'
