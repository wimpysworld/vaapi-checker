package main

import (
    "fmt"
    "os"
    "vaapi-checker/pkg/vaapi"
)

func main() {
    encoders, err := vaapi.GetEncodingCapabilities()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("VA-API Hardware Encoding Capabilities:\n")
    fmt.Printf("Implementation: %s\n", encoders.Implementation)
    fmt.Printf("\nSupported Encoding Formats:\n")
    for _, enc := range encoders.SupportedEncoders {
        fmt.Printf("- %s\n", enc)
    }
}
