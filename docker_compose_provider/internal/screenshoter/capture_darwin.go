//go:build darwin && cgo

package screenshoter

/*
#cgo CFLAGS: -fobjc-arc
#cgo LDFLAGS: -framework AVFoundation -framework CoreMedia -framework CoreVideo -framework Foundation
#include <stdlib.h>
#include "avf_darwin.h"
*/
import "C"

import (
	"context"
	"errors"
	"image"
	"time"
	"unsafe"
)

// warmupFrames are dropped so the camera's auto exposure has time to settle (the very first frames are dark).
const warmupFrames = 10

// captureFrame grabs one frame from the camera sDevice through AVFoundation.
// The wait for the permission prompt and for the frame each stop at oCtx's deadline.
func captureFrame(oCtx context.Context, sDevice string, iWidth, iHeight int) (image.Image, error) {

	iTimeoutMs := 15000
	if oDeadline, ok := oCtx.Deadline(); ok {
		iTimeoutMs = int(time.Until(oDeadline) / time.Millisecond)
	}
	if iTimeoutMs <= 0 {
		return nil, context.DeadlineExceeded
	}

	cDevice := C.CString(sDevice)
	defer C.free(unsafe.Pointer(cDevice))

	var (
		pData                     *C.uchar
		iFrameW, iFrameH, iStride C.int
		pError                    *C.char
	)
	if C.avf_capture_frame(cDevice, C.int(iWidth), C.int(iHeight), warmupFrames, C.int(iTimeoutMs),
		&pData, &iFrameW, &iFrameH, &iStride, &pError) != 0 {
		sMessage := C.GoString(pError)
		C.free(unsafe.Pointer(pError))
		return nil, errors.New(sMessage)
	}
	defer C.free(unsafe.Pointer(pData))

	aBGRA := unsafe.Slice((*byte)(unsafe.Pointer(pData)), int(iStride)*int(iFrameH))

	// AVFoundation gives BGRA; image.RGBA wants RGBA.
	oImage := image.NewRGBA(image.Rect(0, 0, int(iFrameW), int(iFrameH)))
	for y := 0; y < int(iFrameH); y++ {
		aSrc := aBGRA[y*int(iStride) : y*int(iStride)+int(iFrameW)*4]
		aDst := oImage.Pix[y*oImage.Stride : y*oImage.Stride+int(iFrameW)*4]
		for x := 0; x < len(aSrc); x += 4 {
			aDst[x], aDst[x+1], aDst[x+2], aDst[x+3] = aSrc[x+2], aSrc[x+1], aSrc[x], 0xFF
		}
	}
	return oImage, nil
}
