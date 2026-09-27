//go:build darwin && cgo

#import <AVFoundation/AVFoundation.h>
#import <CoreMedia/CoreMedia.h>
#import <CoreVideo/CoreVideo.h>
#include <stdlib.h>
#include <string.h>

#include "avf_darwin.h"

#pragma clang diagnostic ignored "-Wdeprecated-declarations"

// FrameGrabber keeps the first frame that arrives after the warm-up frames.
@interface FrameGrabber : NSObject <AVCaptureVideoDataOutputSampleBufferDelegate>
@property(nonatomic) int warmup;
@property(nonatomic, strong) dispatch_semaphore_t done;
@property(nonatomic) unsigned char *data; // malloc'd BGRA copy, NULL until captured
@property(nonatomic) int width;
@property(nonatomic) int height;
@property(nonatomic) int stride;
@end

@implementation FrameGrabber

- (void)captureOutput:(AVCaptureOutput *)oOutput
    didOutputSampleBuffer:(CMSampleBufferRef)oSample
           fromConnection:(AVCaptureConnection *)oConnection {
    if (self.data != NULL) {
        return;
    }
    if (self.warmup > 0) {
        self.warmup--;
        return;
    }

    CVImageBufferRef oImage = CMSampleBufferGetImageBuffer(oSample);
    if (oImage == NULL) {
        return;
    }

    CVPixelBufferLockBaseAddress(oImage, kCVPixelBufferLock_ReadOnly);
    size_t iWidth = CVPixelBufferGetWidth(oImage);
    size_t iHeight = CVPixelBufferGetHeight(oImage);
    size_t iStride = CVPixelBufferGetBytesPerRow(oImage);
    unsigned char *aCopy = malloc(iStride * iHeight);
    if (aCopy != NULL) {
        memcpy(aCopy, CVPixelBufferGetBaseAddress(oImage), iStride * iHeight);
    }
    CVPixelBufferUnlockBaseAddress(oImage, kCVPixelBufferLock_ReadOnly);

    if (aCopy == NULL) {
        return;
    }
    self.width = (int)iWidth;
    self.height = (int)iHeight;
    self.stride = (int)iStride;
    self.data = aCopy;
    dispatch_semaphore_signal(self.done);
}

@end

static char *copy_error(NSString *sMessage) {
    return strdup([sMessage UTF8String]);
}

// find_camera resolves "default", an index, or a name fragment to a video device.
static AVCaptureDevice *find_camera(const char *cDevice) {
    NSString *sDevice = [NSString stringWithUTF8String:cDevice];
    if (sDevice.length == 0 || [sDevice isEqualToString:@"default"]) {
        return [AVCaptureDevice defaultDeviceWithMediaType:AVMediaTypeVideo];
    }

    AVCaptureDeviceDiscoverySession *oSession = [AVCaptureDeviceDiscoverySession
        discoverySessionWithDeviceTypes:@[ AVCaptureDeviceTypeBuiltInWideAngleCamera, AVCaptureDeviceTypeExternalUnknown ]
                              mediaType:AVMediaTypeVideo
                               position:AVCaptureDevicePositionUnspecified];
    NSArray<AVCaptureDevice *> *aDevices = oSession.devices;

    NSScanner *oScanner = [NSScanner scannerWithString:sDevice];
    NSInteger iIndex = 0;
    if ([oScanner scanInteger:&iIndex] && oScanner.atEnd) {
        return (iIndex >= 0 && iIndex < (NSInteger)aDevices.count) ? aDevices[iIndex] : nil;
    }

    for (AVCaptureDevice *oDevice in aDevices) {
        if ([oDevice.localizedName rangeOfString:sDevice options:NSCaseInsensitiveSearch].location != NSNotFound) {
            return oDevice;
        }
    }
    return nil;
}

static NSString *preset_for(int iWidth, int iHeight) {
    if (iWidth == 1920 && iHeight == 1080) {
        return AVCaptureSessionPreset1920x1080;
    }
    if (iWidth == 1280 && iHeight == 720) {
        return AVCaptureSessionPreset1280x720;
    }
    if (iWidth == 640 && iHeight == 480) {
        return AVCaptureSessionPreset640x480;
    }
    return AVCaptureSessionPresetHigh;
}

// ensure_permission asks for camera access when it was never decided. It returns NULL when access is granted.
static NSString *ensure_permission(int iTimeoutMs) {
    AVAuthorizationStatus iStatus = [AVCaptureDevice authorizationStatusForMediaType:AVMediaTypeVideo];

    if (iStatus == AVAuthorizationStatusNotDetermined) {
        dispatch_semaphore_t oAsked = dispatch_semaphore_create(0);
        [AVCaptureDevice requestAccessForMediaType:AVMediaTypeVideo
                                 completionHandler:^(BOOL bGranted) {
                                     dispatch_semaphore_signal(oAsked);
                                 }];
        if (dispatch_semaphore_wait(oAsked, dispatch_time(DISPATCH_TIME_NOW, (int64_t)iTimeoutMs * NSEC_PER_MSEC)) != 0) {
            return @"等不到攝影機權限的回應,請確認是否有跳出權限對話框";
        }
        iStatus = [AVCaptureDevice authorizationStatusForMediaType:AVMediaTypeVideo];
    }

    if (iStatus != AVAuthorizationStatusAuthorized) {
        return @"沒有攝影機權限,請到 系統設定 > 隱私權與安全性 > 攝影機,允許啟動這個程式的終端機";
    }
    return nil;
}

int avf_capture_frame(const char *cDevice, int iWidth, int iHeight, int iWarmup, int iTimeoutMs,
                      unsigned char **oData, int *oWidth, int *oHeight, int *oStride, char **oError) {
    @autoreleasepool {
        NSString *sProblem = ensure_permission(iTimeoutMs);
        if (sProblem != nil) {
            *oError = copy_error(sProblem);
            return 1;
        }

        AVCaptureDevice *oCamera = find_camera(cDevice);
        if (oCamera == nil) {
            *oError = copy_error([NSString stringWithFormat:@"找不到攝影機 %s", cDevice]);
            return 1;
        }

        NSError *oInputError = nil;
        AVCaptureDeviceInput *oInput = [AVCaptureDeviceInput deviceInputWithDevice:oCamera error:&oInputError];
        if (oInput == nil) {
            *oError = copy_error([NSString stringWithFormat:@"無法開啟攝影機 %@: %@", oCamera.localizedName,
                                                            oInputError.localizedDescription]);
            return 1;
        }

        AVCaptureSession *oSession = [[AVCaptureSession alloc] init];
        if (![oSession canAddInput:oInput]) {
            *oError = copy_error(@"攝影機無法加入擷取流程");
            return 1;
        }
        [oSession addInput:oInput];

        FrameGrabber *oGrabber = [[FrameGrabber alloc] init];
        oGrabber.warmup = iWarmup;
        oGrabber.done = dispatch_semaphore_create(0);

        dispatch_queue_t oQueue = dispatch_queue_create("frame-grabber", DISPATCH_QUEUE_SERIAL);
        AVCaptureVideoDataOutput *oOutput = [[AVCaptureVideoDataOutput alloc] init];
        oOutput.videoSettings = @{(id)kCVPixelBufferPixelFormatTypeKey : @(kCVPixelFormatType_32BGRA)};
        oOutput.alwaysDiscardsLateVideoFrames = YES;
        [oOutput setSampleBufferDelegate:oGrabber queue:oQueue];
        if (![oSession canAddOutput:oOutput]) {
            *oError = copy_error(@"無法取得攝影機影像輸出");
            return 1;
        }
        [oSession addOutput:oOutput];

        // Set after the input is attached, otherwise the device's own default format can win.
        NSString *sPreset = preset_for(iWidth, iHeight);
        if ([oSession canSetSessionPreset:sPreset]) {
            oSession.sessionPreset = sPreset;
        }

        [oSession startRunning];
        dispatch_semaphore_wait(oGrabber.done, dispatch_time(DISPATCH_TIME_NOW, (int64_t)iTimeoutMs * NSEC_PER_MSEC));
        [oSession stopRunning];

        // Let a callback that is still running finish before reading what it stored.
        dispatch_sync(oQueue, ^{
        });

        if (oGrabber.data == NULL) {
            *oError = copy_error(@"讀不到攝影機畫面,請確認沒有被其他程式佔用");
            return 1;
        }

        *oData = oGrabber.data;
        *oWidth = oGrabber.width;
        *oHeight = oGrabber.height;
        *oStride = oGrabber.stride;
        return 0;
    }
}
