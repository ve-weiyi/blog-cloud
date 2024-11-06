{{- range .ImportPkgPaths -}}
{{.}}
{{ end -}}

{{- range .GroupRoutes}}
{{- $prefix := .Prefix}}
{{- range .Routes}}
/** {{ .Summery }} */
export function {{ .Handler }}(

{{- if .Request }}data?: {{ .Request }}{{ end -}}

): Promise<IApiResponse<{{.Response}}>> {
  return request({
    url: "{{$prefix}}{{.Path}}",
    method: "{{.Method}}",
    {{ if .Request }}data: data,{{ end }}
  });
}
{{ end -}}
{{ end -}}
