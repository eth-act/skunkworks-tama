//go:build tamago && riscv64

#include "textflag.h"
#include "../../internal/csr.h"

TEXT ·fcall_big_int256_div(SB), NOSPLIT, $0-24
  MOV a+0(FP), A0
  MOV b+8(FP), A1
  MOV result+16(FP), A2

  CSRS(0x8F2, 10)

  MOV A1, A0
  CSRS(0x8F2, 10)

  CSRWI(0x8C0, 16)

  CSRR(0xFFE); MOV T0, 0(A2)
  CSRR(0xFFE); MOV T0, 8(A2)
  CSRR(0xFFE); MOV T0, 16(A2)
  CSRR(0xFFE); MOV T0, 24(A2)
  CSRR(0xFFE); MOV T0, 32(A2)
  CSRR(0xFFE); MOV T0, 40(A2)
  CSRR(0xFFE); MOV T0, 48(A2)
  CSRR(0xFFE); MOV T0, 56(A2)

  RET
