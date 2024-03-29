package apiparser

// https://www.cnblogs.com/liaozibo/p/openapi-2.html

// Swagger represents the root Swagger .
type SwaggerDefinition struct {
	Swagger             string                 `json:"swagger,omitempty,omitempty" description:"说明 Swagger 规范版本。可被工具用于解析文档。该值必须是 \"2.0\"。" required:"true"`
	Info                *Info                  `json:"info,omitempty" description:"文档元数据。可被工具使用。" required:"true"`
	Host                string                 `json:"host,omitempty" description:"提供 API 服务的主机名或 IP。【必须】是主机名不能包含路径。【可能】包含端口。如果没有指定，默认是文档服务的主机。"`
	BasePath            string                 `json:"basePath,omitempty" description:"API 服务的基础路径。必须以 / 开始。"`
	Schemes             []string               `json:"schemes,omitempty" description:"API 服务的传输协议。【必须】是 \"http\", \"https\", \"ws\", \"wss\"。如果没有指定则与访问文档的协议相同。"`
	Consumes            []string               `json:"consumes,omitempty" description:"全局定义可被 API 消费的 MIME 类型。可被具体接口覆盖。"`
	Produces            []string               `json:"produces,omitempty" description:"全局定义 API 返回的 MIME 类型。可被具体接口覆盖。"`
	Paths               Paths                  `json:"paths,omitempty" description:"定义路径及其支持的操作。" required:"true"`
	Definitions         Definitions            `json:"definitions,omitempty" description:"定义接口使用的数据模型"`
	Parameters          ParametersDefinitions  `json:"parameters,omitempty" description:"定义参数使用数据模型"`
	Responses           ResponsesDefinitions   `json:"responses,omitempty" description:"定义响应使用的数据模型"`
	SecurityDefinitions SecurityDefinitions    `json:"securityDefinitions,omitempty" description:""`
	Security            []SecurityRequirement  `json:"security,omitempty" description:""`
	Tags                []Tag                  `json:"tags,omitempty" description:"可被文档使用的标签对象列表，提供了一些元数据信息。该列表可被工具用于解析顺序。不必把接口使用的标签都定义在这里。标签对象必须唯一。"`
	ExternalDocs        *ExternalDocumentation `json:"externalDocs,omitempty" description:"外部文档"`
}

// Info defines metadata about the API.
type Info struct {
	Title          string   `json:"title,omitempty" description:"API 应用的名称" required:"true"`
	Description    string   `json:"description,omitempty" description:"API 应用的简述"`
	TermsOfService string   `json:"termsOfService,omitempty" description:"提供术语说明的服务"`
	Contact        *Contact `json:"contact,omitempty" description:"联系信息"`
	License        *License `json:"license,omitempty" description:"API 证书信息"`
	Version        string   `json:"version,omitempty" description:"API 版本。不要和规范版本混淆。" required:"true"`
}

// Contact represents the contact information.
type Contact struct {
	Name  string `json:"name,omitempty" description:"个人/组织的名称"`
	URL   string `json:"url,omitempty" description:"指向联系信息的链接"`
	Email string `json:"email,omitempty" description:"联系邮箱"`
}

// License defines API license information.
type License struct {
	Name string `json:"name,omitempty" description:"证书名称" required:"true"`
	URL  string `json:"url,omitempty" description:"证书链接"`
}

// Paths defines API paths and their supported operations.
type Paths map[string]PathItem //API 路径。路径【必须】以斜杠开始。该路径与 basePath 构成完整路径。

type PathItem map[string]Operation

// PathItem defines operations supported by a path.
//type PathItem struct {
//	Ref        string           `json:"$ref,omitempty"` //引用已定义的 Path Item Object
//	Get        Operation        `json:"get"`            //HTTP 方法
//	Put        Operation        `json:"put"`            //HTTP 方法
//	Post       Operation        `json:"post"`           //HTTP 方法
//	Delete     Operation        `json:"delete"`         //HTTP 方法
//	Options    Operation        `json:"options"`        //HTTP 方法
//	Head       Operation        `json:"head"`           //HTTP 方法
//	Patch      Operation        `json:"patch"`          //HTTP 方法
//	Parameters []ParameterOrRef `json:"parameters"`     // 定义该路径下所有操作支持的参数列表，可被具体的操作覆盖。
//}

// Operation defines a specific API operation.
type Operation struct {
	Tags         []string               `json:"tags,omitempty" description:"标签列表，可用于逻辑分组。"`
	Summary      string                 `json:"summary,omitempty" description:"概括接口功能。【应该】少于 120 个字符。"`
	Description  string                 `json:"description,omitempty" description:"详细描述接口功能。"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty" description:"引用外部文档"`
	OperationID  string                 `json:"operationId,omitempty" description:"接口的唯一标识符，【必须】是全局唯一的。"`
	Consumes     []string               `json:"consumes,omitempty" description:"接口消费的 MIME 类型"`
	Produces     []string               `json:"produces,omitempty" description:"接口返回的 MIME 类型"`
	Parameters   []ParameterOrRef       `json:"parameters,omitempty" description:"接口参数列表，可以直接定义 Parameter  或者引用 parameters 定义的参数对象。"`
	Responses    map[string]Response    `json:"responses,omitempty" description:"可能返回响应列表" required:"true"`
	Schemes      []string               `json:"schemes,omitempty" description:"传输协议，\"http\", \"https\", \"ws\", \"wss\"。"`
	Deprecated   bool                   `json:"deprecated,omitempty" description:"是否过时，默认是 false。"`
	Security     []map[string][]string  `json:"security,omitempty" description:""`
}

// ParameterOrRef represents either a parameter  or a reference.
type ParameterOrRef struct {
	Reference string `json:"$ref,omitempty"`
	*Parameter
}

// Parameter defines an API parameter.
type Parameter struct {
	Name        string `json:"name,omitempty" description:"参数名。大小写敏感。" required:"true"`
	In          string `json:"in,omitempty" description:"参数位置。可选值是 \"query\", \"header\", \"path\", \"formData\", \"body\"。" required:"true"`
	Description string `json:"description,omitempty" description:"参数的简述"`
	Required    bool   `json:"required,omitempty" description:"参数是否是强制的。对于路径参数，该值是必填的并且【必须】是 true。其他【可能】false。" required:"true"`

	// 如果 in 的值是 body，还有一个属性名
	Schema *SchemaObject `json:"schema,omitempty" description:"定义 body 的一些属性"`

	// 如果 in 的值不是 body，有一系列属性名来定义参数：
	*SchemaObject
}

// SchemaObject defines a schema .
type SchemaObject struct {
	*Items

	AllOf []AllOf `json:"allOf,omitempty" description:""`
}

// AllOf 是引用类型或者
type AllOf struct {
	Ref string `json:"$ref,omitempty" description:""`

	Type       string                  `json:"type,omitempty" description:"参数类型。该值【必须】是 \"string\", \"number\", \"integer\", \"boolean\", \"array\", \"file\"。"`
	Properties map[string]SchemaObject `json:"properties,omitempty" description:""`
}

// Items defines the type of items in an array.
type Items struct {
	Ref string `json:"$ref,omitempty" description:""`

	Type                 string        `json:"type,omitempty" description:"参数类型。该值【必须】是 \"string\", \"number\", \"integer\", \"boolean\", \"array\", \"file\"。"`
	Description          string        `json:"description,omitempty" description:"参数简述"`
	Example              interface{}   `json:"example,omitempty" description:"参数示例"`
	Format               string        `json:"format,omitempty" description:"具体参数类型"`
	AllowEmptyValue      bool          `json:"allowEmptyValue,omitempty" description:"是否允许空值。只有 in 是 query 或 formData 才会生效。如果允许，可以只传参数名或传空值。默认是 false。"`
	Items                *Items        `json:"items,omitempty" description:"描述数组元素的类型" required:"true"`
	CollectionFormat     string        `json:"collectionFormat,omitempty" description:"数组格式。可选择是：csv，逗号分隔。ssv：空格分隔。tsv：反斜杠分隔。pipes：管道分隔。multi：。默认是 csv。"`
	Default              interface{}   `json:"default,omitempty" description:"参数默认值"`
	Maximum              float64       `json:"maximum,omitempty" description:""`
	ExclusiveMaximum     bool          `json:"exclusiveMaximum,omitempty" description:""`
	Minimum              float64       `json:"minimum,omitempty" description:""`
	ExclusiveMinimum     bool          `json:"exclusiveMinimum,omitempty" description:""`
	MaxLength            int           `json:"maxLength,omitempty" description:""`
	MinLength            int           `json:"minLength,omitempty" description:""`
	Pattern              string        `json:"pattern,omitempty" description:""`
	MaxItems             int           `json:"maxItems,omitempty" description:""`
	MinItems             int           `json:"minItems,omitempty" description:""`
	UniqueItems          bool          `json:"uniqueItems,omitempty" description:""`
	Enum                 []interface{} `json:"enum,omitempty" description:""`
	MultipleOf           float64       `json:"multipleOf,omitempty" description:""`
	AdditionalProperties interface{}   `json:"additionalProperties,omitempty" description:""`
}

// Responses is a container for response s.
type Responses struct {
	Default   Response            `json:"default,omitempty" description:"Reference "`
	Responses map[string]Response `json:"responses,omitempty" description:""`
}

// Response defines a response .
type Response struct {
	Description string         `json:"description,omitempty" description:"必填。对响应的简述。" required:"true"`
	Schema      *SchemaObject  `json:"schema,omitempty" description:"响应结构定义。该值可以是原始类型，数组或对象。如果没有指定该属性，则意味着该响应没有内容。如果它的 \"type\" 是 \"file\"，则【应该】同时声明 produces 的 MIME 类型。"`
	Headers     *HeaderObject  `json:"headers,omitempty" description:"响应头列表"`
	Examples    *ExampleObject `json:"examples,omitempty" description:"响应消息示例"`
}

type HeaderObject map[string]Header

type Header struct {
	Description string `json:"description,omitempty"` // 对响应头的简述。
	SchemaObject
}

// ExampleObject represents an example response.
type ExampleObject struct {
	ApplicationJSON ExampleJSON `json:"application/json,omitempty" description:"该属性名【必须】是接口 produces 属性指定的 MIME 值。并且【应该】和实际需要返回的内容一致。"`
}

// ExampleJSON represents an example JSON response.
type ExampleJSON struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Color  string `json:"color"`
	Gender string `json:"gender"`
	Breed  string `json:"breed"`
}

// Tag represents an API tag.
type Tag struct {
	Name         string                 `json:"name,omitempty" description:"必填。标签名。" required:"true"`
	Description  string                 `json:"description,omitempty" description:"简述。"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty" description:"外部文档。"`
}

// ExternalDocumentation represents external documentation.
type ExternalDocumentation struct {
	Description string `json:"description,omitempty" description:"对外部文档的简短描述"`
	URL         string `json:"url,omitempty" description:"外部"`
}

// Definitions 定义接口使用的数据模型。
type Definitions map[string]Definition

type Definition struct {
	SchemaObject
	Properties map[string]SchemaObject `json:"properties,omitempty" description:""`
}

// ParametersDefinitions 可以被接口复用的参数对象。
type ParametersDefinitions map[string]Parameter

// ResponsesDefinitions 可以被接口复用的响应对象。
type ResponsesDefinitions map[string]Response

// SecurityDefinitions 安全性定义对象。
type SecurityDefinitions map[string]SecurityScheme

// SecurityScheme 安全方案对象。
type SecurityScheme struct {
	Type             string            `json:"type,omitempty"`
	Description      string            `json:"description,omitempty"`
	Name             string            `json:"name,omitempty"`
	In               string            `json:"in,omitempty"`
	Flow             string            `json:"flow,omitempty"`
	AuthorizationURL string            `json:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
}

// SecurityRequirement 安全需求对象。
type SecurityRequirement map[string][]string
