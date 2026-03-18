//go:build tamago && riscv64

#include "textflag.h"
#include "../../internal/csr.h"

TEXT ·fcallParam1(SB), NOSPLIT, $0-8
  MOV val+0(FP), A0
  CSRS(0x8F0, 10)
  RET

TEXT ·fcallGet(SB), NOSPLIT, $0-8
  CSRR(0xFFE)
  MOV T0, ret+0(FP)
  RET

TEXT ·fcallTrigger(SB), NOSPLIT, $0-0
  CSRWI(0x8C0, 18)
  RET
