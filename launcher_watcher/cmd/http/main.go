package http

import (
	"github.com/spf13/cobra"

	"launcher_supervise/bootstrap"
	"launcher_supervise/internal/router"
)

var Command = &cobra.Command{
	Use:   "http",
	Short: "啟動 HTTP 服務",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		sAddr, _ := cmd.Flags().GetString("addr")

		// --watcher:自己不提供服務,改成執行不帶 --watcher 的自己,結束(崩潰或被關掉)就重啟。
		if bWatcher, _ := cmd.Flags().GetBool("watcher"); bWatcher {
			return bootstrap.WatchLauncher("http", "--addr", sAddr)
		}

		return router.New().Run(sAddr)
	},
}

func init() {
	Command.Flags().String("addr", ":8080", "監聽位址")
	Command.Flags().Bool("watcher", false, "由 supervisor 看守服務,結束(崩潰或被關掉)就自動重啟")
}
