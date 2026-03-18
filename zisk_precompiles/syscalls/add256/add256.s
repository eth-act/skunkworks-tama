//go:build tamago && riscv64

#include "textflag.h"
#include "../../internal/csr.h"

TEXT ·syscall_add256(SB), NOSPLIT, $0-16
  MOV ptr+0(FP), A0
  CSRS_RET(0x811, 10)
  MOV T0, ret+8(FP)
  RET
