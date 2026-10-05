//go:build !amd64 && !arm64

package main

func archName() string {
	return "other"
}

func vectorWidth() int {
	return 32
}
