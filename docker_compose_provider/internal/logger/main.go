// Package logger prints what docker compose expects to read: JSON message lines and the metadata.
package logger

import (
	"encoding/json"
	"fmt"
)

// Metadata tells docker compose which options `up` accepts.
const Metadata = `{
  "description": "用宿主機的 AVFoundation 從攝影機截一張圖",
  "up": {
    "parameters": [
      {"name": "device", "description": "攝影機:default、編號(0、1…)或名稱的一部分(system_profiler SPCameraDataType 可查名稱)", "required": false, "type": "string", "default": "default"},
      {"name": "output", "description": "輸出檔案,可用 {time} 代表時間戳", "required": false, "type": "string", "default": "./runtime/desktop/shot-{time}.jpg"}
    ]
  },
  "down": {"parameters": []}
}`

type message struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func send(sType, sMessage string) {
	aLine, _ := json.Marshal(message{Type: sType, Message: sMessage})
	fmt.Println(string(aLine))
}

// Info shows a line in the output of docker compose.
func Info(sMessage string) { send("info", sMessage) }

// Error reports a failure; the command should then exit with a non-zero code.
func Error(sMessage string) { send("error", sMessage) }
