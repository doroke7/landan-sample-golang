package main

// cgo 的 Hello World：Go 呼叫 C 寫的函式。  go run ./sample/cgo_1
//
// 兩種寫 C 程式的方式：
//
//  1. 直接寫在 import "C" 上面的註解裡（下面的 hello_inline）。適合很短的 C 程式。
//  2. 寫成獨立的 .c 和 .h 檔，放在同一個目錄（hello.c、hello.h）。
//     go build 會自動把目錄裡的 .c 檔一起編譯，Go 這邊只要 #include "hello.h" 就能呼叫。
//
// 這個範例不需要 #cgo 指令，因為只用到 C 標準函式庫（stdio），C 編譯器預設就找得到。
// 需要外部函式庫時才要 #cgo CFLAGS / LDFLAGS，範例見 sample/cgo。
//
// 注意：
//   - 說明和下面的 C 程式之間要空一行；註解緊貼 import "C" 的話，會被當成 C 程式。
//   - C 的 printf 和 Go 的 fmt.Println 各自有輸出緩衝區，混著印順序可能亂掉，
//     所以 C 函式印完要 fflush(stdout)。

/*
#include <stdio.h>
#include <stdlib.h>

#include "hello.h"

static void hello_inline(void) {
	printf("Hello, World! (來自 main.go 裡的 C)\n");
	fflush(stdout);
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println("Go: 準備呼叫 C")

	// 1. 呼叫寫在註解裡的 C 函式
	C.hello_inline()

	// 2. 呼叫 hello.c 的函式。Go 字串要轉成 C 字串，用完要 free。
	pName := C.CString("cgo")
	defer C.free(unsafe.Pointer(pName))
	C.hello_from_file(pName)

	fmt.Println("Go: C 呼叫完了")
}
