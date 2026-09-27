package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"landan-desktop-fyne/sample/launcher/cmd/http"
)

var rootCmd = &cobra.Command{
	Use:           "http",
	Short:         "http",
	SilenceUsage:  true,
	SilenceErrors: true, // Execute prints it once
}

func init() {
	rootCmd.AddCommand(http.Command)
}

// Execute runs the root command; called from main.go.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
