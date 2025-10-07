#!/bin/bash
set -e

# Configuration
ZISK_REPO="https://github.com/kevaundray/zisk.git"
ZISK_BRANCH="kw/embed-softfloat"
ZISK_DIR="zisk"

echo "=== Setting up ZisK ==="

# Check if zisk directory already exists
if [ -d "$ZISK_DIR" ]; then
    echo "ZisK directory already exists at $ZISK_DIR"

    if [ -d "$ZISK_DIR/.git" ]; then
        echo "Updating existing ZisK repository..."
        cd "$ZISK_DIR"
        git fetch origin
        git checkout "$ZISK_BRANCH"
        git pull origin "$ZISK_BRANCH"
        cd ..
    else
        echo "ZisK directory exists (not a git repo), skipping setup..."
    fi
else
    echo "Cloning ZisK repository (shallow clone, no history)..."
    git clone --depth 1 --branch "$ZISK_BRANCH" "$ZISK_REPO" "$ZISK_DIR"
fi

echo "=== ZisK setup complete ==="
if [ -d "$ZISK_DIR/.git" ]; then
    echo "Branch: $ZISK_BRANCH"
    cd "$ZISK_DIR" && git log -1 --oneline && cd ..
fi
