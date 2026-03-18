//go:build tamago && riscv64

#include "textflag.h"
#include "../../internal/csr.h"

TEXT ·dma_inputcpy(SB), NOSPLIT, $0-16
  MOV dst+0(FP), A0
  MOV size+8(FP), A1
  CSRS(0x815, 10)
  WORD $0x00B50033
  RET
