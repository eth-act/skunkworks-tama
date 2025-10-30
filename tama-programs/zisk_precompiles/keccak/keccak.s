//go:build tamago && riscv64

#include "textflag.h"

TEXT ·keccak(SB), NOSPLIT, $0-8
  // Load the pointer argument from stack
  // ptr is at 0(FP) - first argument
  MOV ptr+0(FP), A0

  // Execute CSRS instruction
  // csrs 0x800, a0
  // Which is just a shorthand for CSRRS ZERO, 0x800, A0
  // We can't use the latter because 0x800 is not known to GO assembler
  // That's why we have to use encoding of the instruction
  // csrrs rd,offset,rs1
  // CSRRS format: csr[31:20] | rs1[19:15] | funct3[14:12] | rd[11:7] | opcode[6:0]
  // In our case:
  // csr = 0x800, rs1 = 10 (A0), funct3 = 2 (CSRRS), rd = 0 (ZERO), opcode = 0x73
  WORD $0x80052073  // This is the machine code for: csrs 0x800, a0

  // Return
  RET
