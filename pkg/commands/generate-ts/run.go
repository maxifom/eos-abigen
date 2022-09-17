package generate

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/maxifom/eos-abigen/pkg/abitypes"
	"github.com/spf13/afero"

	"github.com/stoewer/go-strcase"
)

type Opts struct {
	ContractFilePath     string
	ContractNameOverride string
	GeneratedFolder      string
	Version              string
	FS                   afero.Fs
}

func generateIBrackets(count int) string {
	if count == 0 {
		return ""
	}

	var s strings.Builder
	for i := 0; i < count; i++ {
		fmt.Fprintf(&s, "[i%d]", i)
	}
	return s.String()
}

//go:embed templates/ts/*.gotmpl
var templatesFs embed.FS

func Run(opts Opts) error {
	fs := opts.FS
	if fs == nil {
		return fmt.Errorf("no fs specified")
	}

	type StructForNestedArray struct {
		F         Field
		I         int
		ArrayType string
	}

	type StructForFieldMapper struct {
		F            Field
		IsLast       bool
		UseFullTypes bool
		Tabs         int
	}

	clientTemplate, err := templatesFs.ReadFile("templates/ts/client.gotmpl")
	if err != nil {
		return err
	}

	t, err := template.New("client").Funcs(sprig.TxtFuncMap()).Funcs(map[string]any{
		"genStructForFieldMapper": func(f Field, isLast bool, useFullTypes bool, tabs int) StructForFieldMapper {
			return StructForFieldMapper{
				F:            f,
				IsLast:       isLast,
				UseFullTypes: useFullTypes,
				Tabs:         tabs,
			}
		},
		"genStructForNestedArray": func(i int, f Field) StructForNestedArray {
			array := StructForNestedArray{
				F: f,
				I: i,
			}

			if f.FullType != "" {
				array.ArrayType = f.FullType
			} else {
				array.ArrayType = f.Type
			}

			for j := 0; j < i; j++ {
				array.ArrayType = string([]rune(array.ArrayType)[:len([]rune(array.ArrayType))-2])
			}

			return array
		},
		"generateTabs": func(count int) string {
			return strings.Repeat("\t", count)
		},
		"sub": func(i int, sub int) int {
			return i - sub
		},
		"generateIBrackets": generateIBrackets,
		"generatePush": func(f Field, i int, obj string) string {
			var s string
			if f.Func != "" {
				s = fmt.Sprintf("%s(%s.%s%s)", f.Func, obj, f.Name, generateIBrackets(i+1))
			} else {
				s = fmt.Sprintf("%s.%s%s", obj, f.Name, generateIBrackets(i+1))
			}

			if f.Method != "" {
				s = fmt.Sprintf("%s.%s()", s, f.Method)
			}

			return s
		},
	}).Parse(string(clientTemplate))
	if err != nil {
		return err
	}

	indexTemplate, err := templatesFs.ReadFile("templates/ts/index.gotmpl")
	if err != nil {
		return err
	}

	t, err = t.New("index").Parse(string(indexTemplate))
	if err != nil {
		return err
	}

	structTemplate, err := templatesFs.ReadFile("templates/ts/struct.gotmpl")
	if err != nil {
		return err
	}

	t, err = t.New("struct").Parse(string(structTemplate))
	if err != nil {
		return err
	}

	tableRowsTemplate, err := templatesFs.ReadFile("templates/ts/table_rows.gotmpl")
	if err != nil {
		return err
	}

	t, err = t.New("table_rows").Parse(string(tableRowsTemplate))
	if err != nil {
		return err
	}

	typesTemplate, err := templatesFs.ReadFile("templates/ts/types.gotmpl")
	if err != nil {
		return err
	}

	t, err = t.New("types").Parse(string(typesTemplate))
	if err != nil {
		return err
	}

	actionBuilderTemplate, err := templatesFs.ReadFile("templates/ts/action_builder.gotmpl")
	if err != nil {
		return err
	}
	t, err = t.New("action_builder").Parse(string(actionBuilderTemplate))
	if err != nil {
		return err
	}

	fieldMapperTemplate, err := templatesFs.ReadFile("templates/ts/field_mapper.gotmpl")
	if err != nil {
		return err
	}
	t, err = t.New("map_field").Parse(string(fieldMapperTemplate))
	if err != nil {
		return err
	}

	contractName := strings.TrimSuffix(filepath.Base(opts.ContractFilePath), filepath.Ext(opts.ContractFilePath))

	abiF, err := afero.ReadFile(fs, opts.ContractFilePath)
	if err != nil {
		return fmt.Errorf("failed to read contract file %s: %w", opts.ContractFilePath, err)
	}

	var abi abitypes.ABI

	err = json.Unmarshal(abiF, &abi)
	if err != nil {
		return fmt.Errorf("failed to unmarshal contract JSON: %w", err)
	}

	fPath := filepath.Join(opts.GeneratedFolder, contractName)
	err = fs.MkdirAll(fPath, os.ModePerm)
	if err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("failed to create output folder %s: %w", fPath, err)
		}
	}

	realContractName := contractName
	if opts.ContractNameOverride != "" {
		realContractName = opts.ContractNameOverride
	}

	return gen(abi, contractName, realContractName, opts.GeneratedFolder, opts.Version, t, fs)
}

type Struct struct {
	Name   string
	Fields []Field
}

type Method struct {
	MethodName     string
	TableName      string
	ReturnName     string
	Struct         Struct
	ReturnNameRows string
}

type Action struct {
	Name       string
	ParamsName string
}

func gen(abi abitypes.ABI, contractName string, realContractName string, generatedFolder string, version string, t *template.Template, fs afero.Fs) error {
	typesF, err := fs.Create(filepath.Join(generatedFolder, contractName, "types.ts"))
	if err != nil {
		return err
	}
	defer typesF.Close()

	indexF, err := fs.Create(filepath.Join(generatedFolder, contractName, "index.ts"))
	if err != nil {
		return err
	}
	defer indexF.Close()

	clientF, err := fs.Create(filepath.Join(generatedFolder, contractName, "client.ts"))
	if err != nil {
		return err
	}
	defer clientF.Close()

	actionBuilderF, err := fs.Create(filepath.Join(generatedFolder, contractName, "action_builder.ts"))
	if err != nil {
		return err
	}
	defer actionBuilderF.Close()

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
	for _, s := range ss {
		structMap[s.Name] = Struct{
			Name:   s.Name,
			Fields: s.Fields,
		}

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

	var methods []Method
	for _, t := range abi.Tables {
		methods = append(methods, Method{
			MethodName:     strings.ReplaceAll(strcase.LowerCamelCase(t.Name), ".", ""),
			TableName:      t.Name,
			ReturnNameRows: strcase.UpperCamelCase(t.Type) + "Rows",
			ReturnName:     strcase.UpperCamelCase(t.Type),
			Struct:         structMap[strcase.UpperCamelCase(t.Type)],
		})
	}

	err = t.ExecuteTemplate(clientF, "client", map[string]interface{}{
		"Version": version,
		"Methods": methods,
	})
	if err != nil {
		return err
	}

	var actions []Action
	for _, action := range abi.Actions {
		actions = append(actions, Action{
			Name:       action.Name,
			ParamsName: strcase.UpperCamelCase(action.Type),
		})
	}

	err = t.ExecuteTemplate(actionBuilderF, "action_builder", map[string]interface{}{
		"Version": version,
		"Actions": actions,
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
	Type                 string
	Func                 string
	Method               string
	IntermediateType     string
	IntermediateFullType string
	FullType             string
	RawType              string
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
			Type:             "string",
			IntermediateType: "number",
			Method:           "toString",
			RawType:          "string | number",
		},
	},
	"uint64": {
		"ts": {
			Type:             "string",
			IntermediateType: "number",
			Method:           "toString",
			RawType:          "string | number",
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
			Type:                 "ExtendedAsset",
			Func:                 "new types.ExtendedAsset",
			IntermediateType:     "ExtendedAssetType",
			FullType:             "types.ExtendedAsset",
			IntermediateFullType: "types.ExtendedAssetType",
		},
	},
}

type Field struct {
	Name                 string
	Type                 string
	FullType             string
	IntermediateType     string
	IntermediateFullType string
	RawType              string
	Func                 string
	Method               string
	ArraysCount          int
	IsStruct             bool
	GenerateMapper       bool
}

type S struct {
	Name   string
	Fields []Field
}

func (f Field) FormatArrayValue(n string, useFullType bool) string {
	var s strings.Builder
	if f.Func != "" {
		name := n
		fmt.Fprintf(&s, "%s(%s)", f.Func, name)
	} else {
		fmt.Fprintf(&s, "%s", n)
	}

	if f.Method != "" {
		fmt.Fprintf(&s, ".%s()", f.Method)
	}

	ss := s.String()
	if !useFullType {
		ss = strings.ReplaceAll(ss, "types.", "")
	}

	return ss

}

func (f Field) FormatNameValue(obj string, useFullType bool) string {
	if f.GenerateMapper && f.ArraysCount > 0 {
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

	ss := s.String()
	if !useFullType {
		ss = strings.ReplaceAll(ss, "types.", "")
	}

	return ss
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
				intermediateFullType := ""
				fullType := ""
				rawType := ""
				if fMapping.IntermediateType != "" {
					intermediateType = fMapping.IntermediateType + listsSuffix
				}
				if fMapping.IntermediateFullType != "" {
					intermediateFullType = fMapping.IntermediateFullType + listsSuffix
				}
				if fMapping.FullType != "" {
					fullType = fMapping.FullType + listsSuffix
				}

				if fMapping.RawType != "" {
					splitted := strings.Split(fMapping.RawType, "|")
					var rawTypes []string
					for _, s := range splitted {
						rawTypes = append(rawTypes, strings.TrimSpace(s)+listsSuffix)
					}
					rawType = strings.Join(rawTypes, " | ")
				}

				s.Fields = append(s.Fields, Field{
					Name:                 fieldName,
					Type:                 realType,
					IntermediateType:     intermediateType,
					RawType:              rawType,
					IntermediateFullType: intermediateFullType,
					Func:                 fMapping.Func,
					Method:               fMapping.Method,
					ArraysCount:          arraysCount,
					FullType:             fullType,
				})
			} else {
				if structName, ok := newStructsMap[fieldType]; ok {
					realType := structName + listsSuffix

					if fmapping, ok := LanguageMapping[structName]["ts"]; ok {
						intermediateType := ""
						intermediateFullType := ""
						rawType := ""
						if fMapping.IntermediateType != "" {
							intermediateType = fMapping.IntermediateType + listsSuffix
						}
						if fMapping.IntermediateFullType != "" {
							intermediateFullType = fMapping.IntermediateFullType + listsSuffix
						}
						if fMapping.RawType != "" {
							splitted := strings.Split(fmapping.RawType, "|")
							var rawTypes []string
							for _, s := range splitted {
								rawTypes = append(rawTypes, strings.TrimSpace(s)+listsSuffix)
							}
							rawType = strings.Join(rawTypes, " | ")
						}

						s.Fields = append(s.Fields, Field{
							Name:                 fieldName,
							Type:                 realType,
							FullType:             "types." + realType,
							IntermediateType:     intermediateType,
							RawType:              rawType,
							IntermediateFullType: intermediateFullType,
							Func:                 fmapping.Func,
							Method:               fMapping.Method,
							ArraysCount:          arraysCount,
							IsStruct:             true,
						})
					} else {
						s.Fields = append(s.Fields, Field{
							Name:        fieldName,
							Type:        realType,
							FullType:    "types." + realType,
							ArraysCount: arraysCount,
							IsStruct:    true,
						})
					}

				} else {
					s.Fields = append(s.Fields, Field{
						Name:        fieldName,
						Type:        "unknown" + listsSuffix,
						ArraysCount: arraysCount,
					})
				}
			}
		}

		ss = append(ss, s)
	}

	for i := range ss {
		for j := range ss[i].Fields {
			if ss[i].Fields[j].IsStruct {
				ss[i].Fields[j].Func = fmt.Sprintf("types.map%s", strcase.UpperCamelCase(strings.ReplaceAll(ss[i].Fields[j].Type, "[]", "")))
				ss[i].Fields[j].IntermediateType = strcase.UpperCamelCase(strings.ReplaceAll(ss[i].Fields[j].Type, "[]", "")) + "Interm" + strings.Repeat("[]", ss[i].Fields[j].ArraysCount)
				if ss[i].Fields[j].IntermediateFullType != "" {
					ss[i].Fields[j].IntermediateFullType = strcase.UpperCamelCase(strings.ReplaceAll(ss[i].Fields[j].IntermediateFullType, "[]", "")) + "Interm" + strings.Repeat("[]", ss[i].Fields[j].ArraysCount)
				} else {
					ss[i].Fields[j].IntermediateFullType = "types." + strcase.UpperCamelCase(strings.ReplaceAll(ss[i].Fields[j].IntermediateType, "[]", "")) + strings.Repeat("[]", ss[i].Fields[j].ArraysCount)
				}
			}

			generateMapper := ss[i].Fields[j].ArraysCount > 0 && (ss[i].Fields[j].Func != "" || ss[i].Fields[j].Method != "" || ss[i].Fields[j].IntermediateType != "")

			ss[i].Fields[j].GenerateMapper = generateMapper
		}
	}

	return ss
}
