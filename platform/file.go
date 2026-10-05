package main

import (
	"fmt"
	"runtime"
)

/*
	檔名後綴決定平台：
	  file_windows.go  只在 GOOS=windows 編譯
	  file_linux.go    只在 GOOS=linux 編譯
	  file_darwin.go   只在 GOOS=darwin 編譯
	  file.go          沒有後綴，所有平台都編譯，放共用邏輯與 main

	後綴只能是 Go 認得的 GOOS / GOARCH（例如 _windows、_amd64、_linux_arm64）。
	每個平台檔都要實作同一組函式（這裡是 platformConfigDir、platformOpenCommand），
	不然那個平台編譯不過。

	檔名規則不夠用時（例如 linux 和 darwin 共用一份），改用檔案開頭的 build tag：
	  //go:build linux || darwin
	build tag 前後要空一行，且要在 package 之前。

	用法：
	  go run ./sample/platform
	  GOOS=windows go build -o /dev/null ./sample/platform   # 交叉編譯，確認每個平台都能過
	  GOOS=linux   go build -o /dev/null ./sample/platform
*/

func main() {
	fmt.Println("os:        ", runtime.GOOS)
	fmt.Println("config dir:", ConfigDir("casino-edge"))

	aCommand := OpenCommand("/tmp")
	fmt.Println("open cmd:  ", aCommand)
}

// ConfigDir 是共用邏輯：每個平台的根目錄不同，由 platformConfigDir 提供，拼接的部分一樣。
func ConfigDir(sApp string) string {
	return platformConfigDir() + string(pathSeparator) + sApp
}

// OpenCommand 回傳用系統預設程式打開 sPath 的指令（不執行，只印出來給你看）。
func OpenCommand(sPath string) []string {
	return append(platformOpenCommand(), sPath)
}
