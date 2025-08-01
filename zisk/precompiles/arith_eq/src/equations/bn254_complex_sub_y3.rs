// code generated
//
// equation: y1-y2-y3+p*q2-p*offset
//
// p: 0x30644E72E131A029B85045B68181585D97816A916871CA8D3C208C16D87CFD47
// offset: 0x8
// (p*offset): 0x183227397098D014DC2822DB40C0AC2ECBC0B548B438E5469E10460B6C3E7EA38
//
// chunks:16
// chunk_bits:16
// terms_by_clock: 2

pub struct Bn254ComplexSubY3 {}

impl Bn254ComplexSubY3 {
    #[allow(clippy::too_many_arguments)]
    pub fn calculate(
        icol: u8,
        y1: &[i64; 16],
        y2: &[i64; 16],
        y3: &[i64; 16],
        q2: &[i64; 16],
    ) -> i64 {
        match icol {
            0 => y1[0] - y2[0] - y3[0] + 0xFD47 * q2[0] - 0xEA38,
            1 => y1[1] - y2[1] - y3[1] + 0xD87C * q2[0] + 0xFD47 * q2[1] - 0xC3E7,
            2 => y1[2] - y2[2] - y3[2] + 0x8C16 * q2[0] + 0xD87C * q2[1] + 0xFD47 * q2[2] - 0x60B6,
            3 => {
                y1[3] - y2[3] - y3[3]
                    + 0x3C20 * q2[0]
                    + 0x8C16 * q2[1]
                    + 0xD87C * q2[2]
                    + 0xFD47 * q2[3]
                    - 0xE104
            }
            4 => {
                y1[4] - y2[4] - y3[4]
                    + 0xCA8D * q2[0]
                    + 0x3C20 * q2[1]
                    + 0x8C16 * q2[2]
                    + 0xD87C * q2[3]
                    + 0xFD47 * q2[4]
                    - 0x5469
            }
            5 => {
                y1[5] - y2[5] - y3[5]
                    + 0x6871 * q2[0]
                    + 0xCA8D * q2[1]
                    + 0x3C20 * q2[2]
                    + 0x8C16 * q2[3]
                    + 0xD87C * q2[4]
                    + 0xFD47 * q2[5]
                    - 0x438E
            }
            6 => {
                y1[6] - y2[6] - y3[6]
                    + 0x6A91 * q2[0]
                    + 0x6871 * q2[1]
                    + 0xCA8D * q2[2]
                    + 0x3C20 * q2[3]
                    + 0x8C16 * q2[4]
                    + 0xD87C * q2[5]
                    + 0xFD47 * q2[6]
                    - 0x548B
            }
            7 => {
                y1[7] - y2[7] - y3[7]
                    + 0x9781 * q2[0]
                    + 0x6A91 * q2[1]
                    + 0x6871 * q2[2]
                    + 0xCA8D * q2[3]
                    + 0x3C20 * q2[4]
                    + 0x8C16 * q2[5]
                    + 0xD87C * q2[6]
                    + 0xFD47 * q2[7]
                    - 0xBC0B
            }
            8 => {
                y1[8] - y2[8] - y3[8]
                    + 0x585D * q2[0]
                    + 0x9781 * q2[1]
                    + 0x6A91 * q2[2]
                    + 0x6871 * q2[3]
                    + 0xCA8D * q2[4]
                    + 0x3C20 * q2[5]
                    + 0x8C16 * q2[6]
                    + 0xD87C * q2[7]
                    + 0xFD47 * q2[8]
                    - 0xC2EC
            }
            9 => {
                y1[9] - y2[9] - y3[9]
                    + 0x8181 * q2[0]
                    + 0x585D * q2[1]
                    + 0x9781 * q2[2]
                    + 0x6A91 * q2[3]
                    + 0x6871 * q2[4]
                    + 0xCA8D * q2[5]
                    + 0x3C20 * q2[6]
                    + 0x8C16 * q2[7]
                    + 0xD87C * q2[8]
                    + 0xFD47 * q2[9]
                    - 0xC0A
            }
            10 => {
                y1[10] - y2[10] - y3[10]
                    + 0x45B6 * q2[0]
                    + 0x8181 * q2[1]
                    + 0x585D * q2[2]
                    + 0x9781 * q2[3]
                    + 0x6A91 * q2[4]
                    + 0x6871 * q2[5]
                    + 0xCA8D * q2[6]
                    + 0x3C20 * q2[7]
                    + 0x8C16 * q2[8]
                    + 0xD87C * q2[9]
                    + 0xFD47 * q2[10]
                    - 0x2DB4
            }
            11 => {
                y1[11] - y2[11] - y3[11]
                    + 0xB850 * q2[0]
                    + 0x45B6 * q2[1]
                    + 0x8181 * q2[2]
                    + 0x585D * q2[3]
                    + 0x9781 * q2[4]
                    + 0x6A91 * q2[5]
                    + 0x6871 * q2[6]
                    + 0xCA8D * q2[7]
                    + 0x3C20 * q2[8]
                    + 0x8C16 * q2[9]
                    + 0xD87C * q2[10]
                    + 0xFD47 * q2[11]
                    - 0xC282
            }
            12 => {
                y1[12] - y2[12] - y3[12]
                    + 0xA029 * q2[0]
                    + 0xB850 * q2[1]
                    + 0x45B6 * q2[2]
                    + 0x8181 * q2[3]
                    + 0x585D * q2[4]
                    + 0x9781 * q2[5]
                    + 0x6A91 * q2[6]
                    + 0x6871 * q2[7]
                    + 0xCA8D * q2[8]
                    + 0x3C20 * q2[9]
                    + 0x8C16 * q2[10]
                    + 0xD87C * q2[11]
                    + 0xFD47 * q2[12]
                    - 0x14D
            }
            13 => {
                y1[13] - y2[13] - y3[13]
                    + 0xE131 * q2[0]
                    + 0xA029 * q2[1]
                    + 0xB850 * q2[2]
                    + 0x45B6 * q2[3]
                    + 0x8181 * q2[4]
                    + 0x585D * q2[5]
                    + 0x9781 * q2[6]
                    + 0x6A91 * q2[7]
                    + 0x6871 * q2[8]
                    + 0xCA8D * q2[9]
                    + 0x3C20 * q2[10]
                    + 0x8C16 * q2[11]
                    + 0xD87C * q2[12]
                    + 0xFD47 * q2[13]
                    - 0x98D
            }
            14 => {
                y1[14] - y2[14] - y3[14]
                    + 0x4E72 * q2[0]
                    + 0xE131 * q2[1]
                    + 0xA029 * q2[2]
                    + 0xB850 * q2[3]
                    + 0x45B6 * q2[4]
                    + 0x8181 * q2[5]
                    + 0x585D * q2[6]
                    + 0x9781 * q2[7]
                    + 0x6A91 * q2[8]
                    + 0x6871 * q2[9]
                    + 0xCA8D * q2[10]
                    + 0x3C20 * q2[11]
                    + 0x8C16 * q2[12]
                    + 0xD87C * q2[13]
                    + 0xFD47 * q2[14]
                    - 0x7397
            }
            15 => {
                y1[15] - y2[15] - y3[15]
                    + 0x3064 * q2[0]
                    + 0x4E72 * q2[1]
                    + 0xE131 * q2[2]
                    + 0xA029 * q2[3]
                    + 0xB850 * q2[4]
                    + 0x45B6 * q2[5]
                    + 0x8181 * q2[6]
                    + 0x585D * q2[7]
                    + 0x9781 * q2[8]
                    + 0x6A91 * q2[9]
                    + 0x6871 * q2[10]
                    + 0xCA8D * q2[11]
                    + 0x3C20 * q2[12]
                    + 0x8C16 * q2[13]
                    + 0xD87C * q2[14]
                    + 0xFD47 * q2[15]
                    - 0x8322
            }
            16 => {
                0x3064 * q2[1]
                    + 0x4E72 * q2[2]
                    + 0xE131 * q2[3]
                    + 0xA029 * q2[4]
                    + 0xB850 * q2[5]
                    + 0x45B6 * q2[6]
                    + 0x8181 * q2[7]
                    + 0x585D * q2[8]
                    + 0x9781 * q2[9]
                    + 0x6A91 * q2[10]
                    + 0x6871 * q2[11]
                    + 0xCA8D * q2[12]
                    + 0x3C20 * q2[13]
                    + 0x8C16 * q2[14]
                    + 0xD87C * q2[15]
                    - 0x1
            }
            17 => {
                0x3064 * q2[2]
                    + 0x4E72 * q2[3]
                    + 0xE131 * q2[4]
                    + 0xA029 * q2[5]
                    + 0xB850 * q2[6]
                    + 0x45B6 * q2[7]
                    + 0x8181 * q2[8]
                    + 0x585D * q2[9]
                    + 0x9781 * q2[10]
                    + 0x6A91 * q2[11]
                    + 0x6871 * q2[12]
                    + 0xCA8D * q2[13]
                    + 0x3C20 * q2[14]
                    + 0x8C16 * q2[15]
            }
            18 => {
                0x3064 * q2[3]
                    + 0x4E72 * q2[4]
                    + 0xE131 * q2[5]
                    + 0xA029 * q2[6]
                    + 0xB850 * q2[7]
                    + 0x45B6 * q2[8]
                    + 0x8181 * q2[9]
                    + 0x585D * q2[10]
                    + 0x9781 * q2[11]
                    + 0x6A91 * q2[12]
                    + 0x6871 * q2[13]
                    + 0xCA8D * q2[14]
                    + 0x3C20 * q2[15]
            }
            19 => {
                0x3064 * q2[4]
                    + 0x4E72 * q2[5]
                    + 0xE131 * q2[6]
                    + 0xA029 * q2[7]
                    + 0xB850 * q2[8]
                    + 0x45B6 * q2[9]
                    + 0x8181 * q2[10]
                    + 0x585D * q2[11]
                    + 0x9781 * q2[12]
                    + 0x6A91 * q2[13]
                    + 0x6871 * q2[14]
                    + 0xCA8D * q2[15]
            }
            20 => {
                0x3064 * q2[5]
                    + 0x4E72 * q2[6]
                    + 0xE131 * q2[7]
                    + 0xA029 * q2[8]
                    + 0xB850 * q2[9]
                    + 0x45B6 * q2[10]
                    + 0x8181 * q2[11]
                    + 0x585D * q2[12]
                    + 0x9781 * q2[13]
                    + 0x6A91 * q2[14]
                    + 0x6871 * q2[15]
            }
            21 => {
                0x3064 * q2[6]
                    + 0x4E72 * q2[7]
                    + 0xE131 * q2[8]
                    + 0xA029 * q2[9]
                    + 0xB850 * q2[10]
                    + 0x45B6 * q2[11]
                    + 0x8181 * q2[12]
                    + 0x585D * q2[13]
                    + 0x9781 * q2[14]
                    + 0x6A91 * q2[15]
            }
            22 => {
                0x3064 * q2[7]
                    + 0x4E72 * q2[8]
                    + 0xE131 * q2[9]
                    + 0xA029 * q2[10]
                    + 0xB850 * q2[11]
                    + 0x45B6 * q2[12]
                    + 0x8181 * q2[13]
                    + 0x585D * q2[14]
                    + 0x9781 * q2[15]
            }
            23 => {
                0x3064 * q2[8]
                    + 0x4E72 * q2[9]
                    + 0xE131 * q2[10]
                    + 0xA029 * q2[11]
                    + 0xB850 * q2[12]
                    + 0x45B6 * q2[13]
                    + 0x8181 * q2[14]
                    + 0x585D * q2[15]
            }
            24 => {
                0x3064 * q2[9]
                    + 0x4E72 * q2[10]
                    + 0xE131 * q2[11]
                    + 0xA029 * q2[12]
                    + 0xB850 * q2[13]
                    + 0x45B6 * q2[14]
                    + 0x8181 * q2[15]
            }
            25 => {
                0x3064 * q2[10]
                    + 0x4E72 * q2[11]
                    + 0xE131 * q2[12]
                    + 0xA029 * q2[13]
                    + 0xB850 * q2[14]
                    + 0x45B6 * q2[15]
            }
            26 => {
                0x3064 * q2[11]
                    + 0x4E72 * q2[12]
                    + 0xE131 * q2[13]
                    + 0xA029 * q2[14]
                    + 0xB850 * q2[15]
            }
            27 => 0x3064 * q2[12] + 0x4E72 * q2[13] + 0xE131 * q2[14] + 0xA029 * q2[15],
            28 => 0x3064 * q2[13] + 0x4E72 * q2[14] + 0xE131 * q2[15],
            29 => 0x3064 * q2[14] + 0x4E72 * q2[15],
            30 => 0x3064 * q2[15],
            _ => 0,
        }
    }
}
