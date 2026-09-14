package cmd

import (
	"context"
	"example/internal/container"
	"example/internal/register"
	"log"

	"github.com/spf13/cobra"
)

var oConsumerCommand = &cobra.Command{
	Use:   "consumer",
	Short: "啟動 consumer 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oContainer, _ := container.InitContainer()

		oConsumerRouter := register.ConsumerInit(oContainer)

		oContext := context.Background()
		if err := oConsumerRouter.Serve(oContext); err != nil {
			log.Printf("consumer stopped: %v", err)
		}

	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oConsumerCommand)
}
