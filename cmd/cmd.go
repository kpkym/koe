package cmd

import (
	"github.com/spf13/cobra"
	"kpk-koe/cmd/web"
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
