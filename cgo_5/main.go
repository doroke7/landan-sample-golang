package main

// cgo 的最小範例，機制和 pkg/inference/openvino 一樣，只是把 OpenVINO 換成系統自帶的東西，
// 任何機器都能直接跑：  go run ./sample/cgo
//
// 這個範例示範三件事：
//
//  1. #cgo CFLAGS：各系統的 cgo_<系統>.go 用 -D 定義巨集 PLATFORM（值是 macOS、Linux、Windows），下面的 C 函式 platform_name() 回傳它。
//     同一份 main.go，在不同系統編譯會得到不同結果。對應 OpenVINO 的 cgo_darwin.go 等檔案
//     （那邊 CFLAGS 是 -I<標頭檔目錄>，這裡是 -D<巨集>，都是傳給 C 編譯器的參數）。
//  2. #cgo LDFLAGS: -lm：連結數學函式庫 libm，C.sqrt、C.pow 的實作在裡面。
//     對應 OpenVINO 的 -lopenvino_c，C.ov_core_create 等函式的實作在 libopenvino_c 裡。
//     拿掉 -lm 的話，在 Linux 上會連結失敗：undefined reference to `sqrt'。
//     （macOS 的 libm 併在系統函式庫裡，拿掉也能過，所以用 Linux 看這個差別。）
//  3. Go 和 C 之間的轉換：Go 字串要用 C.CString 轉成 C 字串，用完一定要 C.free，
//     結果再用 C.GoString 轉回來。model.go 載入模型路徑時也是同樣的寫法。
//
// 注意：說明和下面的 C 程式之間要空一行；註解緊貼 import "C" 的話，會被當成 C 程式。

/*
#cgo LDFLAGS: -lm

#include <math.h>
#include <stdlib.h>
#include <string.h>

// PLATFORM 由 cgo_<系統>.go 的 #cgo CFLAGS: -DPLATFORM=... 提供，這個檔案本身不知道是哪個系統。
// cgo 不允許參數裡有引號（-DPLATFORM="macOS" 會報 malformed #cgo argument），
// 所以 CFLAGS 給的是沒有引號的 macOS，這裡用 # 把它變成字串 "macOS"。
// 要兩層巨集（STR 呼叫 STR_），PLATFORM 才會先展開成 macOS，再被 # 變成字串。
#define STR_(x) #x
#define STR(x) STR_(x)

static const char* platform_name(void) {
	return STR(PLATFORM);
}
*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

func main() {
	// 1. CFLAGS：C 端看到的平台名稱，和 Go 端的 runtime.GOOS 對照
	fmt.Println("GOOS:           ", runtime.GOOS)
	fmt.Println("C PLATFORM 巨集:", C.GoString(C.platform_name()))

	// 2. LDFLAGS：呼叫 libm 的函式。C.double 和 Go 的 float64 要明確轉換。
	fmt.Println("sqrt(2)        =", float64(C.sqrt(C.double(2))))
	fmt.Println("pow(2, 10)     =", float64(C.pow(C.double(2), C.double(10))))

	// 3. 字串轉換：Go → C → Go
	sText := "你好 cgo"
	pText := C.CString(sText) // 在 C 的記憶體配置一份複製，Go 的垃圾回收不會管它
	defer C.free(unsafe.Pointer(pText))
	fmt.Println("strlen(", sText, ") =", int(C.strlen(pText)), "（C 算的是位元組數，不是字數）")
}
