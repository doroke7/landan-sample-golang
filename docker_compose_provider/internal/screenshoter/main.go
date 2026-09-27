// Package screenshoter takes one still picture from the Mac camera through AVFoundation (no ffmpeg needed).
package screenshoter

import (
	"context"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Take saves one frame from the camera sDevice to sOutput
// ("{time}" in sOutput is replaced by a timestamp) and returns the absolute path written.
// sDevice is "default", an index, or part of the camera's name.
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

	// 留夠時間給第一次執行時使用者按下攝影機權限的對話框。
	oCtx, fnCancel := context.WithTimeout(oCtx, 30*time.Second)
	defer fnCancel()

	oImage, err := captureFrame(oCtx, sDevice, 1920, 1080)
	if err != nil {
		return "", err
	}

	oFile, err := os.Create(sOutput)
	if err != nil {
		return "", err
	}
	if err := jpeg.Encode(oFile, oImage, &jpeg.Options{Quality: 90}); err != nil {
		oFile.Close()
		return "", err
	}
	if err := oFile.Close(); err != nil {
		return "", err
	}
	return sOutput, nil
}
