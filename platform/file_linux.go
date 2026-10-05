package main

import (
	"os"
)

const pathSeparator = '/'

func platformConfigDir() string {
	if sDir := os.Getenv("XDG_CONFIG_HOME"); sDir != "" {
		return sDir
	}
	sHome, _ := os.UserHomeDir()
	return sHome + "/.config"
}

func platformOpenCommand() []string {
	return []string{"xdg-open"}
}
