# 使用docker部署服务

## 1. 安装docker

启动docker
docker-compose -f docker-compose.yaml up -d



## b部署nacos、prometheus、grafana

1. 启动nacos
```shell
docker-compose -f nacos/docker-compose.yaml up -d
```

nacos服务开启prometheus监控
首先进入nacos容器添加一行配置,配置文件路径在/home/nacos/conf/application.properties
vim /home/nacos/conf/application.properties

找个位置添加这行配置  暴露metrics数据
management.endpoints.web.exposure.include=*

访问nacos的prometheus数据
http://120.79.136.81:8848/nacos/actuator/prometheus

2. 启动prometheus
```shell
docker-compose -f prometheus/docker-compose.yaml up -d
```

查看看板
http://120.79.136.81:9090/targets
