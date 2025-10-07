#!/bin/bash
set -e

# Configuration
ZISK_REPO="https://github.com/kevaundray/zisk.git"
ZISK_BRANCH="kw/embed-softfloat"
ZISK_DIR="zisk"

echo "=== Setting up ZisK ==="

if [ -d "$ZISK_DIR" ]; then
    echo "ZisK directory already exists at $ZISK_DIR"
else
    echo "Cloning ZisK repository (shallow clone, no history)..."
    git clone --depth 1 --branch "$ZISK_BRANCH" "$ZISK_REPO" "$ZISK_DIR"
    echo "ZisK setup complete"
fi
