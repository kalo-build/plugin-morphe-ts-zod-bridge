#!/bin/bash
set -euo pipefail

PLUGIN_NAME="morphe-ts-zod-bridge"
VERSION="v1.0.0"
OUTPUT_DIR="dist"

mkdir -p "$OUTPUT_DIR"

echo "Building $PLUGIN_NAME $VERSION..."
GOOS=wasip1 GOARCH=wasm go build \
  -o "$OUTPUT_DIR/${PLUGIN_NAME}-${VERSION}.wasm" \
  ./cmd/plugin/

echo "Build complete: $OUTPUT_DIR/${PLUGIN_NAME}-${VERSION}.wasm"
