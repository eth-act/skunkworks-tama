//go:build tamago && riscv64

#include "textflag.h"
#include "../../internal/csr.h"

TEXT ·dma_memcmp(SB), NOSPLIT, $0-32
  MOV b+8(FP), A0
  MOV a+0(FP), A1
  MOV size+16(FP), A2
  CSRS(0x814, 10)
  WORD $0x00C582B3
  MOV T0, ret+24(FP)
  RET
