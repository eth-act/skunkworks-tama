// CSRS: csrrs x0, csr, rs1 -- set CSR bits from register (syscall trigger, fcall param write)
#define CSRS(csr, rs1) WORD $(((csr)<<20) | ((rs1)<<15) | (2<<12) | 0x73)

// CSRS_RET: csrrs T0, csr, rs1 -- set CSR bits and read result into T0 (add256 carry)
#define CSRS_RET(csr, rs1) WORD $(((csr)<<20) | ((rs1)<<15) | (2<<12) | (5<<7) | 0x73)

// CSRWI: csrrwi x0, csr, uimm -- write immediate to CSR (fcall trigger)
#define CSRWI(csr, uimm) WORD $(((csr)<<20) | ((uimm)<<15) | (5<<12) | 0x73)

// CSRR: csrrs T0, csr, x0 -- read CSR value into T0 (fcall output read)
#define CSRR(csr) WORD $(((csr)<<20) | (2<<12) | (5<<7) | 0x73)
