#!/bin/bash -e

echo "Build binary..."
GOOS=linux GOARCH=arm go build -v -ldflags="-s -w" .

type upx && {
    echo "Packing using upx..."
    upx -9 hello-exporter
} || echo "upx not available"
