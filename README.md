# eos-abigen-go

CLI for generating RPC Client and Tables structures to read contracts on EOS-like blockchains

## Installing
```shell
go install github.com/maxifom/eos-abigen-go@latest
```


### Global Options

```
      --config string   config file (default is .eos-abigen-go.yaml)
  -h, --help            help for eos-abigen-go
```

## Generate command
Generate client and table structures from ABI contract file.
You can also provide .eos-abigen-go.yaml file to generate multiple contracts with one command

```
eos-abigen-go generate [flags] [abi_file]
```

### Options

```
  -c, --contract_name_override string   contract name to use in calls to RPC. (default abi filename without extension)
  -f, --folder string                   folder for generated files output (default "generated")
  -h, --help                            help for generate
```

### Options inherited from parent commands

```
      --config string   config file (default is .eos-abigen-go.yaml)
```

## Get contract command

Downloads contract ABI from specified RPC

```
eos-abigen-go get-contract [flags] [...contract_names]
```

### Options

```
  -h, --help             help for get-contract
  -o, --output string    Folder to output contract ABI to (default "contracts")
  -u, --rpc_url string   RPC URL to download ABI file from (default "https://eos.greymass.com")
```