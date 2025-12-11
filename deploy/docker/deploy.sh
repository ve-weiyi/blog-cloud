# 拉取镜像
docker pull ghcr.io/ve-weiyi/blog-rpc:latest

docker pull ghcr.io/ve-weiyi/blog-api:latest

docker pull ghcr.io/ve-weiyi/admin-api:latest


# 部署镜像
docker run -d \
--name blog-rpc \
--restart always \
-p 9999:9999 \
-v ./runtime:/app/runtime \
ghcr.io/ve-weiyi/blog-rpc:latest \
  ./blog \
  -nacos-host "veweiyi.cn" \
  -nacos-port "8848" \
  -nacos-username "nacos" \
  -nacos-password "nacos" \
  -nacos-namespace "prod" \
  -nacos-group "veweiyi.cn" \
  -nacos-data-id "blog-rpc"


docker run -d \
--name blog-api \
--restart always \
-p 9090:9090 \
-v ./runtime:/app/runtime \
ghcr.io/ve-weiyi/blog-api:latest \
  ./blog \
  -nacos-host "veweiyi.cn" \
  -nacos-port "8848" \
  -nacos-username "nacos" \
  -nacos-password "nacos" \
  -nacos-namespace "prod" \
  -nacos-group "veweiyi.cn" \
  -nacos-data-id "blog-api"


docker run -d \
--name admin-api \
--restart always \
-p 9091:9091 \
-v ./runtime:/app/runtime \
ghcr.io/ve-weiyi/admin-api:latest \
  ./admin \
  -nacos-host "veweiyi.cn" \
  -nacos-port "8848" \
  -nacos-username "nacos" \
  -nacos-password "nacos" \
  -nacos-namespace "prod" \
  -nacos-group "veweiyi.cn" \
  -nacos-data-id "admin-api"
