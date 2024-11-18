//go:build !linux
// +build !linux

package vaapi

import "errors"

func GetEncodingCapabilities() (*EncodingCapabilities, error) {
    return nil, errors.New("VA-API is only supported on Linux")
}
