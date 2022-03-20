package ts

import "strings"

var TableRowsTemplate = strings.TrimSpace(`export interface {{.TableName}}Rows {
    more: boolean;
    next_key: string;
    rows: {{.TableName}}[];
}`)
