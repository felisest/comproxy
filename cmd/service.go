package cmd

import (
	"fmt"

	"github.com/felisest/comproxy/internal/infrastructure/di"

	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Start proxy",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Service starting with args %s\n", args)
		return di.InitService()
	},
}
