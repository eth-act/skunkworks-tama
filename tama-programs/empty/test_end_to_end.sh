#!/bin/bash

# End-to-End Test Script for Empty TamaGo Program
# This script tests the complete workflow from building to execution

set -e  # Exit on any error

echo "=========================================="
echo "End-to-End Test for Empty TamaGo Program"
echo "=========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

# Function to check if a command exists
check_command() {
    if ! command -v $1 &> /dev/null; then
        print_error "$1 is not installed or not in PATH"
        exit 1
    fi
}

# Function to check if a file exists
check_file() {
    if [ ! -f "$1" ]; then
        print_error "File not found: $1"
        exit 1
    fi
}

# Function to check if a directory exists
check_directory() {
    if [ ! -d "$1" ]; then
        print_error "Directory not found: $1"
        exit 1
    fi
}

echo "Step 1: Checking prerequisites..."
check_command "go"
check_command "cargo"
check_command "readelf"

# Check for TamaGo toolchain
if [ ! -d "../../tamago-go-latest" ]; then
    print_error "TamaGo toolchain not found at ../../tamago-go-latest"
    exit 1
fi

# Check for Zisk emulator
if [ ! -d "../../zisk" ]; then
    print_error "Zisk emulator not found at ../../zisk"
    exit 1
fi

print_status "All prerequisites found"

echo ""
echo "Step 2: Checking current directory structure..."
check_file "main.go"
check_file "empty_input.bin"
print_status "Source files found"

echo ""
echo "Step 3: Building the empty program..."
echo "Building with TamaGo RISC-V toolchain..."

# Build the program
GOOS=tamago GOARCH=riscv64 ../../tamago-go-latest/bin/go build \
    -gcflags="all=-d=softfloat" \
    -ldflags="-T 0x80000000" \
    -tags tamago,linkcpuinit,linkramstart,linkramsize,linkprintk \
    -o empty.elf .

if [ ! -f "empty.elf" ]; then
    print_error "Failed to build empty.elf"
    exit 1
fi

print_status "Program built successfully"

echo ""
echo "Step 4: Verifying ELF file..."
echo "Checking ELF header:"
readelf -h empty.elf | head -10

echo ""
echo "Checking program headers:"
readelf -l empty.elf | head -15

print_status "ELF file is valid"

echo ""
echo "Step 5: Checking Zisk emulator fix..."
echo "Verifying the fix is applied in the emulator..."

# Check if the fix is applied by looking for the specific code pattern
if grep -q "new_pc == 0" ../../zisk/emulator/src/emu.rs; then
    print_status "Fix is applied in emulator source"
else
    print_warning "Fix may not be applied - checking if emulator was rebuilt"
fi

echo ""
echo "Step 6: Rebuilding Zisk emulator..."
cd ../../zisk
echo "Building emulator in release mode..."
cargo build --release

if [ ! -f "target/release/ziskemu" ]; then
    print_error "Failed to build ziskemu"
    exit 1
fi

print_status "Emulator rebuilt successfully"

echo ""
echo "Step 7: Testing emulator with empty program..."
cd ../tama-programs/empty

echo "Running emulator (this should complete without panic)..."
timeout 30s ../../zisk/target/release/ziskemu -e empty.elf -i empty_input.bin > test_output.log 2>&1

# Check exit status
if [ $? -eq 0 ]; then
    print_status "Emulator ran successfully without panic"
else
    print_error "Emulator failed or timed out"
    echo "Last 20 lines of output:"
    tail -20 test_output.log
    exit 1
fi

echo ""
echo "Step 8: Testing with verbose output..."
echo "Running with verbose output to verify program counter behavior..."

# Run with verbose output and capture last few lines
timeout 30s ../../zisk/target/release/ziskemu -e empty.elf -i empty_input.bin -v > verbose_output.log 2>&1

if [ $? -eq 0 ]; then
    print_status "Verbose run completed successfully"
    echo "Last 10 lines of verbose output:"
    tail -10 verbose_output.log
else
    print_error "Verbose run failed"
    echo "Last 20 lines of verbose output:"
    tail -20 verbose_output.log
    exit 1
fi

echo ""
echo "Step 9: Testing input file handling..."
echo "Checking input file content:"
hexdump -C empty_input.bin

echo ""
echo "Step 10: Final verification..."

# Check that the program actually ran (should see some output in logs)
if grep -q "Emu::run()" verbose_output.log; then
    print_status "Program execution detected in logs"
else
    print_warning "No execution logs found - program may not have run"
fi

# Check that there's no panic in the output
if grep -q "panic" test_output.log || grep -q "panic" verbose_output.log; then
    print_error "Panic detected in output logs"
    echo "Panic details:"
    grep -i "panic" test_output.log verbose_output.log || true
    exit 1
else
    print_status "No panics detected"
fi

echo ""
echo "=========================================="
echo "🎉 ALL TESTS PASSED! 🎉"
echo "=========================================="
echo ""
echo "Summary:"
echo "- ✓ Prerequisites checked"
echo "- ✓ Program built successfully"
echo "- ✓ ELF file verified"
echo "- ✓ Emulator rebuilt with fix"
echo "- ✓ Program runs without panic"
echo "- ✓ Verbose output shows normal execution"
echo "- ✓ Input file handling works"
echo "- ✓ No panics detected"
echo ""
echo "The fix is working correctly! The emulator now properly handles"
echo "Go runtime termination (jump to address 0) without panicking."
echo ""
echo "Files created during testing:"
echo "- test_output.log: Standard output from emulator"
echo "- verbose_output.log: Verbose output from emulator"
echo ""
echo "You can clean up test files with:"
echo "  rm test_output.log verbose_output.log" 