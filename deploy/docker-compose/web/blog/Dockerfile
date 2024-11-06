# 使用Node.js的最新官方镜像作为基础镜像，也可以固定版本
FROM node:20-alpine as build-stage

# 设置新构建容器的工作目录
WORKDIR /

# 安装Git
RUN yum  update && yum install -y git

# 下载项目
RUN git init
RUN git clone https://github.com/ve-weiyi/ve-blog-vite.git

# 切换到项目目录
WORKDIR /ve-blog-vite

# 切换到目标分支并拉取最新代码
RUN git checkout feature/blog

# 安装依赖
RUN npm install -g pnpm

# 执行 Vite 构建命令，生成 dist 目录
RUN pnpm install
RUN pnpm run build

# 使用 Nginx 镜像作为运行时镜像

FROM nginx:stable-alpine as production-stage

# 将 Vite 项目的 dist 目录复制到 Nginx 的默认静态文件目录
COPY --from=build-stage /ve-blog-vite/blog /usr/share/nginx/html

# 暴露容器的 80 端口
EXPOSE 80

# Nginx 会在容器启动时自动运行，无需手动设置 CMD
