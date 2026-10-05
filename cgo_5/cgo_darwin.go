package main

// 只在 GOOS=darwin 編譯。#cgo CFLAGS 是傳給 C 編譯器的參數，這裡定義巨集 PLATFORM。
// 值不加引號（cgo 不允許），main.go 的 C 端會用 # 轉成字串。

/*
#cgo CFLAGS: -DPLATFORM=macOS
*/
import "C"
