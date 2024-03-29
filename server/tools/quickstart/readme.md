## 快速代码构建器

根据数据库模型生成 router、controller、service、repository、model 代码，并自动注册到 gin 框架中。

### 1. 介绍


inject 使用ast在文件中指定位置注入代码
invent 使用template生成代码文件


与go-zero框架对应层级

controller --- api.handler 

service --- api.logic 

repository--- rpc.handler、logic 

entity--- model 
