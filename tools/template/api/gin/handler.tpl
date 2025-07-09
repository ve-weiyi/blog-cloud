package {{.Package}}

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/response"
	{{- range .Imports}}
    {{.}}
    {{- end}}
)

type {{.Group}}Controller struct {
	svcCtx *svctx.ServiceContext
}

func New{{.Group}}Controller(svcCtx *svctx.ServiceContext) *{{.Group}}Controller {
	return &{{.Group}}Controller{
		svcCtx: svcCtx,
	}
}

{{- range .GroupRoutes}}
{{- $prefix := .Prefix}}
{{- range .Routes}}

// @Tags		{{$.Group}}
// @Summary		"{{.Doc}}"
// @accept		application/json
// @Produce		application/json
    {{- if .Request }}
// @Param		data	body		{{commentTypes .Request}}		true	"请求参数"
    {{- end }}
    {{- if .Response }}
// @Success		200		{object}	response.Body{data={{commentTypes .Response}}}	"返回信息"
    {{- end }}
// @Router		{{$prefix}}{{.Path}} [{{.Method}}]
func (s *{{$.Group}}Controller) {{.Handler}}(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	{{- if .Request }}
	var req {{pkgTypes .Request}}
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
    {{- end }}

    {{ if .Response }} data ,err : {{- else -}} err {{- end -}}
    = logic.New{{$.Group}}Logic(s.svcCtx).{{.Handler}}(reqCtx
    {{- if .Request }}, req {{- else -}} {{- end -}})
    if err != nil {
        response.ResponseError(c, err)
        return
    }

	response.ResponseOk(c, {{- if .Response }} data {{- else -}} nil {{- end }})
}

{{- end}}
{{- end}}
