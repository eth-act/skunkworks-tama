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

### Compile a TamaGo Program

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
- Floating-point instructions (But we can add this)
- Hardware timers (time is simulated/deterministic)
- Hardware Interrupts
- MMU (Memory Management Unit)
- Hardware RNG (random numbers must be deterministic)
- Traditional peripherals (UART, GPIO, network, storage, display)

The only system call is `ecall` for program termination and to call special functions.

## Clean

Remove all built artifacts:
```bash
make clean
```