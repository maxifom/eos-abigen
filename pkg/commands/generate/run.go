package generate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/dave/jennifer/jen"
	"github.com/maxifom/eos-abigen-go/pkg/abitypes"
	"github.com/stoewer/go-strcase"
)

type Opts struct {
	ContractFilePath     string
	ContractNameOverride string
	GeneratedFolder      string
	Version              string
}

func Run(opts Opts) error {
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

	err = generateStructs(abi, contractName, opts.GeneratedFolder, opts.Version)
	if err != nil {
		return fmt.Errorf("failed to generate structs: %w", err)
	}

	realContractName := contractName
	if opts.ContractNameOverride != "" {
		realContractName = opts.ContractNameOverride
	}

	err = generateTables(abi, contractName, realContractName, opts.GeneratedFolder, opts.Version)
	if err != nil {
		return fmt.Errorf("failed to generate tables: %w", err)
	}

	return nil
}

func generateTables(abi abitypes.ABI, contractName string, realContractName string, generatedFolder string, version string) error {
	f := NewFile(contractName)
	f.HeaderComment(fmt.Sprintf("Generated by eos-abigen-go version %s", version))

	f.ImportAlias("github.com/maxifom/eos-abigen-go/pkg/client", "rpcClient")
	f.Const().Defs(
		Id("contractName").Op("=").Lit(realContractName),
	)
	f.Add(Type().Id("Client").Struct(
		Op("*").Qual("github.com/maxifom/eos-abigen-go/pkg/client", "RPCClient"),
	))

	constructor := Func().Id("NewClient").Params(Id("rpcClient").Op("*").Qual("github.com/maxifom/eos-abigen-go/pkg/client", "RPCClient")).Params(Id("*" + ("Client"))).Block(
		Return(Op("&").Id("Client").Values(
			Dict{
				Id("RPCClient"): Id("rpcClient"),
			},
		)),
	)
	f.Add(constructor)

	declaredRowsNames := map[string]struct{}{}

	for _, table := range abi.Tables {
		returnType := strcase.UpperCamelCase(table.Type)
		rowsName := returnType + "Rows"
		if _, ok := declaredRowsNames[rowsName]; !ok {
			f.Type().Id(rowsName).Struct(
				Qual("github.com/maxifom/eos-abigen-go/pkg/base", "BaseRows"),
				Id("Rows").Id("[]"+returnType).Tag(map[string]string{"json": "rows"}),
			)

			declaredRowsNames[rowsName] = struct{}{}
		}

		tableActionName := strings.ReplaceAll(strcase.UpperCamelCase(table.Name), ".", "")
		body := Empty()
		body.Id("result").Op(":=").Op("&").Id(rowsName).Values().Line().Line()

		body.Id("options").Op("=").Append(
			Index().Qual("github.com/maxifom/eos-abigen-go/pkg/client", "RequestOption").Values(
				Qual("github.com/maxifom/eos-abigen-go/pkg/client", "Table").Call(Lit(table.Name)),
				Qual("github.com/maxifom/eos-abigen-go/pkg/client", "Scope").Call(Id("contractName")),
				Qual("github.com/maxifom/eos-abigen-go/pkg/client", "Code").Call(Id("contractName")),
				Qual("github.com/maxifom/eos-abigen-go/pkg/client", "Limit").Call(Lit(-1)),
			),

			Id("options").Op("...")).Line()

		body.Err().Op(":=").Id("c").Dot("GetTableRows").Call(Id("ctx"), Id("result"), Id("options").Op("...")).Line()
		body.If(Err().Op("!=").Nil()).Block(
			Return(Nil(), Err()),
		).Line()
		returnStatement := Return(Id("result"), Nil())
		body.Add(returnStatement)
		id := Func().
			Params(Id("c").Id("*"+("Client"))).
			Id(tableActionName).
			Params(Id("ctx").Qual("context", "Context"), Id("options").Op("...").Qual("github.com/maxifom/eos-abigen-go/pkg/client", "RequestOption")).
			Params(Id("*"+rowsName), Id("error")).Block(body)
		f.Add(id)
	}

	fPath := filepath.Join(generatedFolder, contractName, "client.go")
	file, err := os.OpenFile(fPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	err = f.Render(file)
	if err != nil {
		return err
	}
	return nil

}

func generateStructs(abi abitypes.ABI, contractName string, generatedFolder string, version string) error {
	f := NewFile(contractName)
	f.HeaderComment(fmt.Sprintf("Generated by eos-abigen-go version %s", version))

	newTypesMap := map[string]string{}
	for _, t := range abi.Types {
		newTypesMap[t.NewTypeName] = t.Type
	}

	newStructsMap := map[string]string{}

	for _, abiStruct := range abi.Structs {
		newStructsMap[abiStruct.Name] = strcase.UpperCamelCase(abiStruct.Name)
	}

	for _, abiStruct := range abi.Structs {
		s := Type().Id(strcase.UpperCamelCase(abiStruct.Name))
		fields := make([]Code, 0, len(abiStruct.Fields))
		for _, field := range abiStruct.Fields {
			fieldName := strcase.UpperCamelCase(field.Name)
			fieldGen := Id(fieldName)
			fieldType := field.Type

			isList := false
			if strings.Contains(fieldType, "[]") {
				isList = true
				fieldType = strings.ReplaceAll(fieldType, "[]", "")
			}
			if realFieldType, ok := newTypesMap[fieldType]; ok {
				fieldType = realFieldType
			}
			switch fieldType {
			case "bool":
				if isList {
					fieldGen.List(Bool())
				}
				fieldGen.Op("*").Qual("github.com/maxifom/eos-abigen-go/pkg/base", "Bool")
			case "int8":
				if isList {
					fieldGen.Index()
				}
				fieldGen.Int8()
			case "uint8":
				if isList {
					fieldGen.Index()
				}
				fieldGen.Uint8()
			case "int16":
				if isList {
					fieldGen.Index()
				}
				fieldGen.Int16()
			case "uint16":
				if isList {
					fieldGen.Index()
				}
				fieldGen.Uint16()
			case "int32", "varint32":
				if isList {
					fieldGen.Index()
				}
				fieldGen.Int32()
			case "uint32", "varuint32":
				if isList {
					fieldGen.Index()
				}
				fieldGen.Uint32()
			case "int64":
				if isList {
					fieldGen.Index()
				}
				fieldGen.Int64()

			case "uint64":
				if isList {
					fieldGen.Index()
				}
				fieldGen.Op("*").Qual("github.com/maxifom/eos-abigen-go/pkg/base", "UInt64")
			case "int128", "uint128":
				if isList {
					fieldGen.Index()
				}
				fieldGen.Op("*").Qual("github.com/maxifom/eos-abigen-go/pkg/base", "BigInt")
			case "float32":
				if isList {
					fieldGen.Index()
				}
				fieldGen.Op("*").Qual("github.com/maxifom/eos-abigen-go/pkg/base", "Float32")
			case "float64":
				if isList {
					fieldGen.Index()
				}
				fieldGen.Float64()
			case "float128":
				if isList {
					fieldGen.Index()
				}
				fieldGen.Float64()
			case "time_point":
				if isList {
					fieldGen.Index()
				}
				fieldGen.String()
			case "time_point_sec":
				if isList {
					fieldGen.Index()
				}
				fieldGen.String()
			case "block_timestamp_type":
				if isList {
					fieldGen.Index()
				}
				fieldGen.String()
			case "name":
				if isList {
					fieldGen.Index()
				}
				fieldGen.String()
			case "bytes":
				if isList {
					fieldGen.Index()
				}
				fieldGen.String()
			case "string":
				if isList {
					fieldGen.Index()
				}
				fieldGen.String()
			case "checksum160":
				if isList {
					fieldGen.Index()
				}
				fieldGen.String()
			case "checksum256":
				if isList {
					fieldGen.Index()
				}
				fieldGen.String()
			case "checksum512":
				if isList {
					fieldGen.Index()
				}
				fieldGen.String()
			case "public_key":
				if isList {
					fieldGen.Index()
				}
				fieldGen.String()
			case "signature":
				if isList {
					fieldGen.Index()
				}
				fieldGen.String()
			case "symbol":
				if isList {
					fieldGen.Index()
				}
				fieldGen.String()
			case "symbol_code":
				if isList {
					fieldGen.Index()
				}
				fieldGen.String()
			case "asset":
				if isList {
					fieldGen.Index()
				}
				fieldGen.String()
			case "extended_asset":
				if isList {
					fieldGen.Index()
				}
				fieldGen.Qual("github.com/maxifom/eos-abigen-go/pkg/base", "ExtendedAsset")
			default:
				if structName, ok := newStructsMap[fieldType]; ok {
					if isList {
						fieldGen.Index()
					}
					fieldGen.Id(structName)
				} else {
					fieldGen.Qual("encoding/json", "RawMessage")
				}
			}

			fieldGen.Tag(map[string]string{"json": field.Name})
			fields = append(fields, fieldGen)
		}

		s.Struct(fields...)
		f.Add(s)
	}

	fPath := filepath.Join(generatedFolder, contractName, "types.go")
	file, err := os.OpenFile(fPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	err = f.Render(file)
	if err != nil {
		return err
	}
	return nil
}
