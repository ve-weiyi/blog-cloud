package {{ .Package }}

import (
	{{- range .Imports}}
    {{.}}
    {{- end}}
)

// register your handlers here
func RegisterHandlers(r *gin.RouterGroup, svCtx *svctx.ServiceContext) {
    {{- range .Groups}}
    router.New{{.}}Router(svCtx).Register(r)
    {{- end}}
}

