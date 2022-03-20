package generate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/maxifom/eos-abigen/pkg/abitypes"
	"github.com/maxifom/eos-abigen/templates/ts"
	"github.com/stoewer/go-strcase"
)

type Opts struct {
	ContractFilePath     string
	ContractNameOverride string
	GeneratedFolder      string
	Version              string
}

func Run(opts Opts) error {
	t, err := template.New("client").Parse(ts.ClientTemplate)
	if err != nil {
		return err
	}
	t, err = t.New("index").Parse(ts.IndexTemplate)
	if err != nil {
		return err
	}
	t, err = t.New("struct").Parse(ts.StructTemplate)
	if err != nil {
		return err
	}
	t, err = t.New("table_rows").Parse(ts.TableRowsTemplate)
	if err != nil {
		return err
	}
	t, err = t.New("types").Parse(ts.TypesTemplate)
	if err != nil {
		return err
	}

	contractName := strings.TrimSuffix(filepath.Base(opts.ContractFilePath), filepath.Ext(opts.ContractFilePath))

	abiF, err := os.ReadFile(opts.ContractFilePath)
	if err != nil {
		return fmt.Errorf("failed to read contract file %s: %w", opts.ContractFilePath, err)
	}

	var abi abitypes.ABI

	err = json.Unmarshal(abiF, &abi)
	if err != nil {
		return fmt.Errorf("failed to unmarshal contract JSON: %w", err)
	}

	fPath := filepath.Join(opts.GeneratedFolder, contractName)
	err = os.MkdirAll(fPath, os.ModePerm)
	if err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("failed to create output folder %s: %w", fPath, err)
		}
	}

	realContractName := contractName
	if opts.ContractNameOverride != "" {
		realContractName = opts.ContractNameOverride
	}

	return gen(abi, contractName, realContractName, opts.GeneratedFolder, opts.Version, t)
}

type Method struct {
	MethodName string
	TableName  string
	ReturnName string
}

func gen(abi abitypes.ABI, contractName string, realContractName string, generatedFolder string, version string, t *template.Template) error {
	typesF, err := os.Create(filepath.Join(generatedFolder, contractName, "types.ts"))
	if err != nil {
		return err
	}
	defer typesF.Close()

	indexF, err := os.Create(filepath.Join(generatedFolder, contractName, "index.ts"))
	if err != nil {
		return err
	}
	defer indexF.Close()

	clientF, err := os.Create(filepath.Join(generatedFolder, contractName, "client.ts"))
	if err != nil {
		return err
	}
	defer clientF.Close()

	err = t.ExecuteTemplate(indexF, "index", map[string]interface{}{
		"Version": version,
	})
	if err != nil {
		return err
	}

	var methods []Method
	for _, t := range abi.Tables {
		methods = append(methods, Method{
			MethodName: strings.ReplaceAll(strcase.LowerCamelCase(t.Name), ".", ""),
			TableName:  t.Name,
			ReturnName: strcase.UpperCamelCase(t.Type) + "Rows",
		})
	}

	err = t.ExecuteTemplate(clientF, "client", map[string]interface{}{
		"Version": version,
		"Methods": methods,
	})
	if err != nil {
		return err
	}

	var rowTypes strings.Builder
	var structTypes strings.Builder
	declaredRowsNames := map[string]struct{}{}
	for _, table := range abi.Tables {
		tableName := strcase.UpperCamelCase(table.Type)
		if _, ok := declaredRowsNames[tableName]; !ok {
			err = t.ExecuteTemplate(&rowTypes, "table_rows", map[string]interface{}{
				"TableName": tableName,
			})
			if err != nil {
				return err
			}

			declaredRowsNames[tableName] = struct{}{}
		}
	}

	ss := genStructs(abi, getNewTypesMap(abi), getNewStructsMap(abi))
	for _, s := range ss {
		err = t.ExecuteTemplate(&structTypes, "struct", map[string]interface{}{
			"Name":   s.Name,
			"Fields": s.Fields,
		})
	}

	err = t.ExecuteTemplate(typesF, "types", map[string]interface{}{
		"Version":      version,
		"ContractName": realContractName,
		"RowTypes":     rowTypes.String(),
		"StructTypes":  structTypes.String(),
	})
	if err != nil {
		return err
	}

	return nil
}

func getNewTypesMap(abi abitypes.ABI) map[string]string {
	newTypesMap := map[string]string{}
	for _, t := range abi.Types {
		newTypesMap[t.NewTypeName] = t.Type
	}

	return newTypesMap
}

func getNewStructsMap(abi abitypes.ABI) map[string]string {
	newStructsMap := map[string]string{}
	for _, abiStruct := range abi.Structs {
		newStructsMap[abiStruct.Name] = DefaultTSCaseFunc(abiStruct.Name)
	}

	return newStructsMap
}

type LanguageFieldMapping struct {
	Type string
}

var DefaultTSCaseFunc = strcase.UpperCamelCase

var LanguageMapping = map[string]map[string]LanguageFieldMapping{
	"bool": {
		"ts": {
			Type: "boolean",
		},
	},
	"int8": {
		"ts": {
			Type: "number",
		},
	},
	"uint8": {
		"ts": {
			Type: "number",
		},
	},
	"int16": {
		"ts": {
			Type: "number",
		},
	},
	"uint16": {
		"ts": {
			Type: "number",
		},
	},
	"int32": {
		"ts": {
			Type: "number",
		},
	},
	"varint32": {
		"ts": {
			Type: "number",
		},
	},
	"uint32": {
		"ts": {
			Type: "number",
		},
	},
	"varuint32": {
		"ts": {
			Type: "number",
		},
	},
	"int64": {
		"ts": {
			Type: "number",
		},
	},
	"uint64": {
		"ts": {
			Type: "number",
		},
	},
	"int128": {
		"ts": {
			Type: "number",
		},
	},
	"uint128": {
		"ts": {
			Type: "number",
		},
	},
	"float32": {
		"ts": {
			Type: "number",
		},
	},
	"float64": {
		"ts": {
			Type: "number",
		},
	},
	"float128": {
		"ts": {
			Type: "number",
		},
	},
	"time_point": {
		"ts": {
			Type: "string",
		},
	},
	"time_point_sec": {
		"ts": {
			Type: "string",
		},
	},
	"block_timestamp_type": {
		"ts": {
			Type: "string",
		},
	},
	"name": {
		"ts": {
			Type: "string",
		},
	},
	"bytes": {
		"ts": {
			Type: "string",
		},
	},
	"string": {
		"ts": {
			Type: "string",
		},
	},
	"checksum160": {
		"ts": {
			Type: "string",
		},
	},
	"checksum256": {
		"ts": {
			Type: "string",
		},
	},
	"checksum512": {
		"ts": {
			Type: "string",
		},
	},
	"public_key": {
		"ts": {
			Type: "string",
		},
	},
	"signature": {
		"ts": {
			Type: "string",
		},
	},
	"symbol": {
		"ts": {
			Type: "string",
		},
	},
	"symbol_code": {
		"ts": {
			Type: "string",
		},
	},
	"asset": {
		"ts": {
			Type: "string",
		},
	},
	"extended_asset": {
		"ts": {
			Type: "string",
		},
	},
}

type Field struct {
	Name string
	Type string
}

type S struct {
	Name   string
	Fields []Field
}

func genStructs(abi abitypes.ABI, newTypesMap map[string]string, newStructsMap map[string]string) []S {
	ss := make([]S, 0, len(abi.Structs))
	for _, abiStruct := range abi.Structs {
		s := S{
			Name: DefaultTSCaseFunc(abiStruct.Name),
		}

		for _, field := range abiStruct.Fields {
			fieldName := strings.ToLower(field.Name)
			fieldType := field.Type
			listsSuffix := strings.Repeat("[]", strings.Count(fieldType, "[]"))
			fieldType = strings.ReplaceAll(fieldType, "[]", "")
			if realFieldType, ok := newTypesMap[fieldType]; ok {
				fieldType = realFieldType
			}

			if fMapping, ok := LanguageMapping[fieldType]["ts"]; ok {
				realType := fMapping.Type + listsSuffix
				s.Fields = append(s.Fields, Field{
					Name: fieldName,
					Type: realType,
				})
			} else {
				if structName, ok := newStructsMap[fieldType]; ok {
					realType := structName + listsSuffix
					s.Fields = append(s.Fields, Field{
						Name: fieldName,
						Type: realType,
					})
				} else {
					s.Fields = append(s.Fields, Field{
						Name: fieldName,
						Type: "unknown",
					})
				}
			}
		}

		ss = append(ss, s)
	}

	return ss
}
