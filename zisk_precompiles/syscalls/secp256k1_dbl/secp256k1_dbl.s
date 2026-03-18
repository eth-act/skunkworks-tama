//go:build tamago && riscv64

#include "textflag.h"
#include "../../internal/csr.h"

TEXT ·syscall_secp256k1_dbl(SB), NOSPLIT, $0-8
  MOV ptr+0(FP), A0

  CSRS(0x804, 10)

  RET
