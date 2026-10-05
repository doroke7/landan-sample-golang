package main

func archName() string {
	return "arm64 (Apple Silicon / ARM)"
}

// vectorWidth 是一次能處理的位元數：arm64 保證有 NEON（128 bit）。
func vectorWidth() int {
	return 128
}
