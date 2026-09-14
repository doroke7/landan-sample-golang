package cmd

import (
	"example/bootstrap"
	"example/internal/container"
	"example/internal/register"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/spf13/cobra"
)

var oHttpCommand = &cobra.Command{
	Use:   "http",
	Short: "啟動 Gin HTTP 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oContainer, err := container.InitContainer()
		if err != nil {
			log.Fatal(err)
		}
		oGin := gin.Default()

		oEngine := register.HttpInit(oGin, oContainer)
		sAddress := ":" + bootstrap.CONFIG.HTTP.PORT
		oErr := oEngine.Run(sAddress)
		log.Fatal(oErr)
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oHttpCommand)
}
