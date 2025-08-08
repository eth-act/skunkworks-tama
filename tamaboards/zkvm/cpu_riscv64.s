//go:build tamago && riscv64 && linkcpuinit

#include "textflag.h"

// Entry point - this is where execution starts
TEXT cpuinit(SB),NOSPLIT|NOFRAME,$0
	JMP	_rt0_tamago_start(SB)
