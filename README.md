
<div align=center>
<img src="https://mms1.baidu.com/it/u=2815887849,1501151317&fm=253&app=138&f=JPEG" width=300" height="300" />
</div>
<div align=center>
<img src="https://img.shields.io/badge/golang-1.20-blue"/>
<img src="https://img.shields.io/badge/gin-1.9.0-lightBlue"/>
<img src="https://img.shields.io/badge/gorm-1.24.7-red"/>
<img src="https://img.shields.io/badge/redis-9.0.2-brightgreen"/>
<img src="https://img.shields.io/badge/swagger-v1.5.3-green"/>

</div>

联系方式：qq 791422171

简体中文

## 1. 基本介绍

### 1.1 项目介绍

ve-blog 是一个前后端分离的博客系统，项目采用了前后端分离的架构进行开发。前后端通过RESTful API进行数据交互。  

博客前台展示页面使用 vite4+vue3+vuetify 开发。  

博客后台管理系统使用 vite4+vue3+element-plus 开发。  

博客后端服务使用golang语言+ gin框架开发,jwt鉴权+RBAC接口权限管理。    
接入了nacos、mysql、redis、rabbitmq、chatgpt服务，可以使用 飞书、QQ、微博、微信 快速登录。   
使用 cobra、viper、gorm、zap 等golang主流框架。   
拥有一套自动化代码生成工具，可以一键完成数据库表到[增删改查]接口的编写，减少重复工作。


[博客展示前台项目地址](https://github.com/ve-weiyi/ve-blog-vite)

[博客管理后台项目地址](https://github.com/ve-weiyi/ve-admin-vite)

[博客后端服务项目地址](https://github.com/ve-weiyi/ve-blog-golang)

### 预览页面

![img.jpg](images%2Fimg.jpg)

![img_1.jpg](images%2Fimg_1.jpg)
## 项目启动

### 1.拉取submodule
```shell
git submodule update
```

### 2.拉取golang依赖库
```shell
cd server
go mod tidy
```

### 3.修改配置文件

在server目录下新增config.yaml文件，文件内容参考[config.default.yaml](server%2Fconfig.default.yaml)  
修改config.yaml配置文件中的配置信息，然后执行以下命令启动项目

### 4.项目本地启动
确保安装了mysql、redis、rabbitmq 等服务,如果没有安装可以在[server.go](server%2Fcmd%2Fserver.go)OnInitialize 函数中，注释掉相关配置。


打开项目启动目录
```shell
cd server
```

使用本地配置文件yaml启动项目
```shell
go run main.go server --config=./config.yaml 
```

初始化数据库
```shell
go run main.go migrate --action=create --host=127.0.0.1 --port=3306 --username=root --password=mysql7914 -n=blog-veweiyi  --file=./blog-veweiyi.sql
```

使用本地配置文件yaml启动项目
```shell
go run main.go server --config=./config.yaml 
```

使用nacos的配置文件启动项目
```shell
go run main.go server --use-nacos=true --n-ns=test
```


## 5.其他

### 1.格式化

格式化代码
```shell
go fmt ./...
```

格式化导入包

```shell
go get golang.org/x/tools/cmd/goimports
go install golang.org/x/tools/cmd/goimports
goimports -w .
./scripts/gofmt.sh
```

格式化swagger注解

```shell
swag fmt
```

```shell
swag init
```


### 2.cobra使用

1. 安装cobra

```shell
go get -u github.com/spf13/cobra/cobra
```

2. 初始化项目

```shell
cobra-cli init
cobra-cli init --author "791422171@qq.com"
cobra-cli init --license apache
```

3. 添加命令

```shell
cobra-cli add version
cobra-cli add migrate
```

```shell
go run main.go migrate -h
go run main.go migrate --help

go run main.go migrate --action=reset
go run main.go server --use-nacos=true --n-ns=test
```
