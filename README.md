# skunkworks-tama

## Prerequisites

- Go 1.21 or later (for bootstrapping TamaGo)
- Rust toolchain (for building emulator)
- Make
- Git

## Building

### First Time Setup

Both TamaGo and ZisK are managed as git submodules and will be initialized automatically during the build process. You can also initialize them manually:

```bash
git submodule update --init --depth 1
```

### Using docker

Run `run-docker.sh` to enter the shell with all dependencies for Zisk, TamaGo, Rust installed.

### Build Everything
```bash
make
```

This will:
1. Initialize git submodules (TamaGo and ZisK)
2. Build TamaGo
3. Build the emulator

### Build Components Individually

Initialize submodules:
```bash
make init-submodules
```

Build only TamaGo (automatically initializes submodules):
```bash
make build-tamago
```

Build only ZisK emulator (automatically initializes submodules):
```bash
make build-zisk
```

## Project Structure

- `tamago-go-latest/` - TamaGo compiler (git submodule: https://github.com/eth-act/tamago-go, branch: tamago1.25.2-zkvm-dev)
- `zisk/` - ZisK emulator (git submodule: https://github.com/0xPolygonHermez/zisk.git, branch: feature/transpile_unsupported_instructions)
- `tamaboards/zkvm/` - Board support package
- `tama-programs/` - Example programs
  - `empty/` - Minimal "Hello World" program
- `zisk_precompiles` - Go bindings to Zisk precompiles

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

## Package artifacts

Programs in `tama-programs` directory can be built and packaged using the target
```bash
make package-tama-programs
```

That builds programs in `tama-programs` and creates `tama-programs.tar.gz`
