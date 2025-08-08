//go:build tamago && riscv64

#include "textflag.h"

// setRegisters sets A0 to INPUT_ADDR and A1 to OUTPUT_ADDR using hardcoded values
TEXT ·setRegisters(SB),NOSPLIT|NOFRAME,$0
	// Hardcoded values - must match board.go constants
	MOV	$0x90000000, A0  // INPUT_ADDR
	MOV	$0xa0010000, A1  // OUTPUT_ADDR
	RET
