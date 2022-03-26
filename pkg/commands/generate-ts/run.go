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
	type StructForNestedArray struct {
		F Field
		I int
	}

	t, err := template.New("client").Funcs(map[string]any{
		"genStructForNestedArray": func(i int, f Field) StructForNestedArray {
			return StructForNestedArray{
				F: f,
				I: i,
			}
		},
		"generateTabs": func(count int) string {
			return strings.Repeat("\t", count)
		},
		"sub": func(i int, sub int) int {
			return i - sub
		},
		"generateIBrackets": func(count int) string {
			if count == 0 {
				return ""
			}

			var s strings.Builder
			for i := 0; i < count; i++ {
				fmt.Fprintf(&s, "[i%d]", i)
			}
			return s.String()
		},
	}).Parse(ts.ClientTemplate)
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

	t, err = t.New("nested").Parse(ts.NestedArrayTemplate)
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

type Struct struct {
	Name   string
	Fields []Field
}

type Method struct {
	MethodName string
	TableName  string
	ReturnName string
	Struct     Struct
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
	structMap := map[string]Struct{}
	for i, s := range ss {
		structMap[s.Name] = Struct{
			Name:   s.Name,
			Fields: s.Fields,
		}

		err = t.ExecuteTemplate(&structTypes, "struct", map[string]interface{}{
			"Name":   s.Name,
			"Fields": s.Fields,
			"IsLast": i == len(ss)-1,
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

	var methods []Method
	for _, t := range abi.Tables {
		methods = append(methods, Method{
			MethodName: strings.ReplaceAll(strcase.LowerCamelCase(t.Name), ".", ""),
			TableName:  t.Name,
			ReturnName: strcase.UpperCamelCase(t.Type) + "Rows",
			Struct:     structMap[strcase.UpperCamelCase(t.Type)],
		})
	}

	err = t.ExecuteTemplate(clientF, "client", map[string]interface{}{
		"Version": version,
		"Methods": methods,
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
	Type             string
	Func             string
	Method           string
	IntermediateType string
	FullType         string
}

var DefaultTSCaseFunc = strcase.UpperCamelCase

var LanguageMapping = map[string]map[string]LanguageFieldMapping{
	"bool": {
		"ts": {
			Type:             "boolean",
			IntermediateType: "number",
			Func:             "!!",
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
			Type:   "string",
			Method: "toString",
		},
	},
	"uint64": {
		"ts": {
			Type:   "string",
			Method: "toString",
		},
	},
	"int128": {
		"ts": {
			Type: "string",
		},
	},
	"uint128": {
		"ts": {
			Type: "string",
		},
	},
	"float32": {
		"ts": {
			Type:             "number",
			Func:             "Number.parseFloat",
			IntermediateType: "string",
		},
	},
	"float64": {
		"ts": {
			Type:             "number",
			Func:             "Number.parseFloat",
			IntermediateType: "string",
		},
	},
	"float128": {
		"ts": {
			Type: "string",
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
			Type:             "Symbol",
			Func:             "new types.Symbol",
			FullType:         "types.Symbol",
			IntermediateType: "string",
		},
	},
	"symbol_code": {
		"ts": {
			Type: "string",
		},
	},
	"asset": {
		"ts": {
			Type:             "Asset",
			Func:             "new types.Asset",
			IntermediateType: "string",
			FullType:         "types.Asset",
		},
	},
	"extended_asset": {
		"ts": {
			Type:             "ExtendedAsset",
			Func:             "new types.ExtendedAsset",
			IntermediateType: "ExtendedAssetType",
			FullType:         "types.ExtendedAsset",
		},
	},
}

type Field struct {
	Name                string
	Type                string
	FullType            string
	IntermediateType    string
	Func                string
	Method              string
	ArraysCount         int
	ArraysCountIterator []int
}

type S struct {
	Name   string
	Fields []Field
}

func (f Field) FormatNameValue(obj string) string {
	if f.ArraysCount > 0 {
		return "[]"
	}

	if obj == "" {
		var s strings.Builder
		if f.Func != "" {
			fmt.Fprintf(&s, "%s(%s)", f.Func, f.Name)
		} else {
			fmt.Fprintf(&s, "%s", f.Name)
		}

		if f.Method != "" {
			fmt.Fprintf(&s, ".%s()", f.Method)
		}

		return s.String()
	}

	var s strings.Builder
	if f.Func != "" {
		fmt.Fprintf(&s, "%s(%s.%s)", f.Func, obj, f.Name)
	} else {
		fmt.Fprintf(&s, "%s.%s", obj, f.Name)
	}

	if f.Method != "" {
		fmt.Fprintf(&s, ".%s()", f.Method)
	}

	return s.String()
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
			arraysCount := strings.Count(fieldType, "[]")
			var aci []int
			for i := 0; i < arraysCount; i++ {
				aci = append(aci, i)
			}
			listsSuffix := strings.Repeat("[]", arraysCount)
			fieldType = strings.ReplaceAll(fieldType, "[]", "")
			if realFieldType, ok := newTypesMap[fieldType]; ok {
				fieldType = realFieldType
			}

			if fMapping, ok := LanguageMapping[fieldType]["ts"]; ok {
				realType := fMapping.Type + listsSuffix
				intermediateType := ""
				if fMapping.IntermediateType != "" {
					intermediateType = fMapping.IntermediateType + listsSuffix
				}

				s.Fields = append(s.Fields, Field{
					Name:                fieldName,
					Type:                realType,
					IntermediateType:    intermediateType,
					Func:                fMapping.Func,
					Method:              fMapping.Method,
					ArraysCount:         arraysCount,
					ArraysCountIterator: aci,
					FullType:            fMapping.FullType + listsSuffix,
				})
			} else {
				if structName, ok := newStructsMap[fieldType]; ok {
					realType := structName + listsSuffix

					if fmapping, ok := LanguageMapping[structName]["ts"]; ok {
						intermediateType := ""
						if fMapping.IntermediateType != "" {
							intermediateType = fMapping.IntermediateType + listsSuffix
						}
						s.Fields = append(s.Fields, Field{
							Name:                fieldName,
							Type:                realType,
							IntermediateType:    intermediateType,
							Func:                fmapping.Func,
							Method:              fMapping.Method,
							ArraysCount:         arraysCount,
							ArraysCountIterator: aci,
							FullType:            fMapping.FullType + listsSuffix,
						})
					} else {
						s.Fields = append(s.Fields, Field{
							Name:                fieldName,
							Type:                realType,
							ArraysCount:         arraysCount,
							ArraysCountIterator: aci,
						})
					}

				} else {
					s.Fields = append(s.Fields, Field{
						Name:                fieldName,
						Type:                "unknown" + listsSuffix,
						ArraysCount:         arraysCount,
						ArraysCountIterator: aci,
					})
				}
			}
		}

		ss = append(ss, s)
	}

	return ss
}
