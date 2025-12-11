export const {{ .Name }}Perm = {

{{- range $name, $group := .Groups }}
{{- range $group }}
  {{ $prefix := .Prefix}}

  {{- range .Routes}}
  {{ .Handler }}: "{{$prefix}}{{.Path}}",
  {{- end }}

{{- end }}
{{- end }}
}
