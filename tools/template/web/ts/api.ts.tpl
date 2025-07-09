{{- range .ImportPkgPaths -}}
{{.}}
{{ end }}
export const {{ .Name }}API = {
{{- range .TsApiGroups }}
{{- $prefix := .Prefix -}}
{{- range .Routes}}
  /** {{ .Summery }} */
  {{ .Handler }}Api(
  {{- if .Request }}data?: {{ .Request }}{{ end -}}
  ): Promise<IApiResponse<{{.Response}}>> {
    return request({
      url: "{{$prefix}}{{.Path}}",
      method: "{{.Method}}",
      {{ if .Request }}data: data,{{ end }}
    });
  },
{{ end -}}
{{ end -}}
};
