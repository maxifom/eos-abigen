package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

var VERSION = "0.0.1"

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "eos-abigen-go",
	Short: "CLI for generating RPC Client and Tables structures to read contracts on EOS-like blockchains",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.DisableAutoGenTag = true
	rootCmd.SilenceUsage = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .eos-abigen-go.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".eos-abigen-go")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func GenerateDocs(writer io.Writer) error {
	return doc.GenMarkdown(rootCmd, writer)
}

func GenerateDocsDir(dir string) error {
	return doc.GenMarkdownTree(rootCmd, dir)
}