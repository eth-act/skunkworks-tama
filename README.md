# skunkworks-tama

## Prerequisites

- Go 1.21 or later (for bootstrapping TamaGo)
- Rust toolchain (for building emulator)
- Make
- Git

## Building

### First Time Setup

The ZisK emulator is downloaded automatically during the build process. You can also set it up manually:

```bash
./setup-zisk.sh
```

To use a different branch or tag, edit `setup-zisk.sh` and change the `ZISK_BRANCH` variable.

### Build Everything
```bash
make
```

This will:
1. Download/update ZisK emulator (if needed)
2. Build TamaGo
3. Build the emulator

### Build Components Individually

Setup ZisK emulator:
```bash
make setup-zisk
```

Build only TamaGo:
```bash
make build-tamago
```

Build only ZisK emulator:
```bash
make build-zisk
```

## Project Structure

- `tamago-go-latest/` - TamaGo compiler (modified for softfloat support)
- `zisk/` - ZisK emulator (downloaded from https://github.com/0xPolygonHermez/zisk)
- `tamaboards/zkvm/` - Board support package
- `tama-programs/` - Example programs
  - `empty/` - Minimal "Hello World" program
- `setup-zisk.sh` - Script to download/update ZisK emulator

## Running Programs

### Compile the TamaGo Program

Using make (recommended):
```bash
make compile-empty
```

Or manually with all the correct flags:
```bash
cd tama-programs/empty
GOOS=tamago GOARCH=riscv64 ../../tamago-go-latest/bin/go build \
  -gcflags="all=-d=softfloat" \
  -ldflags="-T 0x80000000" \
  -tags tamago,linkcpuinit,linkramstart,linkramsize,linkprintk \
  -o empty.elf .
```

### Run with ZisK Emulator

Run the compiled program:
```bash
# Run with verbose output and console logging
make run-empty-emu

# Run quietly with just console output
make run-empty-emu-quiet
```

### Setup ROM for ZisK

Generate ROM setup for the program:
```bash
make run-empty-rom
```

Or manually:
```bash
cd tama-programs/empty
# Note: The -c flag is required to see console output!
../../zisk/target/debug/ziskemu --elf empty.elf -c
```

## Environment Variables

After building, the following environment variables are available:
- `TAMAGO` - Path to the TamaGo compiler
- `ZISKEMU` - Path to the ZisK emulator

## VM Peripherals

The VM provides minimal "peripherals".

### Memory-mapped I/O regions:
- **Input Buffer** at `0xa0000000` - Where input data is placed for the program
- **Output Buffer** at `0xa0010000` - Where programs write output data  
- **RAM** starting at `0xa0020000` - Main memory for program execution (~512MB)

### Supported features:
- RISC-V RV64IMA instruction set (`c` is not actually in the go compiler yet)
- Simple I/O model: read input → compute → write output → exit

### Not supported:
- Hardware timers (time is simulated/deterministic)
- Hardware Interrupts
- MMU (Memory Management Unit)
- Hardware RNG (random numbers must be deterministic)
- Traditional peripherals (UART, GPIO, network, storage, display)
- Floating-point instructions (compiled with softfloat, any FP instructions will panic the emulator)

The only system call is `ecall` for program termination and to call special functions.

## Debugging with Instruction Tracing

The ZisK emulator provides several tracing options for debugging:

```bash
# Run with console output (REQUIRED to see program output)
ziskemu --elf program.elf -c

# Verbose mode with console output
ziskemu --elf program.elf -v -c

# Log every step
ziskemu --elf program.elf -l

# Print trace every N steps
ziskemu --elf program.elf -p 100

# Save trace to file
ziskemu --elf program.elf -t trace.out

# Enable RISC-V instruction tracing
ziskemu --elf program.elf -a

# Generate statistics
ziskemu --elf program.elf -x
```

**Important**: The `-c` flag is required to see console output from your program!

## Clean

Remove all built artifacts:
```bash
make clean
```