package cmd

import (
	"fmt"
	"os"

	"github.com/felisest/comproxy/internal/infrastructure/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		_ = config.LoadConfig("config")
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
