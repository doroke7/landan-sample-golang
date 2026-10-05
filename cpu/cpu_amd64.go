package main

func archName() string {
	return "amd64 (Intel / AMD, x86-64)"
}

// vectorWidth 是一次能處理的位元數：x86-64 保證有 SSE2（128 bit），AVX2 是 256 bit 但要執行時判斷。
func vectorWidth() int {
	return 128
}
