package ts

import "strings"

var TableRowsTemplate = strings.TrimSpace(`export type {{.TableName}}Rows = {
    more: boolean;
    next_key: string;
    rows: {{.TableName}}[];
};
export type {{.TableName}}RowsInterm = {
    more: boolean;
    next_key: string;
    rows: {{.TableName}}Interm[];
};
`)
