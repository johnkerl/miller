package bifs

import (
	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

// ================================================================
// Bitwise NOT

func bitwise_not_i_i(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(^input1.AcquireIntValue())
}

var bitwise_not_dispositions = [mlrval.MT_DIM]UnaryFunc{
	/*INT    */ bitwise_not_i_i,
	/*FLOAT  */ _erro1,
	/*BOOL   */ _erro1,
	/*VOID   */ _void1,
	/*STRING */ _erro1,
	/*ARRAY  */ _absn1,
	/*MAP    */ _absn1,
	/*FUNC   */ _erro1,
	/*ERROR  */ _erro1,
	/*NULL   */ _null1,
	/*ABSENT */ _absn1,
}

func BIF_bitwise_not(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return bitwise_not_dispositions[input1.Type()](input1)
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

func bitcount_i_i(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	a := uint64(input1.AcquireIntValue())
	a = (a & _m01) + ((a >> 1) & _m01)
	a = (a & _m02) + ((a >> 2) & _m02)
	a = (a & _m04) + ((a >> 4) & _m04)
	a = (a & _m08) + ((a >> 8) & _m08)
	a = (a & _m16) + ((a >> 16) & _m16)
	a = (a & _m32) + ((a >> 32) & _m32)
	return mlrval.FromInt(int64(a))
}

var bitcount_dispositions = [mlrval.MT_DIM]UnaryFunc{
	/*INT    */ bitcount_i_i,
	/*FLOAT  */ _erro1,
	/*BOOL   */ _erro1,
	/*VOID   */ _void1,
	/*STRING */ _erro1,
	/*ARRAY  */ _absn1,
	/*MAP    */ _absn1,
	/*FUNC   */ _erro1,
	/*ERROR  */ _erro1,
	/*NULL   */ _zero1,
	/*ABSENT */ _absn1,
}

func BIF_bitcount(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return bitcount_dispositions[input1.Type()](input1)
}

// ================================================================
// Bitwise AND

func bitwise_and_i_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(input1.AcquireIntValue() & input2.AcquireIntValue())
}

var bitwise_and_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT               FLOAT  BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {bitwise_and_i_ii, _erro, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _1___},
	/*FLOAT  */ {_erro, _erro, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*VOID   */ {_void, _void, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _erro, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _erro, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn},
	/*ABSENT */ {_2___, _erro, _erro, _absn, _erro, _absn, _absn, _erro, _erro, _absn, _absn},
}

func BIF_bitwise_and(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return bitwise_and_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ----------------------------------------------------------------
// Bitwise OR

func bitwise_or_i_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(input1.AcquireIntValue() | input2.AcquireIntValue())
}

var bitwise_or_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT              FLOAT  BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {bitwise_or_i_ii, _erro, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _1___},
	/*FLOAT  */ {_erro, _erro, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*VOID   */ {_void, _void, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _erro, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _erro, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn},
	/*ABSENT */ {_2___, _erro, _erro, _absn, _erro, _absn, _absn, _erro, _erro, _absn, _absn},
}

func BIF_bitwise_or(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return bitwise_or_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ----------------------------------------------------------------
// Bitwise XOR

func bitwise_xor_i_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(input1.AcquireIntValue() ^ input2.AcquireIntValue())
}

var bitwise_xor_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT               FLOAT  BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {bitwise_xor_i_ii, _erro, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _1___},
	/*FLOAT  */ {_erro, _erro, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*VOID   */ {_void, _void, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn},
	/*ABSENT */ {_2___, _erro, _erro, _absn, _erro, _absn, _absn, _erro, _erro, _absn, _absn},
}

func BIF_bitwise_xor(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return bitwise_xor_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ----------------------------------------------------------------
// Left shift

func lsh_i_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(input1.AcquireIntValue() << uint64(input2.AcquireIntValue()))
}

var left_shift_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT       FLOAT  BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {lsh_i_ii, _erro, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _1___},
	/*FLOAT  */ {_erro, _erro, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*VOID   */ {_void, _void, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _erro, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _erro, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn},
	/*ABSENT */ {_2___, _erro, _erro, _absn, _erro, _absn, _absn, _erro, _erro, _absn, _absn},
}

func BIF_left_shift(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return left_shift_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ----------------------------------------------------------------
// Signed right shift

func srsh_i_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(input1.AcquireIntValue() >> uint64(input2.AcquireIntValue()))
}

var signed_right_shift_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT        FLOAT  BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {srsh_i_ii, _erro, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _1___},
	/*FLOAT  */ {_erro, _erro, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*VOID   */ {_void, _void, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _erro, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _erro, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn},
	/*ABSENT */ {_2___, _erro, _erro, _absn, _erro, _absn, _absn, _erro, _erro, _absn, _absn},
}

func BIF_signed_right_shift(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return signed_right_shift_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ----------------------------------------------------------------
// Unsigned right shift

func ursh_i_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	var ua uint64 = uint64(input1.AcquireIntValue())
	var ub uint64 = uint64(input2.AcquireIntValue())
	var uc = ua >> ub
	return mlrval.FromInt(int64(uc))
}

var unsigned_right_shift_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT        FLOAT  BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {ursh_i_ii, _erro, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _1___},
	/*FLOAT  */ {_erro, _erro, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*VOID   */ {_void, _void, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _erro, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _erro, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn},
	/*ABSENT */ {_2___, _erro, _erro, _absn, _erro, _absn, _absn, _erro, _erro, _absn, _absn},
}

func BIF_unsigned_right_shift(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return unsigned_right_shift_dispositions[input1.Type()][input2.Type()](input1, input2)
}
