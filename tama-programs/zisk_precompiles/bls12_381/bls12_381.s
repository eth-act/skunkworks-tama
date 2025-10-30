//go:build tamago && riscv64

#include "textflag.h"

TEXT ·bls12_381_dbl(SB), NOSPLIT, $0-8
  MOV ptr+0(FP), A0

  WORD $0x80D52073

  RET

TEXT ·bls12_381_add(SB), NOSPLIT, $0-8
  MOV ptr+0(FP), A0

  WORD $0x80C52073

  RET

TEXT ·bls12_381_complex_add(SB), NOSPLIT, $0-8
  MOV ptr+0(FP), A0

  WORD $0x80E52073

  RET

TEXT ·bls12_381_complex_mul(SB), NOSPLIT, $0-8
  MOV ptr+0(FP), A0

  WORD $0x81052073

  RET

TEXT ·bls12_381_complex_sub(SB), NOSPLIT, $0-8
  MOV ptr+0(FP), A0

  WORD $0x80F52073

  RET
