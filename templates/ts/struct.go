package ts

import "strings"

var StructTemplate = strings.TrimSpace(`export interface {{.Name}} {
    {{ range .Fields }}{{.Name}}: {{.Type}}
    {{end}}
}`)
