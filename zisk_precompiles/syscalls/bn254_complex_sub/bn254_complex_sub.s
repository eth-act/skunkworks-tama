//go:build tamago && riscv64

#include "textflag.h"
#include "../../internal/csr.h"

TEXT ·syscall_bn254_complex_sub(SB), NOSPLIT, $0-8
  MOV ptr+0(FP), A0

  CSRS(0x809, 10)

  RET
