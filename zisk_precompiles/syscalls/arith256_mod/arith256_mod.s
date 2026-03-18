//go:build tamago && riscv64

#include "textflag.h"
#include "../../internal/csr.h"

TEXT ·syscall_arith256_mod(SB), NOSPLIT, $0-8
  MOV ptr+0(FP), A0

  CSRS(0x802, 10)

  RET
