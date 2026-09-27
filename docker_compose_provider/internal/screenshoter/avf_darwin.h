//go:build darwin && cgo

#ifndef AVF_DARWIN_H
#define AVF_DARWIN_H

// avf_capture_frame grabs one BGRA frame from a camera through AVFoundation.
//
// cDevice is "default", a zero-based index, or part of the camera's name.
// iWarmup frames are thrown away first so auto exposure can settle.
// iTimeoutMs bounds each wait (the permission prompt, then the frame).
//
// It returns 0 on success: *oData is a malloc'd buffer of oStride*oHeight bytes the caller must free.
// On failure it returns non-zero and *oError is a malloc'd message the caller must free.
int avf_capture_frame(const char *cDevice, int iWidth, int iHeight, int iWarmup, int iTimeoutMs,
                      unsigned char **oData, int *oWidth, int *oHeight, int *oStride, char **oError);

#endif
