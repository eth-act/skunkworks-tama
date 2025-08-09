//go:build tamago && riscv64

package zkvm

import (
	"runtime"
	"unsafe"
	_ "github.com/usbarmory/tamago/riscv64"
)

const (
	// ZisK ROM addresses (matching Zisk memory map from mem.rs)
	ROM_ENTRY = 0x1000      // First BIOS instruction address
	ROM_EXIT  = 0x1004      // Last BIOS instruction address
	ROM_ADDR  = 0x80000000  // First program ROM instruction address
	
	// ZisK I/O addresses (matching Zisk memory map from mem.rs)
	INPUT_ADDR  = 0x90000000  // First input data memory address
	SYS_ADDR    = 0xa0000000  // First system RW memory address (RAM_ADDR)
	// TODO: Check: Can BSS be initialized at system address and UART_ADDR overwrite it
	UART_ADDR   = 0xa0000200  // UART memory address (SYS_ADDR + 512)
	OUTPUT_ADDR = 0xa0010000  // First output RW memory address
)

//go:linkname ramStart runtime.ramStart
var ramStart uint64 = 0xa0020000 // Match ZisK's AVAILABLE_MEM_ADDR

//go:linkname ramSize runtime.ramSize
var ramSize uint64 = 0x1FFE0000 // Match ZisK's RAM size (~512MB)

// ramStackOffset is always defined here as there's no linkramstackoffset build tag
//go:linkname ramStackOffset runtime.ramStackOffset
var ramStackOffset uint64 = 0x100000 // 1MB stack (matching linker script and ZisK)

// TODO: We can probably remove this
// Bloc sets the heap start address to bypass initBloc() 
//go:linkname Bloc runtime.Bloc
var Bloc uintptr = 0xa0120000 // Start heap after stack (ramStart + ramStackOffset)

// printk implementation for zkVM
//go:linkname printk runtime.printk
func printk(c byte) {
	*(*byte)(unsafe.Pointer(uintptr(UART_ADDR))) = c
}

// hwinit1 is now defined in hwinit1.s 
// we use it to set A0/A1 registers to the input and output address

// Use this as a stub timer. It is all single threaded, and there is no concept of time.
// This may return the cycle count in the future.
var timer int64 = 0

//go:linkname nanotime1 runtime.nanotime1
func nanotime1() int64 {
	// Return deterministic time for zkVM
	// Could be based on instruction count or fixed increments
	// Or just return number of clock cycles
	timer++
	return timer * 1000
}

// Init initializes the zkVM 
//go:linkname Init runtime.hwinit1
func Init() {
	// Set up custom exit handler
	runtime.Exit = zkVMExit
}

//go:linkname initRNG runtime.initRNG
func initRNG() {
	// Deterministic RNG initialization
	// TODO: There is no "proper" rng so nothing to init.
}

//go:linkname getRandomData runtime.getRandomData
func getRandomData(b []byte) {
	// Deterministic "random" data
	// In a real zkVM, this might come from the input
	for i := range b {
		b[i] = byte(i & 0xFF)
	}
}

// zkVMExit is our custom exit handler that intercepts runtime.exit calls
func zkVMExit(code int32) {
	// TODO: Add a println about what the exit code is
	// Zisk exit via ecall
	Shutdown()
}

// Shutdown is defined in shutdown.s and uses ecall to exit
func Shutdown()

// setRegisters is defined in hwinit1.s and sets A0/A1 registers
// A0 = INPUT_ADDR (0x90000000)
// A1 = OUTPUT_ADDR (0xa0010000)
func setRegisters()