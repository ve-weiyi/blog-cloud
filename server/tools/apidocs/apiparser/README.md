## 接口参数解析器


### 1. 介绍

#### SwaggerParser ([swagger.go](apiparser%2Fswagger.go))

解析swagger.json，提取接口定义。包括接口路径、接口方法、接口参数、接口返回值等信息。


####  ApiParser ([ast.go](apiparser%2Fast.go))
使用golang ast语法树，解析路径下的接口注释，提取接口定义。包括接口路径、接口方法、接口参数、接口返回值等信息。

### 2. 使用方式

见 [swagger_test.go](swagger_test.go)
