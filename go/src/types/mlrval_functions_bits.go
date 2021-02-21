package types

// ================================================================
// Bitwise NOT

func bitwise_not_i_i(output, input1 *Mlrval) {
	output.SetFromInt(^input1.intval)
}

var bitwise_not_dispositions = [MT_DIM]UnaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*VOID   */ _void1,
	/*STRING */ _erro1,
	/*INT    */ bitwise_not_i_i,
	/*FLOAT  */ _erro1,
	/*BOOL   */ _erro1,
	/*ARRAY  */ _absn1,
	/*MAP    */ _absn1,
}

func MlrvalBitwiseNOT(output, input1 *Mlrval) {
	bitwise_not_dispositions[input1.mvtype](output, input1)
}

// ================================================================
// Bit-count
// https://en.wikipedia.org/wiki/Hamming_weight

const _m01 uint64 = 0x5555555555555555
const _m02 uint64 = 0x3333333333333333
const _m04 uint64 = 0x0f0f0f0f0f0f0f0f
const _m08 uint64 = 0x00ff00ff00ff00ff
const _m16 uint64 = 0x0000ffff0000ffff
const _m32 uint64 = 0x00000000ffffffff

func bitcount_i_i(output, input1 *Mlrval) {
	a := uint64(input1.intval)
	a = (a & _m01) + ((a >> 1) & _m01)
	a = (a & _m02) + ((a >> 2) & _m02)
	a = (a & _m04) + ((a >> 4) & _m04)
	a = (a & _m08) + ((a >> 8) & _m08)
	a = (a & _m16) + ((a >> 16) & _m16)
	a = (a & _m32) + ((a >> 32) & _m32)
	output.SetFromInt(int(a))
}

var bitcount_dispositions = [MT_DIM]UnaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*VOID   */ _void1,
	/*STRING */ _erro1,
	/*INT    */ bitcount_i_i,
	/*FLOAT  */ _erro1,
	/*BOOL   */ _erro1,
	/*ARRAY  */ _absn1,
	/*MAP    */ _absn1,
}

func MlrvalBitCount(output, input1 *Mlrval) {
	bitcount_dispositions[input1.mvtype](output, input1)
}

// ================================================================
// Bitwise AND

func bitwise_and_i_ii(output, input1, input2 *Mlrval) {
	output.SetFromInt(input1.intval & input2.intval)
}

var bitwise_and_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _erro, _erro, _absn, _absn},
	/*VOID   */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, bitwise_and_i_ii, _erro, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _erro, _void, _erro, _erro, _erro, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalBitwiseAND(output, input1, input2 *Mlrval) {
	bitwise_and_dispositions[input1.mvtype][input2.mvtype](output, input1, input2)
}

// ----------------------------------------------------------------
// Bitwise OR

func bitwise_or_i_ii(output, input1, input2 *Mlrval) {
	output.SetFromInt(input1.intval | input2.intval)
}

var bitwise_or_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _erro, _erro, _absn, _absn},
	/*VOID   */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, bitwise_or_i_ii, _erro, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _erro, _void, _erro, _erro, _erro, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalBitwiseOR(output, input1, input2 *Mlrval) {
	bitwise_or_dispositions[input1.mvtype][input2.mvtype](output, input1, input2)
}

// ----------------------------------------------------------------
// Bitwise XOR

func bitwise_xor_i_ii(output, input1, input2 *Mlrval) {
	output.SetFromInt(input1.intval ^ input2.intval)
}

var bitwise_xor_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _erro, _erro, _absn, _absn},
	/*VOID   */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, bitwise_xor_i_ii, _erro, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _erro, _void, _erro, _erro, _erro, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalBitwiseXOR(output, input1, input2 *Mlrval) {
	bitwise_xor_dispositions[input1.mvtype][input2.mvtype](output, input1, input2)
}

// ----------------------------------------------------------------
// Left shift

func lsh_i_ii(output, input1, input2 *Mlrval) {
	output.SetFromInt(input1.intval << uint64(input2.intval))
}

var left_shift_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _erro, _erro, _absn, _absn},
	/*VOID   */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, lsh_i_ii, _erro, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _erro, _void, _erro, _erro, _erro, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalLeftShift(output, input1, input2 *Mlrval) {
	left_shift_dispositions[input1.mvtype][input2.mvtype](output, input1, input2)
}

// ----------------------------------------------------------------
// Signed right shift

func srsh_i_ii(output, input1, input2 *Mlrval) {
	output.SetFromInt(input1.intval >> uint64(input2.intval))
}

var signed_right_shift_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _erro, _erro, _absn, _absn},
	/*VOID   */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, srsh_i_ii, _erro, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _erro, _void, _erro, _erro, _erro, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalSignedRightShift(output, input1, input2 *Mlrval) {
	signed_right_shift_dispositions[input1.mvtype][input2.mvtype](output, input1, input2)
}

// ----------------------------------------------------------------
// Unsigned right shift

func ursh_i_ii(output, input1, input2 *Mlrval) {
	var ua uint64 = uint64(input1.intval)
	var ub uint64 = uint64(input2.intval)
	var uc = ua >> ub
	output.SetFromInt(int(uc))
}

var unsigned_right_shift_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _erro, _erro, _absn, _absn},
	/*VOID   */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, ursh_i_ii, _erro, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _erro, _void, _erro, _erro, _erro, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalUnsignedRightShift(output, input1, input2 *Mlrval) {
	unsigned_right_shift_dispositions[input1.mvtype][input2.mvtype](output, input1, input2)
}
