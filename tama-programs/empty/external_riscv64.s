//go:build tamago && riscv64

#include "textflag.h"

// func myExternalFunction(a, b uint32) uint32
TEXT ·myExternalFunction(SB), NOSPLIT, $0-12
    // Go calling convention on RISC-V64:
    // Arguments are passed on stack in this order:
    //   a at 0(FP)  - 4 bytes
    //   b at 4(FP)  - 4 bytes
    //   ret at 8(FP) - 4 bytes

    // C calling convention on RISC-V64:
    // First argument in a0 (X10)
    // Second argument in a1 (X11)
    // Return value in a0 (X10)

    // Load arguments from Go stack
    MOVWU a+0(FP), X10    // Load first arg (a) into a0 (X10)
    MOVWU b+4(FP), X11    // Load second arg (b) into a1 (X11)

    // Call the external C function
    CALL c_external_func(SB)

    // Store return value back to Go stack
    MOVW X10, ret+8(FP)   // Store a0 (X10) to return value

    RET

// Declare the external C function symbol
// This will be resolved at link time from your .o file
GLOBL c_external_func(SB), RODATA, $8

