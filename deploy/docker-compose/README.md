# 使用docker部署服务


服务器名称	作用
deploy-server.com	部署 gitlab、jenkins、harbor（预先装好 docker、docker-compose）
srv-data.com	    部署 mysql、redis、es 等等，模拟独立环境,k8s 内部连接到此服务器
nginx-gateway.com	网关，独立于 k8s 集群外部
k8s                 集群	K8s 集群
