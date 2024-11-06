
# 网站搭建之https配置websocket连接(wss)  
文章封面:  https://veport.oss-cn-beijing.aliyuncs.com/articles/f9fa18da262910eb13f802b003147915.jpg   
文章类型: 1   
文章分类: 网站搭建   
文章标签: [websocket springboot https]   
创建时间: 2022-02-10 10:57:51 +0800 CST   

文章内容:
## 前言
ws和wss的关系，就像http和https的关系。
当网站发布上线后，由于websocket问题，聊天室功能不能正常使用。具体问题如下：使用http非安全协议访问网站时，可以访问ws://localhost:8088/api/websocket，聊天室功能正常。但是使用https安全协议访问网站时，可以访问ws://localhost:8088/api/websocket失败

![BAC09972C5BC2F039B371679953FCDC7.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/bac09972c5bc2f039b371679953fcdc7.png)


查看报错信息可以知道，因为https不支持ws协议，这个请求被拦截了。所以我们应该想到https只支持wss协议。那么我们只简单的将websocket地址前缀改为wss？修改后依然报错：

![04B43675B833412CB6580A279614331C.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/028181b8fe4cd46c4b13b3267d3d703d.png)

出现这个问题的原因是，后端springboot并不支持https访问，因此我们需要给后端添加SSL证书。

```xml
在终端输入以下命令，生成keystore.jks证书：
keytool -genkey -alias tomcat -keyalg RSA -keystore ./keystore.jks -validity 3650 -ext san=dns:spam,ip:110.42.180.40
```

![3386BE795E832853F5A0BDC23B19C265.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/ee48a2d53494dffe911025b3c1bd1bd7.png)

其中 -keystore是证书路径，-validity 是有效期(可选)，单位为天，-ext 是额外信息(可选)，这里填了发布单位。

然后在项目的配置文件yml里设置

![image20220210104501701.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/4648a2b0bcf1f0417fd56ab5cef7eefa.png)


当可以使用https访问接口时，说明已经配置成功

![B683E1968A284CD698005E6C84ED786A.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/2a93889977a7a0a8053ca8fc15afb0ae.png)

这个时候就可以使用wss连接websocket啦，[websocket在线测试地址](http://www.jsons.cn/websocket/)。  
(注意！！不要使用Chrome浏览器测试，因为我们自己颁发的证书是不被Chrome认可的！！如果想使用)  

**推荐使用官方颁发的SLL证书部署websocket和springboot。具体实现方式可查阅参考链接4和参考链接5。
**

参考链接:

1.[spring boot 集成 websocket 的四种方式](https://www.cnblogs.com/kiwifly/p/11729304.html)
2.[Springboot2.1集成WebSocket配置wss访问](https://blog.csdn.net/u012977315/article/details/84944708)
3.[websocket在线测试地址](http://www.jsons.cn/websocket/)
4.[Nginx 服务器 SSL 证书安装部署[腾讯云]](https://cloud.tencent.com/document/product/400/35244)
5.[在Spring Boot中配置ssl证书实现https](https://www.jianshu.com/p/eb52e0f5ee85)
