//go:build tamago && riscv64

#include "textflag.h"
#include "../../internal/csr.h"

TEXT ·fcall_msb_pos_384(SB), NOSPLIT, $0-24
  MOV x+0(FP), A0
  MOV y+8(FP), A1
  MOV result+16(FP), A2

  CSRS(0x8F3, 10)

  MOV A1, A0
  CSRS(0x8F3, 10)

  CSRWI(0x8C0, 15)

  CSRR(0xFFE); MOV T0, 0(A2)
  CSRR(0xFFE); MOV T0, 8(A2)

  RET
