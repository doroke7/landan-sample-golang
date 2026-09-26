package up

import (
	"context"

	"github.com/spf13/cobra"

	"landan-desktop-fyne/sample/docker_compose_provider/internal/logger"
	"landan-desktop-fyne/sample/docker_compose_provider/internal/screenshoter"
)

var (
	sDevice string
	sOutput string
)

var Command = &cobra.Command{
	Use:   "up SERVICE",
	Short: "截一張圖",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sPath, err := screenshoter.Take(context.Background(), sDevice, sOutput)
		if err != nil {
			logger.Error(args[0] + ": " + err.Error())
			return err
		}
		logger.Info(args[0] + ": 截圖完成 " + sPath)
		return nil
	},
}

func init() {
	Command.Flags().StringVar(&sDevice, "device", "default", "avfoundation 影像裝置")
	Command.Flags().StringVar(&sOutput, "output", "./runtime/desktop/shot-{time}.jpg", "輸出檔案,{time} 代表時間戳")
}
