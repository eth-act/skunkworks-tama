//go:build tamago && riscv64

#include "textflag.h"
#include "../../internal/csr.h"

TEXT ·fcall_input(SB), NOSPLIT, $0-8
  MOV addr+0(FP), A0

  CSRS(0x8F0, 10)

  CSRWI(0x8C0, 22)

  RET
