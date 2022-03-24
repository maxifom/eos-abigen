package ts

import "strings"

var StructTemplate = strings.TrimSpace(`export type {{.Name}} = {
    {{ range .Fields }}{{.Name}}: {{.Type}}
    {{end}}
};
export type {{.Name}}Interm = {
    {{ range .Fields }}{{.Name}}: 	{{ if ne .IntermediateType "" }}{{.IntermediateType}}{{else}}{{.Type}}{{end}}
    {{end}}
};

`)
