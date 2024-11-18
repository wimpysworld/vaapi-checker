//go:build linux
// +build linux

package vaapi

/*
#cgo pkg-config: libva libva-drm

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <va/va.h>
#include <va/va_drm.h>
#include <fcntl.h>
#include <unistd.h>

static int open_device(const char* path) {
    return open(path, O_RDWR);
}

static VADisplay get_display(int fd) {
    return vaGetDisplayDRM(fd);
}

static int check_string_contains(const char* str, const char* substr) {
    return strstr(str, substr) != NULL;
}
*/
import "C"
import (
    "errors"
    "fmt"
    "strings"
    "unsafe"
)

func getDRMDevice() (C.int, error) {
    paths := []string{
        "/dev/dri/renderD128",
        "/dev/dri/renderD129",
        "/dev/dri/card0",
        "/dev/dri/card1",
    }

    for _, path := range paths {
        cPath := C.CString(path)
        fd := C.open_device(cPath)
        C.free(unsafe.Pointer(cPath))
        if fd >= 0 {
            return fd, nil
        }
    }

    return -1, errors.New("no available DRM device found")
}

func GetEncodingCapabilities() (*EncodingCapabilities, error) {
    fd, err := getDRMDevice()
    if err != nil {
        return nil, fmt.Errorf("failed to open DRM device: %v", err)
    }
    defer C.close(fd)

    display := C.get_display(fd)
    if display == nil {
        return nil, errors.New("failed to get VA-API display")
    }
    defer C.vaTerminate(display)

    var major, minor C.int
    if C.vaInitialize(display, &major, &minor) != C.VA_STATUS_SUCCESS {
        return nil, errors.New("failed to initialize VA-API")
    }

    vendorStr := strings.ToLower(C.GoString(C.vaQueryVendorString(display)))

    // Determine implementation
    impl := "Unknown"
    if strings.Contains(vendorStr, "intel") {
        impl = "Intel"
    } else if strings.Contains(vendorStr, "amd") || strings.Contains(vendorStr, "mesa gallium") {
        impl = "AMD"
    }

    var numProfiles C.int
    profiles := make([]C.VAProfile, 20)
    if C.vaQueryConfigProfiles(display, &profiles[0], &numProfiles) != C.VA_STATUS_SUCCESS {
        return nil, errors.New("failed to query profiles")
    }

    var supported []string
    for i := C.int(0); i < numProfiles; i++ {
        var numEntrypoints C.int
        entrypoints := make([]C.VAEntrypoint, 10)

        if C.vaQueryConfigEntrypoints(display, profiles[i], &entrypoints[0], &numEntrypoints) != C.VA_STATUS_SUCCESS {
            continue
        }

        hasEncode := false
        for j := C.int(0); j < numEntrypoints; j++ {
            if entrypoints[j] == C.VAEntrypointEncSlice {
                hasEncode = true
                break
            }
        }

        if hasEncode {
            switch profiles[i] {
            case C.VAProfileH264Main:
                supported = append(supported, "H.264 Main")
            case C.VAProfileH264High:
                supported = append(supported, "H.264 High")
            case C.VAProfileHEVCMain:
                supported = append(supported, "HEVC Main")
            case C.VAProfileVP9Profile0:
                supported = append(supported, "VP9")
            case C.VAProfileAV1Profile0:
                supported = append(supported, "AV1")
            }
        }
    }

    return &EncodingCapabilities{
        Implementation: impl,
        SupportedEncoders: supported,
    }, nil
}
