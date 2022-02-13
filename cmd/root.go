package cmd

import (
	"os"

	"github.com/mehdy/sabet/pkg/sabet"
	"github.com/spf13/cobra"
)

var (
	cfgDir []string

	rootCmd = &cobra.Command{
		Use:   "sabet",
		Short: "A static tool for running jobs in an event-driven environment",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			sabet.NewManager(args).Run()
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
