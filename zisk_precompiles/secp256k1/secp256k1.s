//go:build tamago && riscv64

#include "textflag.h"

TEXT ·syscall_secp256k1_dbl(SB), NOSPLIT, $0-8
  MOV ptr+0(FP), A0

  WORD $0x80452073

  RET

TEXT ·syscall_secp256k1_add(SB), NOSPLIT, $0-8
  MOV ptr+0(FP), A0

  WORD $0x80352073

  RET
