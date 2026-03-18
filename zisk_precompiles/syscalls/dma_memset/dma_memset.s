//go:build tamago && riscv64

#include "textflag.h"
#include "../../internal/csr.h"

TEXT ·dma_memset_zero(SB), NOSPLIT, $0-16
  MOV dst+0(FP), A0
  MOV size+8(FP), A1
  CSRS(0x816, 10)
  WORD $0x00058013
  RET
