Docker--harbor私有仓库部署与管理

https://blog.csdn.net/m0_62231324/article/details/134381529

1. 登录harbor

```sh
docker login veweiyi.cn:10443/ -u admin -p Harbor12345 
```

2. 给镜像打标签

   docker tag <本地镜像名>:<版本> <Harbor地址>/<项目名>/<镜像名>:<版本>

```sh
docker tag blog-rpc:latest veweiyi.cn:10443/blog/blog-rpc:latest
```

3. 推送镜像

   docker push <Harbor地址>/<项目名>/<镜像名>:<版本>

```sh
docker push veweiyi.cn:10443/blog/blog-rpc:latest
```
