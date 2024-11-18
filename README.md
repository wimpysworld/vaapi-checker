# VA-API Hardware Encoder Checker

A utility to check VA-API hardware encoding capabilities on headless Linux systems.

## Prerequisites

- Go 1.18 or later
- `libva-dev` and `libva-drm-dev` packages installed
- DRM render node access rights

### Ubuntu 22.04

An example of installing the required tool using Ubuntu 22.04.

#### Installing

Install Go

```bash
sudo apt-get install golang-go
```
Install the required libraries

```bash
sudo apt-get install libva-dev libdrm-dev pkg-config
```

Check if `pkg-config` can find `libva` and `libva-drm`

```bash
pkg-config --libs --cflags libva libva-drm
```

Make sure your user has access to DRM render nodes:
```bash
sudo usermod -a -G render,video $USER
```

#### Building

Clear the cache if required.

```bash
go clean -cache
```

Build verbosely.

```bash
go build -v -x ./cmd/vaapi-checker
```

#### Running

```bash
./vaapi-checker
```

The output will look something like this:

```
VA-API Hardware Encoding Capabilities:
Implementation: Intel

Supported Encoding Formats:
- H.264 High
- HEVC Main
- AV1
```
