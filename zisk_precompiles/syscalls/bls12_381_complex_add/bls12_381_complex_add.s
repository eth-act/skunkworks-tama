//go:build tamago && riscv64

#include "textflag.h"
#include "../../internal/csr.h"

TEXT ·syscall_bls12_381_complex_add(SB), NOSPLIT, $0-8
  MOV ptr+0(FP), A0

  CSRS(0x80E, 10)

  RET
