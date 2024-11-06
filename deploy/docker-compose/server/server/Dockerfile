# 使用golang:alpine作为基础镜像
FROM golang:alpine AS builder

# 设置当前阶段标签，当前阶段是用于构建 Golang 项目的
LABEL stage=gobuilder

# 用阿里云的镜像地址替换默认的Alpine Linux软件包下载地址
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 设置Golang环境变量
ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

# 安装必要的依赖
RUN apk update --no-cache
RUN apk add --no-cache git
RUN apk add --no-cache tzdata

# 设置新构建容器的工作目录
WORKDIR /

# 下载项目
RUN git init
RUN git clone https://github.com/ve-weiyi/ve-blog-golang.git

# 切换到项目目录
WORKDIR /ve-blog-golang/server

# 切换到目标分支并拉取最新代码
RUN git checkout feature/blog
RUN go mod tidy

# 构建Golang运行程序到 /app 目录下
RUN go build -ldflags="-s -w" -o /app/server/main main.go

# 使用golang:alpine作为构建镜像,上面安装了git,go等运行环境，方便编译项目。直接在golang:alpine上运行时，image的体积约为2.5g
# 使用scratch作为基础镜像,scratch用于构建轻量级的、无操作系统依赖的容器镜像。使用scratch镜像运行，image的体积约为50m

FROM scratch

#维护者信息
LABEL maintainer="791422171@qq.com"

# 设置语言环境为en_US.utf8
ENV LANG=en_US.utf8
# 设置时区环境变量
ENV TZ Asia/Shanghai

# 设置新新容器的工作目录
WORKDIR /app/server

# 复制启动文件到新容器中
COPY --from=builder /app/server/main /app/server/main


# 设置启动命令(使用config.yaml)
#CMD ["./main","api", "-c=file", "-f=./config.yaml"]

# 设置启动命令(使用nacos作为配置中心)
CMD ["./main","api", "-c=nacos","--n-namespace=test"]
