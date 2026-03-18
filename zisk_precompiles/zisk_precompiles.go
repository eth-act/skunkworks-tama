// Package zisk_precompiles provides Go bindings for ZisK v0.16.0 precompiled operations.
//
// # Syscalls
//
// Syscalls are triggered by a single CSR write instruction (CSRS). The hardware
// reads parameters from a memory region pointed to by register A0, executes the
// operation, and writes the result back to the same memory region (in-place).
//
// Syscall list: https://github.com/0xPolygonHermez/zisk/blob/v0.16.0/definitions/src/syscall.rs.
//
// | Package                          | CSR     | Description                |
// | -------------------------------- | ------- | -------------------------- |
// | `syscalls/keccakf`               | `0x800` | Keccak-f[1600] permutation |
// | `syscalls/arith256`              | `0x801` | (Dh, Dl) = A * B + C       |
// | `syscalls/arith256_mod`          | `0x802` | D = (A * B + C) mod N      |
// | `syscalls/secp256k1_add`         | `0x803` | P1 = P1 + P2               |
// | `syscalls/secp256k1_dbl`         | `0x804` | P = 2 * P                  |
// | `syscalls/sha256f`               | `0x805` | SHA-256 compression        |
// | `syscalls/bn254_curve_add`       | `0x806` | P1 = P1 + P2               |
// | `syscalls/bn254_curve_dbl`       | `0x807` | P = 2 * P                  |
// | `syscalls/bn254_complex_add`     | `0x808` | F1 = F1 + F2 (Fp2)         |
// | `syscalls/bn254_complex_sub`     | `0x809` | F1 = F1 - F2 (Fp2)         |
// | `syscalls/bn254_complex_mul`     | `0x80A` | F1 = F1 * F2 (Fp2)         |
// | `syscalls/arith384_mod`          | `0x80B` | D = (A * B + C) mod N      |
// | `syscalls/bls12_381_curve_add`   | `0x80C` | P1 = P1 + P2               |
// | `syscalls/bls12_381_curve_dbl`   | `0x80D` | P = 2 * P                  |
// | `syscalls/bls12_381_complex_add` | `0x80E` | F1 = F1 + F2 (Fp2)         |
// | `syscalls/bls12_381_complex_sub` | `0x80F` | F1 = F1 - F2 (Fp2)         |
// | `syscalls/bls12_381_complex_mul` | `0x810` | F1 = F1 * F2 (Fp2)         |
// | `syscalls/add256`                | `0x811` | (Cout, C) = A + B + Cin    |
// | `syscalls/poseidon2`             | `0x812` | Poseidon2 permutation      |
// | `syscalls/dma_memcpy`            | `0x813` | memcpy(dst, src, size)     |
// | `syscalls/dma_memcmp`            | `0x814` | memcmp(a, b, size)         |
// | `syscalls/dma_inputcpy`          | `0x815` | Copy input to dst          |
// | `syscalls/dma_memset`            | `0x816` | memset(dst, 0, size)       |
// | `syscalls/secp256r1_add`         | `0x817` | P1 = P1 + P2               |
// | `syscalls/secp256r1_dbl`         | `0x818` | P = 2 * P                  |
// | `syscalls/blake2br`              | `0x819` | BLAKE2b round              |
//
// # Fcalls (Function Calls)
//
// Fcalls are multi-step operations: write parameters via CSR 0x8Fx registers,
// trigger via CSRWI to CSR 0x8C0 with the opcode as the immediate, then read
// results via CSR 0xFFE.
//
// Fcall list: https://github.com/0xPolygonHermez/zisk/blob/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls/mod.rs
//
// | Package                  | Fcall ID | Description                              |
// | ------------------------ | -------- | ---------------------------------------- |
// | `fcalls/secp256k1/fp`    | `1`      | a = a^(-1) mod p                         |
// | `fcalls/secp256k1/fn`    | `2`      | a = a^(-1) mod n                         |
// | `fcalls/secp256k1/fp`    | `3`      | a = sqrt(a) mod p                        |
// | `fcalls/msb_pos_256`     | `4`      | MSB position of 256-bit value            |
// | `fcalls/bn254/fp`        | `6`      | a = a^(-1) mod p                         |
// | `fcalls/bn254/fp2`       | `7`      | a = a^(-1) in Fp2                        |
// | `fcalls/bn254/twist`     | `8`      | Twist point addition + line coefficients |
// | `fcalls/bn254/twist`     | `9`      | Twist point doubling + line coefficients |
// | `fcalls/bls12_381/fp`    | `10`     | a = a^(-1) mod p                         |
// | `fcalls/bls12_381/fp`    | `11`     | a = sqrt(a) mod p                        |
// | `fcalls/bls12_381/fp2`   | `12`     | a = a^(-1) in Fp2                        |
// | `fcalls/bls12_381/twist` | `13`     | Twist point addition + line coefficients |
// | `fcalls/bls12_381/twist` | `14`     | Twist point doubling + line coefficients |
// | `fcalls/msb_pos_384`     | `15`     | MSB position of 384-bit value            |
// | `fcalls/big_int256_div`  | `16`     | (Q, R) = A / B (256-bit)                 |
// | `fcalls/big_int_div`     | `17`     | (Q, R) = A / B (variable length)         |
// | `fcalls/bin_decomp`      | `18`     | Binary decomposition (MSB to LSB)        |
// | `fcalls/bls12_381/fp2`   | `19`     | a = sqrt(a) in Fp2                       |
// | `fcalls/secp256k1/ecdsa` | `20`     | ECDSA signature verification             |
// | `fcalls/secp256r1/ecdsa` | `21`     | ECDSA signature verification             |
// | `fcalls/input`           | `22`     | Read program input word                  |
package zisk_precompiles
