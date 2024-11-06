package {{.Package}}

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/response"
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

    data, err := service.New{{$.Group}}Service(s.svcCtx).{{.Handler}}(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
    {{- else }}
    err = service.NewWebsocketService(s.svcCtx).WebSocket(reqCtx)
    if err != nil {
        response.ResponseError(c, err)
        return
    }
    {{- end }}

    {{- if .Response }}
	response.ResponseOk(c, data)
    {{- else }}
    response.ResponseOk(c, nil)
    {{- end }}
}

{{- end}}
{{- end}}
