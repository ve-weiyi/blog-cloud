package {{.Package}}

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/response"
	{{- range .Imports}}
    {{.}}
    {{- end}}
)

type {{.Group}}Logic struct {
	svcCtx *svctx.ServiceContext
}

func New{{.Group}}Logic(svcCtx *svctx.ServiceContext) *{{.Group}}Logic {
	return &{{.Group}}Logic{
		svcCtx: svcCtx,
	}
}

{{- range .GroupRoutes}}
{{- range .Routes}}
// {{.Doc}}
func (s *{{$.Group}}Logic) {{.Handler}}(reqCtx *request.Context
{{- if .Request }}, in{{pkgTypes .Request}} {{- else -}}  {{- end }}) (
{{- if .Response }}out {{pkgTypes .Response}}, {{- else -}}  {{- end }}err error) {
    // todo

    return
}

{{- end}}
{{- end}}
