package main

import (
	"os"
)

const pathSeparator = '\\'

func platformConfigDir() string {
	if sDir := os.Getenv("APPDATA"); sDir != "" {
		return sDir
	}
	sHome, _ := os.UserHomeDir()
	return sHome + `\AppData\Roaming`
}

func platformOpenCommand() []string {
	return []string{"cmd", "/c", "start", ""}
}
