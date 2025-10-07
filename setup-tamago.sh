#!/bin/bash
set -e

# Configuration
TAMAGO_REPO="https://github.com/eth-act/tamago-go.git"
TAMAGO_BRANCH="zkvm-develop"
TAMAGO_DIR="tamago-go-latest"

echo "=== Setting up TamaGo ==="

if [ -d "$TAMAGO_DIR" ]; then
    echo "TamaGo directory already exists at $TAMAGO_DIR"
else
    echo "Cloning TamaGo repository (shallow clone, no history)..."
    git clone --depth 1 --branch "$TAMAGO_BRANCH" "$TAMAGO_REPO" "$TAMAGO_DIR"

    echo "Building TamaGo..."
    cd "$TAMAGO_DIR/src"
    ./make.bash
    cd ../..
    echo "TamaGo setup complete"
fi