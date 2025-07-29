# skunkworks-tama

## Prerequisites

- Go 1.21 or later (for bootstrapping TamaGo)
- Rust toolchain (for building emulator)
- Make

## Building

### Build Everything
```bash
make
```

This will:
1. Download and build TamaGo
2. Build the emulator

### Build Components Individually

Build only TamaGo:
```bash
make build-tamago
```

Build only ZisK emulator:
```bash
make build-zisk
```

## Project Structure

- `tamago-go-latest/` - TamaGo compiler
- `zisk/` - emulator source
- `tamaboards/zkvm/` - Board support package
- `tama-programs/` - Example programs
  - `empty/` - Minimal empty program

## Running Programs

### Compile a TamaGo Program without make

Example with the empty program:
```bash
cd tama-programs/empty
GOOS=tamago GOARCH=riscv64 ../../tamago-go-latest/bin/go build -tags tamago -o empty.elf .
```

### Run with ZisK Emulator

Run the compiled program:
```bash
make run-empty
```

Or manually:
```bash
cd tama-programs/empty
../../zisk/target/release/ziskemu -e empty.elf -i empty_input.bin
```

## Environment Variables

After building, the following environment variables are available:
- `TAMAGO` - Path to the TamaGo compiler
- `ZISKEMU` - Path to the ZisK emulator

## Clean

Remove all built artifacts:
```bash
make clean
```