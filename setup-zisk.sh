#!/bin/bash
set -e

# Configuration
ZISK_REPO="https://github.com/0xPolygonHermez/zisk.git"
ZISK_BRANCH="feature/float_riscv_instructions"
ZISK_DIR="zisk"

echo "=== Setting up ZisK ==="

# Check if zisk directory already exists
if [ -d "$ZISK_DIR" ]; then
    echo "ZisK directory already exists at $ZISK_DIR"
    echo "Checking if it's a git repository..."

    if [ -d "$ZISK_DIR/.git" ]; then
        echo "Updating existing ZisK repository..."
        cd "$ZISK_DIR"
        git fetch origin
        git checkout "$ZISK_BRANCH"
        git pull origin "$ZISK_BRANCH"
        cd ..
    else
        echo "ERROR: $ZISK_DIR exists but is not a git repository"
        echo "Please remove it manually and run this script again"
        exit 1
    fi
else
    echo "Cloning ZisK repository..."
    git clone --branch "$ZISK_BRANCH" "$ZISK_REPO" "$ZISK_DIR"
fi

echo "=== ZisK setup complete ==="
echo "Branch: $ZISK_BRANCH"
cd "$ZISK_DIR" && git log -1 --oneline && cd ..
