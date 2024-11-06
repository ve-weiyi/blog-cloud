## 接口参数解析器


OPENAPI3.0介绍
https://www.cnblogs.com/liaozibo/p/openapi-2.html
https://blog.csdn.net/qq_41971087/article/details/126065322
https://www.cnblogs.com/yaohl0911/p/14567915.html

OpenAPI 中文文档：
https://openapi.apifox.cn/

OpenAPI 2.0 和OpenApi3.0比较
https://koca.szkingdom.com/forum/t/topic/1103

校验OpenAPI文档语法是否正确：
https://swagger.io/tools/swagger-editor/

Q：OpenAPI和RESTful API区别？
A：refer：https://www.cnblogs.com/origin-zy/p/17541703.html

### 1. 介绍

#### SwaggerParser ([swagger.go](tools/parserx/swagger.go))

解析swagger.json，提取接口定义。包括接口路径、接口方法、接口参数、接口返回值等信息。


####  ApiParser ([ast.go](parserx/ast.go))
使用golang ast语法树，解析路径下的接口注释，提取接口定义。包括接口路径、接口方法、接口参数、接口返回值等信息。

### 2. 使用方式

见 [swagger_test.go](swagger_test.go)
