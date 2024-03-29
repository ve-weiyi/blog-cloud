package apidocs

// api文档信息
type TsApiDoc struct {
	Tag              string            // 标签、分类
	ImportPkgPaths   []string          // 导入的包
	ImportModelTypes []string          // 导入的model类型
	ModelDeclares    []*TsModelDeclare // 定义的model
	ApiDeclares      []*TsApiDeclare   // 定义的api方法
}

type TsApiDeclare struct {
	Tag          string
	FunctionName string
	Summary      string

	Base   string
	Route  string
	Method string
	Header []*TsApiParam // 请求头参数
	Path   []*TsApiParam // 路径参数 path
	Query  []*TsApiParam // 查询参数 query
	Form   []*TsApiParam // 表单参数 form-data
	Body   *TsApiParam   // 请求体参数

	Request  string // 请求参数
	Response string // 响应参数
}

// 参数定义
type TsApiParam struct {
	Name        string `json:"name"`                  // 参数名
	Type        string `json:"type"`                  // 类型 object array string integer number boolean
	Description string `json:"description,omitempty"` // 描述
}

type TsModelDeclare struct {
	Name    string
	Extends []string
	Fields  []*TsModelField
}

// model属性定义信息  Name string // 属性名称
type TsModelField struct {
	Name    string // 属性名称  Name
	Type    string // 属性类型  string、int、bool、float、{UpperStartCamelName}
	Comment string // 属性的注释  属性名称
}
