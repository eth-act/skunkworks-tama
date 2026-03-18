//go:build tamago && riscv64

#include "textflag.h"
#include "../../../internal/csr.h"

TEXT ·fcall_verify(SB), NOSPLIT, $0-40
  MOV pk+0(FP), A0
  MOV z+8(FP), A1
  MOV r+16(FP), A2
  MOV s+24(FP), A3
  MOV result+32(FP), A4

  CSRS(0x8F3, 10)

  MOV A1, A0
  CSRS(0x8F2, 10)

  MOV A2, A0
  CSRS(0x8F2, 10)

  MOV A3, A0
  CSRS(0x8F2, 10)

  CSRWI(0x8C0, 21)

  CSRR(0xFFE); MOV T0, 0(A4)
  CSRR(0xFFE); MOV T0, 8(A4)
  CSRR(0xFFE); MOV T0, 16(A4)
  CSRR(0xFFE); MOV T0, 24(A4)
  CSRR(0xFFE); MOV T0, 32(A4)
  CSRR(0xFFE); MOV T0, 40(A4)
  CSRR(0xFFE); MOV T0, 48(A4)
  CSRR(0xFFE); MOV T0, 56(A4)

  RET
