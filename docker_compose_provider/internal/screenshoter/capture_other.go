//go:build !(darwin && cgo)

package screenshoter

import (
	"context"
	"errors"
	"image"
)

// captureFrame is only implemented on macOS with cgo enabled.
func captureFrame(oCtx context.Context, sDevice string, iWidth, iHeight int) (image.Image, error) {
	return nil, errors.New("這個版本不能截圖:只有 macOS 且啟用 cgo 編譯的版本才支援")
}
