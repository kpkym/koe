package cmd

import (
	"github.com/kpkym/koe/cmd/web"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(web.Cmd)
}
