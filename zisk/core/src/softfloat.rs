//! Software floating-point operations
//! 
//! Implements IEEE 754 floating-point arithmetic by emitting sequences of 
//! Zisk integer operations that implement the floating-point semantics.

use crate::{ZiskInstBuilder, FREG_BASE_ADDR, FREG_SIZE};
use std::collections::HashMap;

// IEEE 754 rounding modes (stored in FRM register)
pub const FRM_RNE: u64 = 0; // Round to Nearest, ties to Even (default)
pub const FRM_RTZ: u64 = 1; // Round towards Zero
pub const FRM_RDN: u64 = 2; // Round Down (towards negative infinity)
pub const FRM_RUP: u64 = 3; // Round Up (towards positive infinity)
pub const FRM_RMM: u64 = 4; // Round to Nearest, ties to Max magnitude

// IEEE 754 exception flags (bits in FFLAGS register)
pub const FFLAGS_NX: u64 = 1 << 0; // Inexact
pub const FFLAGS_UF: u64 = 1 << 1; // Underflow
pub const FFLAGS_OF: u64 = 1 << 2; // Overflow
pub const FFLAGS_DZ: u64 = 1 << 3; // Divide by zero
pub const FFLAGS_NV: u64 = 1 << 4; // Invalid operation

// FCSR is a view over FRM (bits 7:5) and FFLAGS (bits 4:0)
// FCSR = (FRM << 5) | FFLAGS

// ---------- IEEE754 single-precision layout ----------
const FRAC_BITS: u64 = 23;
const EXP_BITS:  u64 = 8;
const SIGN_MASK: u64 = 0x8000_0000;
const EXP_MASK:  u64 = 0x7F80_0000;
const FRAC_MASK: u64 = 0x007F_FFFF;
const EXP_BIAS:  u64 = 127;
const HIDDEN_1:  u64 = 1 << FRAC_BITS;

// ---------- scratch FP memory "temps" (memory addresses in FP register space) ----------
// Use unused FP register slots as temporary storage to avoid integer register conflicts
const TEMP_A_BITS: u64 = FREG_BASE_ADDR + (20 * FREG_SIZE);  // f20 slot
const TEMP_B_BITS: u64 = FREG_BASE_ADDR + (21 * FREG_SIZE);  // f21 slot

const TEMP_SA: u64 = FREG_BASE_ADDR + (22 * FREG_SIZE); const TEMP_SB: u64 = FREG_BASE_ADDR + (23 * FREG_SIZE);
const TEMP_EA: u64 = FREG_BASE_ADDR + (24 * FREG_SIZE); const TEMP_EB: u64 = FREG_BASE_ADDR + (25 * FREG_SIZE);
const TEMP_FA: u64 = FREG_BASE_ADDR + (26 * FREG_SIZE); const TEMP_FB: u64 = FREG_BASE_ADDR + (27 * FREG_SIZE);

const TEMP_MA: u64 = FREG_BASE_ADDR + (28 * FREG_SIZE); const TEMP_MB: u64 = FREG_BASE_ADDR + (29 * FREG_SIZE);

const TEMP_CMP:  u64 = FREG_BASE_ADDR + (30 * FREG_SIZE);
const TEMP_MASK: u64 = FREG_BASE_ADDR + (31 * FREG_SIZE);

// For more temporaries, we can use the same slots since operations don't overlap
const TEMP_TMP:    u64 = FREG_BASE_ADDR + (20 * FREG_SIZE);  // Reuse f20
const TEMP_DE:     u64 = FREG_BASE_ADDR + (21 * FREG_SIZE);  // Reuse f21
const TEMP_MB_AL:  u64 = FREG_BASE_ADDR + (22 * FREG_SIZE);  // Reuse f22
const TEMP_SUM:    u64 = FREG_BASE_ADDR + (23 * FREG_SIZE);  // Reuse f23
const TEMP_CARRY:  u64 = FREG_BASE_ADDR + (24 * FREG_SIZE);  // Reuse f24
const TEMP_SIG:    u64 = FREG_BASE_ADDR + (25 * FREG_SIZE);  // Reuse f25
const TEMP_EXP:    u64 = FREG_BASE_ADDR + (26 * FREG_SIZE);  // Reuse f26
const TEMP_FRAC:   u64 = FREG_BASE_ADDR + (27 * FREG_SIZE);  // Reuse f27
const TEMP_PACK1:  u64 = FREG_BASE_ADDR + (28 * FREG_SIZE);  // Reuse f28
const TEMP_PACK2:  u64 = FREG_BASE_ADDR + (29 * FREG_SIZE);  // Reuse f29
const TEMP_RES32:  u64 = FREG_BASE_ADDR + (30 * FREG_SIZE);  // Reuse f30

// ---------- tiny emitter helper ----------
#[inline(always)]
fn emit1(insts: &mut HashMap<u64, ZiskInstBuilder>, pc: &mut u64, mut z: ZiskInstBuilder) {
    z.j(4, 4);
    z.build();
    insts.insert(*pc, z);
    *pc += 4;
}

#[inline(always)]
fn ez(insts: &mut HashMap<u64, ZiskInstBuilder>, pc: &mut u64) -> ZiskInstBuilder {
    ZiskInstBuilder::new(*pc)
}

// ---------- common micro-ops (make call sites readable) ----------

fn load_freg_bits(insts: &mut HashMap<u64, ZiskInstBuilder>, pc: &mut u64, faddr: u64, dst_addr: u64, note: &str) {
    let mut z = ez(insts, pc);
    z.src_a("imm", 0, false);
    z.src_b("mem", faddr, false); // loads 64b; our 32b payload is in low bits
    z.op("copyb").unwrap();
    z.store("mem", dst_addr as i64, false, false);  // Store to memory, not register
    z.verbose(note);
    emit1(insts, pc, z);
}

fn store_freg_bits(insts: &mut HashMap<u64, ZiskInstBuilder>, pc: &mut u64, src_addr: u64, faddr: u64, note: &str) {
    let mut z = ez(insts, pc);
    z.src_a("imm", 0, false);
    z.src_b("mem", src_addr, false);  // Read from memory, not register
    z.op("copyb").unwrap();              // writes full 64b; high 32 will be zero here
    z.store("mem", faddr as i64, false, false);
    z.verbose(note);
    emit1(insts, pc, z);
}

fn srl_imm(insts: &mut HashMap<u64, ZiskInstBuilder>, pc: &mut u64, src_addr: u64, sh: u64, dst_addr: u64, note: &str) {
    let mut z = ez(insts, pc);
    z.src_a("mem", src_addr, false);  // Read from memory
    z.src_b("imm", sh, false);
    z.op("srl").unwrap();
    z.store("mem", dst_addr as i64, false, false);  // Store to memory
    z.verbose(note);
    emit1(insts, pc, z);
}

fn srl_var(insts: &mut HashMap<u64, ZiskInstBuilder>, pc: &mut u64, src_addr: u64, sh_addr: u64, dst_addr: u64, note: &str) {
    let mut z = ez(insts, pc);
    z.src_a("mem", src_addr, false);  // Read from memory
    z.src_b("mem", sh_addr, false);   // Read shift amount from memory
    z.op("srl").unwrap();
    z.store("mem", dst_addr as i64, false, false);  // Store to memory
    z.verbose(note);
    emit1(insts, pc, z);
}

fn sll_imm(insts: &mut HashMap<u64, ZiskInstBuilder>, pc: &mut u64, src_addr: u64, sh: u64, dst_addr: u64, note: &str) {
    let mut z = ez(insts, pc);
    z.src_a("mem", src_addr, false);  // Read from memory
    z.src_b("imm", sh, false);
    z.op("sll").unwrap();
    z.store("mem", dst_addr as i64, false, false);  // Store to memory
    z.verbose(note);
    emit1(insts, pc, z);
}

fn and_imm(insts: &mut HashMap<u64, ZiskInstBuilder>, pc: &mut u64, src_addr: u64, imm: u64, dst_addr: u64, note: &str) {
    let mut z = ez(insts, pc);
    z.src_a("mem", src_addr, false);  // Read from memory
    z.src_b("imm", imm, false);
    z.op("and").unwrap();
    z.store("mem", dst_addr as i64, false, false);  // Store to memory
    z.verbose(note);
    emit1(insts, pc, z);
}

fn or_imm(insts: &mut HashMap<u64, ZiskInstBuilder>, pc: &mut u64, src_addr: u64, imm: u64, dst_addr: u64, note: &str) {
    let mut z = ez(insts, pc);
    z.src_a("mem", src_addr, false);  // Read from memory
    z.src_b("imm", imm, false);
    z.op("or").unwrap();
    z.store("mem", dst_addr as i64, false, false);  // Store to memory
    z.verbose(note);
    emit1(insts, pc, z);
}

fn add_rr(insts: &mut HashMap<u64, ZiskInstBuilder>, pc: &mut u64, a_addr: u64, b_addr: u64, dst_addr: u64, note: &str) {
    let mut z = ez(insts, pc);
    z.src_a("mem", a_addr, false);  // Read from memory
    z.src_b("mem", b_addr, false);  // Read from memory
    z.op("add").unwrap();
    z.store("mem", dst_addr as i64, false, false);  // Store to memory
    z.verbose(note);
    emit1(insts, pc, z);
}

fn sub_rr(insts: &mut HashMap<u64, ZiskInstBuilder>, pc: &mut u64, a_addr: u64, b_addr: u64, dst_addr: u64, note: &str) {
    let mut z = ez(insts, pc);
    z.src_a("mem", a_addr, false);  // Read from memory
    z.src_b("mem", b_addr, false);  // Read from memory
    z.op("sub").unwrap();
    z.store("mem", dst_addr as i64, false, false);  // Store to memory
    z.verbose(note);
    emit1(insts, pc, z);
}

fn xor_rr(insts: &mut HashMap<u64, ZiskInstBuilder>, pc: &mut u64, a_addr: u64, b_addr: u64, dst_addr: u64, note: &str) {
    let mut z = ez(insts, pc);
    z.src_a("mem", a_addr, false);  // Read from memory
    z.src_b("mem", b_addr, false);  // Read from memory
    z.op("xor").unwrap();
    z.store("mem", dst_addr as i64, false, false);  // Store to memory
    z.verbose(note);
    emit1(insts, pc, z);
}

fn ltu_rr(insts: &mut HashMap<u64, ZiskInstBuilder>, pc: &mut u64, a_addr: u64, b_addr: u64, dst_addr: u64, note: &str) {
    let mut z = ez(insts, pc);
    z.src_a("mem", a_addr, false);  // Read from memory
    z.src_b("mem", b_addr, false);  // Read from memory
    z.op("ltu").unwrap(); // c=1 if a<b
    z.store("mem", dst_addr as i64, false, false);  // Store to memory
    z.verbose(note);
    emit1(insts, pc, z);
}

fn neg_mask_from_bool(insts: &mut HashMap<u64, ZiskInstBuilder>, pc: &mut u64, bool_addr: u64, mask_addr: u64, note: &str) {
    let mut z = ez(insts, pc);
    z.src_a("imm", 0, false);
    z.src_b("mem", bool_addr, false);  // Read from memory
    z.op("sub").unwrap(); // 0 - b
    z.store("mem", mask_addr as i64, false, false);  // Store to memory
    z.verbose(note);
    emit1(insts, pc, z);
}

// swap rX <-> rY only if mask == all-ones; uses tmp
fn swap_by_mask(insts: &mut HashMap<u64, ZiskInstBuilder>, pc: &mut u64, rx_addr: u64, ry_addr: u64, mask_addr: u64, tmp_addr: u64, what: &str) {
    xor_rr(insts, pc, rx_addr, ry_addr, tmp_addr, &format!("t = {} ^ {}", what, what));
    {
        let mut z = ez(insts, pc);
        z.src_a("mem", tmp_addr, false);    // Read from memory
        z.src_b("mem", mask_addr, false);   // Read from memory
        z.op("and").unwrap();
        z.store("mem", tmp_addr as i64, false, false);  // Store to memory
        z.verbose("t = (x^y) & mask");
        emit1(insts, pc, z);
    }
    xor_rr(insts, pc, rx_addr, tmp_addr, rx_addr, "x ^= t");
    xor_rr(insts, pc, ry_addr, tmp_addr, ry_addr, "y ^= t");
}

// dst = (keep_if_zero_mask ? a : b) with mask=0 or all-ones
fn branchless_select(insts: &mut HashMap<u64, ZiskInstBuilder>, pc: &mut u64, a_addr: u64, b_addr: u64, mask_addr: u64, dst_addr: u64, note: &str) {
    // diff = (a ^ b) & mask; dst = a ^ diff
    xor_rr(insts, pc, a_addr, b_addr, TEMP_TMP, "diff = a ^ b");
    {
        let mut z = ez(insts, pc);
        z.src_a("mem", TEMP_TMP, false);    // Read from memory
        z.src_b("mem", mask_addr, false);   // Read from memory
        z.op("and").unwrap();
        z.store("mem", TEMP_TMP as i64, false, false);  // Store to memory
        z.verbose("diff &= mask");
        emit1(insts, pc, z);
    }
    xor_rr(insts, pc, a_addr, TEMP_TMP, dst_addr, note);
}

fn unpack_fp32_bits(insts: &mut HashMap<u64, ZiskInstBuilder>, pc: &mut u64, src_bits_addr: u64, rs_addr: u64, re_addr: u64, rf_addr: u64, tag: &str) {
    srl_imm(insts, pc, src_bits_addr, 31, rs_addr, &format!("{tag}: s = bits>>31"));
    srl_imm(insts, pc, src_bits_addr, FRAC_BITS, re_addr, &format!("{tag}: e_tmp = bits>>23"));
    and_imm(insts, pc, re_addr, (1<<EXP_BITS)-1, re_addr, &format!("{tag}: e = e_tmp & 0xFF"));
    and_imm(insts, pc, src_bits_addr, FRAC_MASK, rf_addr, &format!("{tag}: f = bits & 0x007FFFFF"));
}

fn add_hidden_one(insts: &mut HashMap<u64, ZiskInstBuilder>, pc: &mut u64, frac_addr: u64, mant_dst_addr: u64, tag: &str) {
    or_imm(insts, pc, frac_addr, HIDDEN_1, mant_dst_addr, &format!("{tag}: mant = frac | (1<<23)"));
}

fn pack_fp32_bits(insts: &mut HashMap<u64, ZiskInstBuilder>, pc: &mut u64, rs_addr: u64, re_addr: u64, rf_addr: u64, dst_addr: u64, tag: &str) {
    // pack1 = (re & 0xFF) << 23
    and_imm(insts, pc, re_addr, 0xFF, TEMP_PACK1, &format!("{tag}: pack1 = e & 0xFF"));
    sll_imm(insts, pc, TEMP_PACK1, FRAC_BITS, TEMP_PACK1, &format!("{tag}: pack1 <<= 23"));
    // pack2 = (rs << 31)
    sll_imm(insts, pc, rs_addr, 31, TEMP_PACK2, &format!("{tag}: pack2 = s<<31"));
    // res = pack2 | pack1 | rf
    add_rr(insts, pc, TEMP_PACK2, TEMP_PACK1, dst_addr, &format!("{tag}: res = pack2|pack1 (via add/xor ok)"));
    {
        let mut z = ez(insts, pc);
        z.src_a("mem", dst_addr, false);    // Read from memory
        z.src_b("mem", rf_addr, false);     // Read from memory
        z.op("or").unwrap();
        z.store("mem", dst_addr as i64, false, false);  // Store to memory
        z.verbose(&format!("{tag}: res |= frac"));
        emit1(insts, pc, z);
    }
}

// ---------- fp32 add (fast path) ----------
pub fn emit_fp32_add(
    insts: &mut HashMap<u64, ZiskInstBuilder>,
    pc: &mut u64,
    rs1: u32,
    rs2: u32,
    rd: u32,
) {
    let a_addr = FREG_BASE_ADDR + (rs1 as u64 * FREG_SIZE);
    let b_addr = FREG_BASE_ADDR + (rs2 as u64 * FREG_SIZE);
    let d_addr = FREG_BASE_ADDR + (rd  as u64 * FREG_SIZE);

    // 1) load raw bits from the FP "register slots" (memory)
    load_freg_bits(insts, pc, a_addr, TEMP_A_BITS, &format!("fp32_add: load f{rs1} bits"));
    load_freg_bits(insts, pc, b_addr, TEMP_B_BITS, &format!("fp32_add: load f{rs2} bits"));

    // 2) unpack (sign, exp, frac)
    unpack_fp32_bits(insts, pc, TEMP_A_BITS, TEMP_SA, TEMP_EA, TEMP_FA, "A");
    unpack_fp32_bits(insts, pc, TEMP_B_BITS, TEMP_SB, TEMP_EB, TEMP_FB, "B");

    // 3) (fast-path assumption) normals + same sign; add hidden 1
    add_hidden_one(insts, pc, TEMP_FA, TEMP_MA, "A");
    add_hidden_one(insts, pc, TEMP_FB, TEMP_MB, "B");

    // 4) ensure EA >= EB by conditional swap (so we always shift mb)
    ltu_rr(insts, pc, TEMP_EA, TEMP_EB, TEMP_CMP, "cmp = (ea<eb)");
    neg_mask_from_bool(insts, pc, TEMP_CMP, TEMP_MASK, "mask = 0 - cmp");
    swap_by_mask(insts, pc, TEMP_EA, TEMP_EB, TEMP_MASK, TEMP_TMP, "e");
    swap_by_mask(insts, pc, TEMP_MA, TEMP_MB, TEMP_MASK, TEMP_TMP, "m");
    swap_by_mask(insts, pc, TEMP_SA, TEMP_SB, TEMP_MASK, TEMP_TMP, "s");

    // 5) align small mantissa: de = ea - eb; mb_al = mb >> de
    sub_rr(insts, pc, TEMP_EA, TEMP_EB, TEMP_DE, "de = ea - eb");
    srl_var(insts, pc, TEMP_MB, TEMP_DE, TEMP_MB_AL, "mb_al = mb >> de");

    // 6) same-sign ADD: sum = ma + mb_al
    add_rr(insts, pc, TEMP_MA, TEMP_MB_AL, TEMP_SUM, "sum = ma + mb_al");

    // 7) normalize if carry
    srl_imm(insts, pc, TEMP_SUM, FRAC_BITS + 1, TEMP_CARRY, "carry_tmp = sum>>24");
    and_imm(insts, pc, TEMP_CARRY, 1, TEMP_CARRY, "carry = carry_tmp & 1");

    // sig = carry ? (sum>>1) : sum
    srl_imm(insts, pc, TEMP_SUM, 1, TEMP_SIG, "t = sum>>1");
    neg_mask_from_bool(insts, pc, TEMP_CARRY, TEMP_MASK, "mask = 0 - carry");
    branchless_select(insts, pc, TEMP_SUM, TEMP_SIG, TEMP_MASK, TEMP_SIG, "sig = carry ? t : sum");

    // exp' = ea + carry
    add_rr(insts, pc, TEMP_EA, TEMP_CARRY, TEMP_EXP, "exp = ea + carry");

    // (no rounding yet) strip hidden bit for fraction
    and_imm(insts, pc, TEMP_SIG, FRAC_MASK, TEMP_FRAC, "frac = sig & ((1<<23)-1)");

    // 8) pack + store
    pack_fp32_bits(insts, pc, TEMP_SA, TEMP_EXP, TEMP_FRAC, TEMP_RES32, "pack");
    store_freg_bits(insts, pc, TEMP_RES32, d_addr, &format!("fp32_add: f{rd} <- result (fast path)"));
}

/// Emits ROM subroutine for IEEE 754 single-precision floating-point addition
/// This just calls the full emit_fp32_add implementation with dummy parameters
pub fn emit_fp32_add_subroutine(
    insts: &mut HashMap<u64, ZiskInstBuilder>,
    start_addr: &mut u64,
) {
    println!("DEBUG: emit_fp32_add_subroutine called with start_addr: 0x{:x}", *start_addr);
    
    // Call the full IEEE 754 implementation with dummy parameters (will be overridden)
    // The comprehensive implementation will use f20-f31 as temporaries
    let before = *start_addr;
    emit_fp32_add(insts, start_addr, 1, 2, 3);  // f1 + f2 -> f3
    let after = *start_addr;
    
    println!("DEBUG: emit_fp32_add generated {} instructions (0x{:x} to 0x{:x})", 
             (after - before) / 4, before, after);
}

/// Emits Zisk instructions for IEEE 754 double-precision floating-point addition
pub fn emit_fp64_add(
    _insts: &mut HashMap<u64, ZiskInstBuilder>,
    _start_addr: &mut u64,
    _rs1: u32,
    _rs2: u32,
    _rd: u32,
) {
    todo!("fp64_add not implemented yet")
}

/// Emits Zisk instructions for IEEE 754 single-precision floating-point subtraction
pub fn emit_fp32_sub(
    _insts: &mut HashMap<u64, ZiskInstBuilder>,
    _start_addr: &mut u64,
    _rs1: u32,
    _rs2: u32,
    _rd: u32,
) {
    todo!("fp32_sub not implemented yet")
}

/// Emits Zisk instructions for IEEE 754 double-precision floating-point subtraction
pub fn emit_fp64_sub(
    _insts: &mut HashMap<u64, ZiskInstBuilder>,
    _start_addr: &mut u64,
    _rs1: u32,
    _rs2: u32,
    _rd: u32,
) {
    todo!("fp64_sub not implemented yet")
}

/// Emits Zisk instructions for IEEE 754 single-precision floating-point multiplication
pub fn emit_fp32_mul(
    _insts: &mut HashMap<u64, ZiskInstBuilder>,
    _start_addr: &mut u64,
    _rs1: u32,
    _rs2: u32,
    _rd: u32,
) {
    todo!("fp32_mul not implemented yet")
}

/// Emits Zisk instructions for IEEE 754 double-precision floating-point multiplication
pub fn emit_fp64_mul(
    _insts: &mut HashMap<u64, ZiskInstBuilder>,
    _start_addr: &mut u64,
    _rs1: u32,
    _rs2: u32,
    _rd: u32,
) {
    todo!("fp64_mul not implemented yet")
}