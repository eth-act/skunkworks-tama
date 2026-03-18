//go:build tamago && riscv64

#include "textflag.h"
#include "../../../internal/csr.h"

TEXT ·fcall_inv(SB), NOSPLIT, $0-16
  MOV input+0(FP), A0
  MOV result+8(FP), A1

  CSRS(0x8F2, 10)

  CSRWI(0x8C0, 1)

  CSRR(0xFFE); MOV T0, 0(A1)
  CSRR(0xFFE); MOV T0, 8(A1)
  CSRR(0xFFE); MOV T0, 16(A1)
  CSRR(0xFFE); MOV T0, 24(A1)

  RET

TEXT ·fcall_sqrt(SB), NOSPLIT, $0-24
  MOV input+0(FP), A0
  MOV parity+8(FP), A1
  MOV result+16(FP), A2

  CSRS(0x8F2, 10)

  MOV A1, A0
  CSRS(0x8F0, 10)

  CSRWI(0x8C0, 3)

  CSRR(0xFFE); MOV T0, 0(A2)
  CSRR(0xFFE); MOV T0, 8(A2)
  CSRR(0xFFE); MOV T0, 16(A2)
  CSRR(0xFFE); MOV T0, 24(A2)
  CSRR(0xFFE); MOV T0, 32(A2)

  RET
