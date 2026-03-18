//go:build tamago && riscv64

#include "textflag.h"
#include "../../internal/csr.h"

TEXT ·fcall_msb_pos_256(SB), NOSPLIT, $0-32
  MOV count+0(FP), A0
  MOV x+8(FP), A1
  MOV y+16(FP), A2
  MOV result+24(FP), A3

  CSRS(0x8F0, 10)

  MOV A1, A0
  CSRS(0x8F2, 10)

  MOV A2, A0
  CSRS(0x8F2, 10)

  CSRWI(0x8C0, 4)

  CSRR(0xFFE); MOV T0, 0(A3)
  CSRR(0xFFE); MOV T0, 8(A3)

  RET
