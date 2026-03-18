//go:build tamago && riscv64

#include "textflag.h"
#include "../../../internal/csr.h"

TEXT ·fcall_inv(SB), NOSPLIT, $0-16
  MOV input+0(FP), A0
  MOV result+8(FP), A1

  CSRS(0x8F4, 10)

  CSRWI(0x8C0, 12)

  CSRR(0xFFE); MOV T0, 0(A1)
  CSRR(0xFFE); MOV T0, 8(A1)
  CSRR(0xFFE); MOV T0, 16(A1)
  CSRR(0xFFE); MOV T0, 24(A1)
  CSRR(0xFFE); MOV T0, 32(A1)
  CSRR(0xFFE); MOV T0, 40(A1)
  CSRR(0xFFE); MOV T0, 48(A1)
  CSRR(0xFFE); MOV T0, 56(A1)
  CSRR(0xFFE); MOV T0, 64(A1)
  CSRR(0xFFE); MOV T0, 72(A1)
  CSRR(0xFFE); MOV T0, 80(A1)
  CSRR(0xFFE); MOV T0, 88(A1)

  RET

TEXT ·fcall_sqrt(SB), NOSPLIT, $0-16
  MOV element+0(FP), A0
  MOV result+8(FP), A1

  CSRS(0x8F5, 10)

  CSRWI(0x8C0, 19)

  CSRR(0xFFE); MOV T0, 0(A1)
  CSRR(0xFFE); MOV T0, 8(A1)
  CSRR(0xFFE); MOV T0, 16(A1)
  CSRR(0xFFE); MOV T0, 24(A1)
  CSRR(0xFFE); MOV T0, 32(A1)
  CSRR(0xFFE); MOV T0, 40(A1)
  CSRR(0xFFE); MOV T0, 48(A1)
  CSRR(0xFFE); MOV T0, 56(A1)
  CSRR(0xFFE); MOV T0, 64(A1)
  CSRR(0xFFE); MOV T0, 72(A1)
  CSRR(0xFFE); MOV T0, 80(A1)
  CSRR(0xFFE); MOV T0, 88(A1)
  CSRR(0xFFE); MOV T0, 96(A1)

  RET
