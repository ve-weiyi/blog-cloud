package {{.Package}}

import (
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/response"
	{{- range .Imports}}
    {{.}}
    {{- end}}
)

type {{.Group}}Service struct {
	svcCtx *svctx.ServiceContext
}

func New{{.Group}}Service(svcCtx *svctx.ServiceContext) *{{.Group}}Service {
	return &{{.Group}}Service{
		svcCtx: svcCtx,
	}
}

{{- range .GroupRoutes}}
{{- range .Routes}}
// {{.Doc}}
{{- if .Request }}
func (s *{{$.Group}}Service) {{.Handler}}(reqCtx *request.Context, in {{pkgTypes .Request}}) (out {{pkgTypes .Response}}, err error) {
    // todo

    return
}
{{- else }}
func (s *{{$.Group}}Service) {{.Handler}}(reqCtx *request.Context) (err error) {
    // todo

    return
}
{{- end }}

{{- end}}
{{- end}}
