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

练手项目，兴趣是第一驱动力。工作之余持续更新，持续学习新技术。

简体中文

## 1. 基本介绍

### 1.1 项目介绍

ve-blog 是一个前后端分离的博客系统，项目采用了前后端分离的架构进行开发。前后端通过RESTful API进行数据交互。

博客前台展示页面使用 vite4+vue3+vuetify 开发。

博客后台管理系统使用 vite4+vue3+element-plus 开发。

博客后端服务使用golang语言。项目有两种框架：gin和go-zero

+ server/[README.md](server/README.md) 使用[gin](https://github.com/gin-gonic/gin)
  框架开发。gin框架的优点是轻量、快速、易用，适合快速开发API接口，适用于golang后端初学者。
+ zero/[README.md](server/README.md) 使用[go-zero](https://github.com/zeromicro/go-zero)
  框架开发。go-zero框架的优点是功能齐全，自带了多种中间件，适合大型项目和复杂业务的解耦。

1.主流框架。cobra、viper、gorm、zap。   
2.主流中间件。nacos、mysql、redis、rabbitmq、chatgpt服务。    
3.第三方登录的解决方案。可以使用 飞书、QQ、微博、微信 快速登录。    
4.jwt鉴权+RBAC权限管理。可以对用户和角色进行权限管理，可以对菜单页面和接口访问进行权限控制。
5.提供了接口文档。通过导入swagger.json可以在apifox和postman查看和调试接口。   
6.提供了自动化代=码生成工具。从数据库表一键生成业务代码。

[博客展示前台项目地址](https://github.com/ve-weiyi/ve-blog-vite)

[博客管理后台项目地址](https://github.com/ve-weiyi/ve-admin-vite)

[博客后端服务项目地址](https://github.com/ve-weiyi/ve-blog-golang)

### 预览页面

![img.jpg](assets/images/img.jpg)

![img_1.jpg](assets/images/img_1.jpg)

## todo

- [ ] 拆分rpc服务和api服务

## 版本发布记录

**格式**

主版本号.次版本号.修订号
