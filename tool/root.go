package tool

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go run ./cmd/tool/main.go",
	Short: "TODO: A brief description of your application",
	Long:  `TODO: A longer description that spans multiple lines and likely contains examples`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
