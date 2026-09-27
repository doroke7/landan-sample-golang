package down

import (
	"github.com/spf13/cobra"

	"landan-desktop-fyne/sample/docker_compose_provider/internal/logger"
)

var Command = &cobra.Command{
	Use:   "down SERVICE",
	Short: "一次性任務,沒有需要停止的程序",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info(args[0] + ": 一次性任務,沒有需要停止的程序")
	},
}
