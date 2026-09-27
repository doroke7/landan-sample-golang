package http

import (
	"github.com/spf13/cobra"

	"landan-desktop-fyne/sample/launcher/internal/router"
)

var Command = &cobra.Command{
	Use:   "http",
	Short: "啟動 HTTP 服務",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		sAddr, _ := cmd.Flags().GetString("addr")
		return router.New().Run(sAddr)
	},
}

func init() {
	Command.Flags().String("addr", ":8080", "監聽位址")
}
