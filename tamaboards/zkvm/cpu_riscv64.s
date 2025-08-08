//go:build tamago && riscv64 && linkcpuinit

#include "textflag.h"

// Entry point - this is where execution starts
TEXT cpuinit(SB),NOSPLIT|NOFRAME,$0
	// Clear A0/A1 registers so runtime does not think we have
	// argc/argv arguments
	MOV	$0, A0
	MOV	$0, A1
	JMP	_rt0_tamago_start(SB)
