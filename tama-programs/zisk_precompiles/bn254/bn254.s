//go:build tamago && riscv64

#include "textflag.h"

TEXT ·bn254_dbl(SB), NOSPLIT, $0-8
  MOV ptr+0(FP), A0

  WORD $0x80752073

  RET

TEXT ·bn254_add(SB), NOSPLIT, $0-8
  MOV ptr+0(FP), A0

  WORD $0x80652073

  RET

TEXT ·bn254_complex_add(SB), NOSPLIT, $0-8
  MOV ptr+0(FP), A0

  WORD $0x80852073

  RET

TEXT ·bn254_complex_mul(SB), NOSPLIT, $0-8
  MOV ptr+0(FP), A0

  WORD $0x80A52073

  RET

TEXT ·bn254_complex_sub(SB), NOSPLIT, $0-8
  MOV ptr+0(FP), A0

  WORD $0x80952073

  RET
