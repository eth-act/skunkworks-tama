//go:build tamago && riscv64

#include "textflag.h"

TEXT ·arith384_mod(SB), NOSPLIT, $0-8
  MOV ptr+0(FP), A0

  WORD $0x80B52073

  RET
