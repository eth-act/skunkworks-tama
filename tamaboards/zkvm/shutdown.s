//go:build tamago && riscv64

#include "textflag.h"

// Shutdown triggers ZisK exit via ecall with a7=93
TEXT ·Shutdown(SB),NOSPLIT|NOFRAME,$0
	MOV	$93, A7		// CAUSE_EXIT = 93
	ECALL			// System call to exit
	JMP	-4(PC)		// Infinite loop - ensures we never return
