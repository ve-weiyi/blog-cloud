# gin框架运行的服务


## server项目结构

```shell
├── api 服务接口处理逻辑,根据 router,controller,service,repository,model 分层
│   └── blog 博客服务
│      ├── router
│      └── controller
│      ├── service
│      └── repository
│      ├── model
├── cmd 项目启动命令
├── config 配置文件
├── core 核心运行组件
├── docs swagger文档
├── infra 项目基础设施
├── initialize 项目初始化
├── svc service_context服务上下文,持有所有项目运行时的资源
└── tools 工具包,快速代码等工具
```

### 1.安装项目运行环境
- 默认你已经安装好了Golang环境
- 你
- 确保安装了 mysql、redis、rabbitmq 等服务。

如果您不知道如何安装这些服务，可以参考[docker-compose](../deploy/docker-compose)目录下的docker-compose.yaml文件，使用docker-compose启动这些服务。

如果您不会使用docker，可以参考[Docker.md](../zero/resource/docker-compose/Docker.md)文档，或者其他相关博客。

如果你不想使用这些服务，可以在[service_context.go](svctx/service_context.go) NewServiceContext 函数中，注释掉相关初始化代码。

- 初始化数据库,创建表和数据
```shell
go run main.go migrate --action=migrate \
--file=./blog-veweiyi.sql \
--host=127.0.0.1 \
--port=3307 \
--username=root \
--password=mysql7914 \
--name='blog-veweiyi'  
```

### 2.修改配置文件
在server目录下新增config.yaml文件，修改config.yaml配置文件中的配置信息，文件内容参考[config.default.yaml](server/config.default.yaml)


### 3.项目本地启动
启动方式二选一

1. 使用本地配置文件 config.yaml 启动项目
```shell
go run main.go api -c=file -f=./config.yaml
```

2. 使用nacos的配置文件启动项目
```shell
go run main.go api -c=nacos --n-namespace=test
```

### 4.其他

格式化代码
```shell
go fmt ./...
```

格式化导入包
```shell
go install golang.org/x/tools/cmd/goimports
goimports -w .
```

### 5.生成swagger文档
安装swag
```shell
go install github.com/swaggo/swag/cmd/swag@latest
```

生成swagger注解
```shell
swag init
```

格式化swagger注解

```shell
swag fmt
```
