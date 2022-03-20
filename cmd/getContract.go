package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/maxifom/eos-abigen/pkg/client"
	"github.com/spf13/cobra"
)

var getContractCmd = &cobra.Command{
	Use:   "get-contract [flags] [...contract_names]",
	Short: "Downloads contract ABI from specified RPC",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rpcURL, err := cmd.Flags().GetString("rpc_url")
		if err != nil {
			return err
		}

		outputFolder, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		rpcClient, err := client.NewRPCClient(rpcURL, &http.Client{Timeout: 30 * time.Second})
		if err != nil {
			return fmt.Errorf("failed to create rpc client: %w", err)
		}

		err = os.MkdirAll(outputFolder, os.ModePerm)
		if err != nil {
			if !os.IsExist(err) {
				return fmt.Errorf("failed to create output folder %s: %w", outputFolder, err)
			}
		}

		for _, a := range args {
			abi, err := rpcClient.GetABI(context.Background(), a)
			if err != nil {
				return fmt.Errorf("failed to get abi for %s: %w", a, err)
			}

			err = func() error {
				f, err := os.Create(filepath.Join(outputFolder, fmt.Sprintf("%s.json", a)))
				if err != nil {
					return err
				}
				defer f.Close()

				return json.NewEncoder(f).Encode(abi.ABI)
			}()
			if err != nil {
				return err
			}

		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getContractCmd)
	getContractCmd.Flags().StringP("rpc_url", "u", "https://eos.greymass.com", "RPC URL to download ABI file from")
	getContractCmd.Flags().StringP("output", "o", "contracts", "Folder to output contract ABI to")
}
