# 格式化 .api文件
goctl api format -dir ./

# 根据 api 文件生成 swagger 文档。
goctl api swagger --api app.api --dir ../docs

# 根据 api 文件生成 业务代码
goctl api go -api app.api -dir ../ --style go_zero --home ../../../../.goctl/template --type-group
