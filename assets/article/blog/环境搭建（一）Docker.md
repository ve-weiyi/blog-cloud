
# 前言

## 什么是Docker ？

```go
Docker 是一个开源的容器化平台，它使开发者能够创建、部署和运行应用程序在隔离的环境中（称为“容器”）中。容器是一种虚拟化技术，但它不同于传统的虚拟机，因为它不需要完整的操作系统，而是共享宿主操作系统的内核，从而更加轻量、高效。

Docker 的主要特点包括：

轻量和高效：Docker 容器不需要完整的操作系统，而是只包含运行应用所需的代码和依赖，这使得容器启动和运行速度很快，且资源消耗小。

跨平台兼容：Docker 容器可以在任何支持 Docker 的平台上运行，从而实现“构建一次，随处运行”（Build Once, Run Anywhere）的效果。这对于开发、测试和生产环境的迁移非常有帮助。

版本控制：Docker 使用镜像（Images）来定义容器内容，开发者可以创建镜像的不同版本，从而轻松进行版本管理和更新。

隔离和安全：Docker 容器与宿主系统之间有一定的隔离，使得应用更安全，并且不同的容器之间也相互隔离，防止了资源干扰。

Docker的主要组件
镜像 (Image)：Docker 镜像是一个只读的模板，用于创建 Docker 容器。每个镜像包含应用程序运行所需的所有内容，比如代码、库、依赖项等。

容器 (Container)：容器是镜像的一个运行实例，可以视为一个独立的、隔离的进程，它包含应用程序和运行所需的环境。

Docker Hub：一个云端的 Docker 镜像仓库，允许用户共享和下载镜像。

Docker的工作流程
使用 Dockerfile 定义镜像，描述应用程序所需的依赖和配置。
使用 docker build 命令创建镜像。
使用 docker run 命令启动容器，在隔离的环境中运行应用。
使用 docker push 和 docker pull 上传和下载镜像。
Docker 的应用广泛，从本地开发环境设置到大规模的生产环境应用都可以使用。它特别适合微服务架构，简化了应用的部署、管理和扩展。
```

以上内容来源于ChatGPT。

我以个人使用体验，我抽象的对docker进行比喻：

假设你需要运行一个项目，你旧需要搭建运行环境，如安装Mysql、安装Redis、安装Java运行环境等等。

我们知道安装各种环境需要花费时间去学习安装命令，而且很容易出现各种问题，如版本不兼容、依赖冲突等等。

并且你在安装之后，还会对你操作系统的配置产生影响，如端口占用、环境变量配置等等。例如你在安装Java时，需要配置环境变量，才能使用Java命令。

Docker就是一个解决方案，它将你的项目运行环境打包成一个容器，你只需要安装Docker，然后运行你的容器，就可以运行你的项目了。

**Docker像是一个公司，你的各类环境就是其中的一个team（例如mysql、redis），你不需要去考虑这个team中存在哪些人**。当需要这个team时，你就招聘这个team，当不需要时，你就解散这个team，不会产生任何影响。


# 用途

Docker用于快速安装 Mysql、Redis、Rabbitmq等等服务。

# 安装

centos7安装Docker
```shell
sudo yum update

sudo yum install git

sudo yum install docker

sudo yum install docker-compose
```

## 安装Docker-Desktop
Docker-Desktop 是Docker官方提供的桌面版Docker，可以使用可视化界面管理Docker，它支持Windows和macOS。

## 检测安装结果
```shell
docker --version
```
