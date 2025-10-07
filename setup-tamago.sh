#!/bin/bash
set -e

# Configuration
TAMAGO_REPO="https://github.com/eth-act/tamago-go.git"
TAMAGO_BRANCH="zkvm-develop"
TAMAGO_DIR="tamago-go-latest"

echo "=== Setting up TamaGo ==="

# Check if tamago directory already exists
if [ -d "$TAMAGO_DIR" ]; then
    echo "TamaGo directory already exists at $TAMAGO_DIR"

    if [ -d "$TAMAGO_DIR/.git" ]; then
        echo "Updating existing TamaGo repository..."
        cd "$TAMAGO_DIR"
        git fetch origin
        git checkout "$TAMAGO_BRANCH"
        git pull origin "$TAMAGO_BRANCH"
        cd ..
    else
        echo "TamaGo directory exists (not a git repo), skipping setup..."
    fi
else
    echo "Cloning TamaGo repository (shallow clone, no history)..."
    git clone --depth 1 --branch "$TAMAGO_BRANCH" "$TAMAGO_REPO" "$TAMAGO_DIR"

    echo "=== Building TamaGo ==="
    cd "$TAMAGO_DIR/src"
    ./make.bash
    cd ../..
fi

echo "=== TamaGo setup complete ==="
if [ -d "$TAMAGO_DIR/.git" ]; then
    echo "Branch: $TAMAGO_BRANCH"
    cd "$TAMAGO_DIR" && git log -1 --oneline && cd ..
fi
