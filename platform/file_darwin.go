package main

import (
	"os"
)

const pathSeparator = '/'

func platformConfigDir() string {
	sHome, _ := os.UserHomeDir()
	return sHome + "/Library/Application Support"
}

func platformOpenCommand() []string {
	return []string{"open"}
}
