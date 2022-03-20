package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/maxifom/eos-abigen-go/pkg/commands/generate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ContractConfig struct {
	File         string `json:"file" yaml:"file" mapstructure:"file"`
	NameOverride string `json:"name_override" yaml:"name_override" mapstructure:"name_override"`
}

var generateCmd = &cobra.Command{
	Use:   "generate [flags] [abi_file]",
	Short: "Generate client and table structures from ABI contract file",
	Long:  "Generate client and table structures from ABI contract file. \nYou can also provide .eos-abigen-go.yaml file to generate multiple contracts with one command",
	Args: func(cmd *cobra.Command, args []string) error {
		var contracts []ContractConfig
		viper.UnmarshalKey("generate.contracts", &contracts)

		if len(contracts) > 0 {
			for _, c := range contracts {
				exists, err := fileExists(c.File)
				if err != nil {
					return err
				}
				if !exists {
					return fmt.Errorf("file %s does not exist", c.File)
				}
			}

			return nil
		}

		if len(args) != 1 {
			return fmt.Errorf("1 file is required. provided: %d", len(args))
		}

		exists, err := fileExists(args[0])
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("file %s does not exist. \nYou can use eos-abigen-go get-contract %s command to download it", args[0], args[0])
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		contractNameOverride, err := cmd.Flags().GetString("contract_name_override")
		if err != nil {
			return err
		}

		if len(args) > 0 {
			opts := generate.Opts{
				ContractFilePath:     args[0],
				ContractNameOverride: contractNameOverride,
				GeneratedFolder:      viper.GetString("generate.folder"),
				Version:              VERSION,
			}

			err = generate.Run(opts)
			if err != nil {
				return fmt.Errorf("failed to generate for %s: %w", args[0], err)
			}

			return nil
		}

		var contracts []ContractConfig
		err = viper.UnmarshalKey("generate.contracts", &contracts)
		if err != nil {
			return fmt.Errorf("failed to get contracts from config file: %w", err)
		}

		for _, c := range contracts {
			opts := generate.Opts{
				ContractFilePath:     c.File,
				ContractNameOverride: c.NameOverride,
				GeneratedFolder:      viper.GetString("generate.folder"),
				Version:              VERSION,
			}

			err = generate.Run(opts)
			if err != nil {
				return fmt.Errorf("failed to generate for %s: %w", c.File, err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringP("contract_name_override", "c", "", "contract name to use in calls to RPC. (default abi filename without extension)")
	generateCmd.Flags().StringP("folder", "f", "generated", "folder for generated files output")
	err := viper.BindPFlag("generate.folder", generateCmd.Flags().Lookup("folder"))
	if err != nil {
		panic(err)
	}
}

func fileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
