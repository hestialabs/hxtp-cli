#!/usr/bin/env bash

set -euo pipefail

# HestiaLabs CLI Build Script
# Cross-compiles hxtpctl for multiple platforms and architectures.

APP_NAME="hxtp-cli"
VERSION="1.0.0"
DIST_DIR="dist"
BUILD_PATH="./cmd/hxtp-cli"

# Platforms to build for: "os/arch"
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

# Clean and create dist directory
echo "==> Cleaning up old builds..."
rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

echo "==> Starting cross-compilation for version $VERSION..."

for platform in "${PLATFORMS[@]}"; do
    # Split "os/arch"
    IFS="/" read -r OS ARCH <<< "$platform"
    
    BINARY_NAME="$APP_NAME"
    if [[ "$OS" == "windows" ]]; then
        BINARY_NAME="${APP_NAME}.exe"
    fi
    
    OUTPUT_NAME="${APP_NAME}-${OS}-${ARCH}"
    if [[ "$OS" == "windows" ]]; then
        OUTPUT_NAME="${APP_NAME}-${OS}-${ARCH}"
    fi
    
    echo "    Building for $OS/$ARCH..."
    
    # Run go build
    GOOS=$OS GOARCH=$ARCH go build -ldflags "-s -w -X main.version=$VERSION" -o "$DIST_DIR/$BINARY_NAME" "$BUILD_PATH"
    
    # Create archive
    cd "$DIST_DIR"
    if [[ "$OS" == "windows" ]]; then
        # Use zip for Windows
        if command -v zip >/dev/null 2>&1; then
            zip -q "${OUTPUT_NAME}.zip" "$BINARY_NAME"
            echo "    Created ${OUTPUT_NAME}.zip"
        else
            tar -czf "${OUTPUT_NAME}.tar.gz" "$BINARY_NAME"
            echo "    Created ${OUTPUT_NAME}.tar.gz (zip recommended for Windows)"
        fi
    else
        # Use tar.gz for Unix
        tar -czf "${OUTPUT_NAME}.tar.gz" "$BINARY_NAME"
        echo "    Created ${OUTPUT_NAME}.tar.gz"
    fi
    
    # Remove the raw binary to keep dist clean (optional)
    rm "$BINARY_NAME"
    cd ..
done

echo ""
echo " Build complete! Binaries are in the '$DIST_DIR' directory."
ls -lh "$DIST_DIR"
