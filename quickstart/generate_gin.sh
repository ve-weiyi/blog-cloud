
# api生成dto文件
go run main.go api dto \
-f='../zero/service/api/blog/proto/blog.api' \
-t='./resource/gin/template' \
-o='../server/api/blog'  \
-n='%s.go'

# api生成service文件
go run main.go api service \
-f='../zero/service/api/blog/proto/blog.api' \
-t='./resource/gin/template' \
-c='github.com/ve-weiyi/ve-blog-golang/server/svctx'  \
-o='../server/api/blog'  \
-n='%s.sv.go'

# api生成controller文件
go run main.go api controller \
-f='../zero/service/api/blog/proto/blog.api' \
-t='./resource/gin/template' \
-c='github.com/ve-weiyi/ve-blog-golang/server/svctx'  \
-o='../server/api/blog'  \
-n='%s.ctl.go'

# api生成router文件
go run main.go api router \
-f='../zero/service/api/blog/proto/blog.api' \
-t='./resource/gin/template' \
-c='github.com/ve-weiyi/ve-blog-golang/server/svctx'  \
-o='../server/api/blog'  \
-n='%s.rt.go'


###############################

# api生成dto文件
go run main.go api dto \
-f='../zero/service/api/admin/proto/admin.api' \
-t='./resource/gin/template' \
-o='../server/api/admin'  \
-n='%s.go'

# api生成service文件
go run main.go api service \
-f='../zero/service/api/admin/proto/admin.api' \
-t='./resource/gin/template' \
-c='github.com/ve-weiyi/ve-blog-golang/server/svctx'  \
-o='../server/api/admin'  \
-n='%s.sv.go'

# api生成controller文件
go run main.go api controller \
-f='../zero/service/api/admin/proto/admin.api' \
-t='./resource/gin/template' \
-c='github.com/ve-weiyi/ve-blog-golang/server/svctx'  \
-o='../server/api/admin'  \
-n='%s.ctl.go'

# api生成router文件
go run main.go api router \
-f='../zero/service/api/admin/proto/admin.api' \
-t='./resource/gin/template' \
-c='github.com/ve-weiyi/ve-blog-golang/server/svctx'  \
-o='../server/api/admin'  \
-n='%s.rt.go'
