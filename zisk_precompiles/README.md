# Rust and Go syscall parameters types compatibility

The syscall parameter types are a direct translation from Rust to Go.
As it happens their memory representation is the same. Go doesn't provide
any guarantees for the memory layout of Go struct thus the current implementation
is depedent on the compiler version. For the sake of PoC that's sufficient.
Ultimately parameters to syscalls could be represented just as arrays of bytes
that are serialized and deserialized to Go types or CGO could be used if
supported by TamaGo.

# Syscall invocation

All syscall invocations are similar to each other.
The reason is that in Zisk just a single pointer is passed to the syscall.

Let's analyze `keccak`:

```
TEXT ·syscall_keccak(SB), NOSPLIT, $0-8
  MOV ptr+0(FP), A0

  WORD $0x80052073

  RET
```

`MOV ptr+0(FP), A0` loads the pointer argument from stack (GO ABI) to
a register (syscall ABI).

`WORD $0x80052073` puts a binary encoding of CSRS instruction `CSRS 0x800, A0`.
The value 0x800 is associated with `kecccak` syscall. That is just a shorthand
for `CSRRS ZERO, 0x800, A0`. We can't use the latter because 0x800 literal can't
be used in a place where a named CSR register is expected by GO assembler.

The encoding of the instruction `CSRRS rd,offset,rs1` is
`csr[31:20] | rs1[19:15] | funct3[14:12] | rd[11:7] | opcode[6:0]`
which in our case yields `0x80052073`.

# Testing
Use `make test-zisk-precompiles` target to execute tests.

Basic tests are provided for all syscalls except for sha256.

