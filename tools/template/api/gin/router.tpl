package {{.Package}}

import (
	{{- range .Imports}}
    {{.}}
    {{- end}}
)

type {{.Group}}Router struct {
	svcCtx *svctx.ServiceContext
}

func New{{.Group}}Router(svcCtx *svctx.ServiceContext) *{{.Group}}Router {
	return &{{.Group}}Router{
		svcCtx: svcCtx,
	}
}

func (s *{{.Group}}Router) Register(r *gin.RouterGroup) {
    {{- range .GroupRoutes}}
    // {{$.Group}}
    // {{.Middleware}}
    {
        group := r.Group("{{.Prefix}}")
        {{- range .Middleware}}
        group.Use(s.svcCtx.{{.}})
        {{- end }}

        h := handler.New{{$.Group}}Controller(s.svcCtx)
        {{- range .Routes}}
        // {{.Doc}}
        group.{{.Method}}("{{.Path}}", h.{{.Handler}})
        {{- end}}
    }

    {{- end}}
}

