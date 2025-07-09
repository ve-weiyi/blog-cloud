{{- range . -}}
{{ .Comment }}
export interface {{ .Name }} {{ if .Extends }}extends {{Join .Extends}} {{ end }}{
  {{- range .Fields }}
  {{ .Name }}{{ if .Nullable }}?{{ end }}: {{ .Type }}; {{ .Comment }}
  {{- end }}
}
{{ end -}}
