package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"landan-desktop-fyne/sample/docker_compose_provider/cmd/desktop"
)

var oRootCommand = &cobra.Command{
	Use:           "desktop",
	Short:         "用宿主機的 ffmpeg 截圖,由 docker compose 驅動",
	SilenceUsage:  true, // 失敗時 docker compose 只需要 JSON 訊息,不要 usage
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	oRootCommand.AddCommand(desktop.Command)
}

// Execute 供 main.go 調用
func Execute() {
	if err := oRootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
