package ts

var StructTemplate = `
export type {{.Name}} = {
	{{$lenFields := len .Fields -}}
	{{ range $i, $f := .Fields -}}{{$f.Name}}: {{$f.Type}}
{{if (lt $i (sub $lenFields 1))}}	{{end}}{{end -}}
};

export type {{.Name}}Interm = {
	{{ range $i, $f := .Fields -}}{{$f.Name}}: {{$f.Type}}{{ if ne $f.IntermediateType "" }} | {{$f.IntermediateType}}{{end}}
{{if (lt $i (sub $lenFields 1))}}	{{end}}{{end -}}
};

export type {{.Name}}Raw = {
	{{ range $i, $f := .Fields -}}{{$f.Name}}: {{ if ne $f.RawType ""}}{{ $f.RawType }}{{ else }}{{ if ne $f.IntermediateType "" }}{{ $f.IntermediateType }}{{ else }}{{ $f.Type }}{{ end }}{{ end }}
{{if (lt $i (sub $lenFields 1))}}	{{end}}{{end -}}
};
`
