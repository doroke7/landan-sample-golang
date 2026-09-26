// Package screenshoter takes one still picture from the Mac camera with the host's ffmpeg.
package screenshoter

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Take saves one frame from the avfoundation video device sDevice to sOutput
// ("{time}" in sOutput is replaced by a timestamp) and returns the absolute path written.
func Take(oCtx context.Context, sDevice, sOutput string) (string, error) {

	if sDevice == "" {
		sDevice = "default"
	}
	if sOutput == "" {
		sOutput = "./runtime/desktop/shot-{time}.jpg"
	}
	sOutput = strings.ReplaceAll(sOutput, "{time}", time.Now().Format("20060102-150405"))

	sOutput, err := filepath.Abs(sOutput)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(sOutput), 0o755); err != nil {
		return "", err
	}
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return "", fmt.Errorf("找不到 ffmpeg,請先安裝(brew install ffmpeg)")
	}

	oCtx, fnCancel := context.WithTimeout(oCtx, 20*time.Second)
	defer fnCancel()

	// 只取 1 張。framerate/video_size 是 Mac 內建相機支援的模式,別台不一定支援。
	oCmd := exec.CommandContext(oCtx, "ffmpeg",
		"-hide_banner", "-loglevel", "error",
		"-f", "avfoundation",
		"-framerate", "30",
		"-video_size", "1280x720",
		"-i", sDevice+":none",
		"-frames:v", "1",
		"-y", sOutput,
	)

	aOutput, err := oCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg 失敗: %v\n%s", err, aOutput)
	}
	return sOutput, nil
}
