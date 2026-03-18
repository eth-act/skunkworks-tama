//go:build tamago && riscv64

#include "textflag.h"
#include "../../internal/csr.h"

TEXT ·dma_memcpy(SB), NOSPLIT, $0-24
  MOV src+8(FP), A0
  MOV dst+0(FP), A1
  MOV size+16(FP), A2
  CSRS(0x813, 10)
  WORD $0x00C58033
  RET
