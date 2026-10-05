package main

// cgo 的 Hello World：Go 呼叫 C 寫的函式。  go run ./sample/cgo_1
//
// C 程式直接寫在 import "C" 正上方的註解裡，這段註解就是 C 程式（cgo 稱為 preamble）。
// 裡面定義的函式，在 Go 這邊用 C.函式名 呼叫。
//
// 這個範例不需要 #cgo 指令，因為只用到 C 標準函式庫（stdio、stdlib），C 編譯器預設就找得到。
// 需要外部函式庫時才要 #cgo CFLAGS / LDFLAGS，範例見 sample/cgo。
//
// 注意：
//   - 說明和下面的 C 程式之間要空一行；註解緊貼 import "C" 的話，會被當成 C 程式。
//   - C 的 printf 和 Go 的 fmt.Println 各自有輸出緩衝區，混著印順序可能亂掉，
//     所以 C 函式印完要 fflush(stdout)。

/*
#include <stdio.h>
#include <stdlib.h>

static void hello(void) {
	printf("Hello, World! (來自 C)\n");
	fflush(stdout);
}

static void hello_name(const char* name) {
	printf("Hello, %s! (來自 C)\n", name);
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

	// 1. 不帶參數
	C.hello()

	// 2. 帶字串參數：Go 字串要轉成 C 字串，用完要 free。
	pName := C.CString("cgo")
	defer C.free(unsafe.Pointer(pName))
	C.hello_name(pName)

	fmt.Println("Go: C 呼叫完了")
}
