// 声明 proto 语法版本，固定值
syntax = "proto3";

// proto 包名
package {{.Tag}};

// 生成 golang 代码后的包名
option go_package = "{{.Tag}}";

import "any.proto";

{{- range $key, $value := .ModelDeclares}}
message {{$value.Name}} {
    {{- range $index, $field := $value.Fields}}
      {{messageType $field.Type}} {{$field.Json}} = {{add $index 1}};{{if $field.Comment}}// {{$field.Comment}}{{end}}
    {{- end}}
}
{{end}}

service {{.Tag}}Rpc {
    {{- range .ApiDeclares}}
    // {{.Summary}}
    rpc {{.FunctionName}}({{.Request}}) returns({{.Response}});
    {{- end}}
}
