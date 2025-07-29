# zkVM Board Support

## Key Differences from Hardware Boards

1. **No I/O During Execution**
   - Output is collected in memory and returned at the end
   - Input is provided at the start of execution (not as part of argc/argv)
   - No UART or MMIO access

2. **Deterministic Execution**
   - Fixed time values (no real clock)
   - Deterministic "random" data (in fact, we should disable randomness or just return the clock cycles)

3. **Memory Constraints**
   - Limited to 16MB RAM
   - No special memory regions
   - All memory is regular RAM

4. **No Hardware Features**
   - No CSR access
   - No interrupts (but we can use ECALL like a software interrupt)
   - No privileged modes
   - No FPU (software float only, but we can modify )

## API

### Input/Output

Input is provided by the emulator by filling the memory of the program at a particular address `INPUT_ADDRESS.

Output is left in memory at a particular address `OUTPUT_ADDRESS`

We will supply methods to read these into the program once we have shown it to work with no input. These methods will look approximately like:

- `ReadInput[Data]() []byte` - Read enough input to deserialize into data from `OUTPUT_ADDRESS`
- `WriteOutput([]byte)` - Append to output buffer at address `OUTPUT_ADDRESS`

### Standard Functions
- `Init()` - Initialize the zkVM "board"
- `Shutdown()` - Halt the program

## Usage (Once it all works)

```go
package main

import "tamagotest/boards/zkvm"

func main() {
    zkvm.Init()
    
    // Read input
    input := zkvm.ReadInput[Data]()
    
    // Process...
    result := compute(input)
    
    // Write output
    zkvm.WriteOutput(result)
    
    zkvm.Shutdown()
}
```

## Limitations

1. **No Real-Time Operations**
   - Can't read actual time
   - Can't wait or sleep
   - No timeouts

2. **No External Communication**
   - Can't make network calls
   - Can't read files
   - Can't access hardware

3. **Fixed Resources**
   - Memory is limited
   - Stack size is fixed
   - No dynamic loading