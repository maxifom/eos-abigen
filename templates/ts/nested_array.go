package ts

import "strings"

var NestedArrayTemplate = strings.TrimSpace(`
	{{ define "nested" }}

	row.{{.F.Name}} = r.{{.F.Name}}.map(a => {{ "nested" . }})


`)

// {{ generateTabs (sub .I -4)}}let arr{{.I}}{{ if ne .ArrayType "" }}: {{.ArrayType}} {{ end }} = [];
// {{ generateTabs (sub .I -4)}}for (let i{{.I}} = 0; i{{.I}} < r.{{.F.Name}}{{generateIBrackets .I}}.length; i{{.I}}++) {
//			{{- if eq .I (sub .F.ArraysCount 1) }}
// {{ generateTabs (sub .I -5)}}arr{{.I}}.push({{ generatePush .F .I "r" }})
//			{{- else -}}
//			{{- template "nested" (genStructForNestedArray (sub .I -1) .F) -}}
//			{{- end }}
// {{ generateTabs (sub .I -4) -}} }
//		{{ if gt .I 0 -}}
// {{ generateTabs (sub .I -2) -}}arr{{sub .I 1}}.push(arr{{.I}});
//		{{- end -}}
//	{{end -}}
