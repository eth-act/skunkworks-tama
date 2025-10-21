//go:build tamago && riscv64

#include "textflag.h"

TEXT ·sha256(SB), NOSPLIT, $0-8
  MOV ptr+0(FP), A0

  WORD $0x80552073

  RET
