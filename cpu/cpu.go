package main

import (
	"fmt"
	"runtime"
)

/*
	檔名後綴也可以是 GOARCH：
	  cpu_amd64.go   只在 GOARCH=amd64 編譯（Intel / AMD）
	  cpu_arm64.go   只在 GOARCH=arm64 編譯（Apple Silicon、樹莓派 4/5、ARM 伺服器）
	  cpu.go         沒有後綴，所有 CPU 都編譯，放共用邏輯與 main

	OS 和 CPU 可以組合，順序固定是 _GOOS_GOARCH：
	  cpu_darwin_arm64.go   只有 macOS + arm64（Apple Silicon）才編譯
	  cpu_linux_amd64.go    只有 Linux + amd64 才編譯

	每個 CPU 檔都要實作同一組函式（這裡是 archName、vectorWidth），不然那個 CPU 編譯不過。
	amd64、arm64 以外（例如 386、riscv64）也要有檔案，這裡用 cpu_other.go 的 build tag 兜底：
	  //go:build !amd64 && !arm64

	用法：
	  go run ./sample/cpu
	  GOARCH=amd64 go build -o /dev/null ./sample/cpu   # 交叉編譯，確認每種 CPU 都能過
	  GOARCH=arm64 go build -o /dev/null ./sample/cpu
	  GOOS=linux GOARCH=riscv64 go build -o /dev/null ./sample/cpu   # darwin 只有 amd64、arm64，其他 CPU 要搭配別的 GOOS

	不用寫檔案就能知道的資訊：runtime.GOARCH、runtime.NumCPU()
	要判斷 CPU 有沒有某個指令集（AVX2、NEON）用 golang.org/x/sys/cpu，那是「執行時」判斷，
	跟這裡「編譯時」選檔案是兩回事。
*/

func main() {
	fmt.Println("GOOS/GOARCH:", runtime.GOOS+"/"+runtime.GOARCH)
	fmt.Println("NumCPU:     ", runtime.NumCPU())
	fmt.Println("arch:       ", archName())
	fmt.Println("vector bits:", vectorWidth())
	fmt.Println("sum:        ", Sum([]float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}))
}

// Sum 是共用邏輯：每次累加 vectorWidth()/32 個 float32，模擬各 CPU 一次處理的寬度不同。
// 這裡只是範例，實際加速要用組合語言（.s 檔）或 SIMD 套件，檔名一樣用 _amd64.s、_arm64.s 區分。
func Sum(aValues []float32) float32 {
	iStep := vectorWidth() / 32
	var fTotal float32
	for i := 0; i < len(aValues); i += iStep {
		iEnd := min(i+iStep, len(aValues))
		for _, fValue := range aValues[i:iEnd] {
			fTotal += fValue
		}
	}
	return fTotal
}
