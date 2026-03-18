//go:build tamago && riscv64

#include "textflag.h"
#include "../../internal/csr.h"

TEXT ·syscall_bn254_complex_mul(SB), NOSPLIT, $0-8
  MOV ptr+0(FP), A0

  CSRS(0x80A, 10)

  RET
