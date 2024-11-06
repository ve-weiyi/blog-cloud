syntax = "v1"

import "base.api"
import "types.api"

type (
    {{- range $key, $value := .ModelDeclares}}
    {{$value.Name}} {
        {{- range $value.Fields}}
        {{.Name}} {{messageType .Type}} `json:"{{.Json}},optional"` {{if .Comment}}// {{.Comment}}{{end}}
    	{{- end}}
    }
    {{end}}
)

@server(
    prefix: /api/v1
    group: {{.Tag}}
)

service blog-api {
    {{- range .ApiDeclares}}
    @doc "{{.Summary}}"
    @handler {{.FunctionName}}
    {{.Method}} {{.Route}} ({{.Request}}) returns ({{.Response}})
    {{end}}
}
