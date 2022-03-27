# eos-abigen

CLI for generating type-safe clients for EOS-compatible contracts.

## Features

**Golang**:

* Struct generation and Client generation for read-access get_table_rows with sensible defaults
* Utility classes for Asset, ExtendedAsset, Symbol, time structs
* Support for nested arrays and structs

**Typescript**:

* Struct generation and client generation for read-access using get_table_rows with sensible defaults
* Action builder
* Utility classes for Asset, ExtendedAsset, Symbol, time structs
* Support for nested arrays and structs

## Installing

MacOS Intel:

```shell
wget https://github.com/maxifom/eos-abigen/releases/latest/download/eos-abigen_macos_amd64
mv eos-abigen_macos_amd64 eos-abigen
chmod +x eos-abigen
```

MacOS M1:

```shell
wget https://github.com/maxifom/eos-abigen/releases/latest/download/eos-abigen_macos_arm64
mv eos-abigen_macos_arm64 eos-abigen
chmod +x eos-abigen
```

Linux:

```shell
wget https://github.com/maxifom/eos-abigen/releases/latest/download/eos-abigen_linux_amd64
mv eos-abigen_linux_amd64 eos-abigen
chmod +x eos-abigen
```

Windows:

```shell
Download https://github.com/maxifom/eos-abigen/releases/latest/download/eos-abigen_windows_amd64.exe
```

Or install using go install:

```shell
go install github.com/maxifom/eos-abigen@latest
```

## Getting started

To get started you need a ABI JSON file for the contract. You can download it
using `eos-abigen get-contract <CONTRACT_NAME>` command. Note that default is `https://eos.greymass.com`, if you need
another chain specify RPC Node URL using `-u` param.

After you have downloaded contract, generate code using `eos-abigen generate <ABI_JSON_PATH>` for Golang
and `eos-abigen generate-ts <ABI_JSON_PATH>` for Typescript.

Also you can provide `.eos-abigen.yaml` config file and invoke generation using `eos-abigen generate`
or `eos-abigen generate-ts` without arguments

Example .eos-abigen.yaml:

```yaml
generate:
  folder: generated # Folder to which save generated code from current dir
  contracts:
    - file: ./contracts/eosio.json # Path to ABI JSON file
      name_override: eosio123 # Contract name override, it will be used for get_table_rows and Action builder
    - file: /some/folder/contracts/eosio.json
      name_override: eosio123
```
