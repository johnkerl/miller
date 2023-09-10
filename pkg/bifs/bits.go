package bifs

import (
	"github.com/johnkerl/miller/pkg/mlrval"
)

// ================================================================
// Bitwise NOT

func bitwise_not_i_i(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(^input1.AcquireIntValue())
}

func bitwise_not_te(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromTypeErrorUnary("~", input1)
}

var bitwise_not_dispositions = [mlrval.MT_DIM]UnaryFunc{
	/*INT    */ bitwise_not_i_i,
	/*FLOAT  */ bitwise_not_te,
	/*BOOL   */ bitwise_not_te,
	/*VOID   */ _void1,
	/*STRING */ bitwise_not_te,
	/*ARRAY  */ _absn1,
	/*MAP    */ _absn1,
	/*FUNC   */ bitwise_not_te,
	/*ERROR  */ bitwise_not_te,
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

func bitcount_te(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromTypeErrorUnary("bitcount", input1)
}

var bitcount_dispositions = [mlrval.MT_DIM]UnaryFunc{
	/*INT    */ bitcount_i_i,
	/*FLOAT  */ bitcount_te,
	/*BOOL   */ bitcount_te,
	/*VOID   */ _void1,
	/*STRING */ bitcount_te,
	/*ARRAY  */ _absn1,
	/*MAP    */ _absn1,
	/*FUNC   */ bitcount_te,
	/*ERROR  */ bitcount_te,
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

func bwandte(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromTypeErrorBinary("&", input1, input2)
}

var bitwise_and_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT               FLOAT  BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {bitwise_and_i_ii, bwandte, bwandte, _void, bwandte, _absn, _absn, bwandte, bwandte, bwandte, _1___},
	/*FLOAT  */ {bwandte, bwandte, bwandte, _void, bwandte, _absn, _absn, bwandte, bwandte, bwandte, bwandte},
	/*BOOL   */ {bwandte, bwandte, bwandte, bwandte, bwandte, _absn, _absn, bwandte, bwandte, bwandte, bwandte},
	/*VOID   */ {_void, _void, bwandte, _void, bwandte, _absn, _absn, bwandte, bwandte, bwandte, _absn},
	/*STRING */ {bwandte, bwandte, bwandte, bwandte, bwandte, _absn, _absn, bwandte, bwandte, bwandte, bwandte},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, bwandte, _absn, bwandte, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, bwandte, _absn, bwandte, _absn},
	/*FUNC   */ {bwandte, bwandte, bwandte, bwandte, bwandte, bwandte, bwandte, bwandte, bwandte, bwandte, bwandte},
	/*ERROR  */ {bwandte, bwandte, bwandte, bwandte, bwandte, _absn, _absn, bwandte, bwandte, bwandte, bwandte},
	/*NULL   */ {bwandte, bwandte, bwandte, bwandte, bwandte, bwandte, bwandte, bwandte, bwandte, bwandte, _absn},
	/*ABSENT */ {_2___, bwandte, bwandte, _absn, bwandte, _absn, _absn, bwandte, bwandte, _absn, _absn},
}

func BIF_bitwise_and(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return bitwise_and_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ----------------------------------------------------------------
// Bitwise OR

func bitwise_or_i_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(input1.AcquireIntValue() | input2.AcquireIntValue())
}

func bworte(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromTypeErrorBinary("|", input1, input2)
}

var bitwise_or_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT              FLOAT  BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {bitwise_or_i_ii, bworte, bworte, _void, bworte, _absn, _absn, bworte, bworte, bworte, _1___},
	/*FLOAT  */ {bworte, bworte, bworte, _void, bworte, _absn, _absn, bworte, bworte, bworte, bworte},
	/*BOOL   */ {bworte, bworte, bworte, bworte, bworte, _absn, _absn, bworte, bworte, bworte, bworte},
	/*VOID   */ {_void, _void, bworte, _void, bworte, _absn, _absn, bworte, bworte, bworte, _absn},
	/*STRING */ {bworte, bworte, bworte, bworte, bworte, _absn, _absn, bworte, bworte, bworte, bworte},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, bworte, _absn, bworte, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, bworte, _absn, bworte, _absn},
	/*FUNC   */ {bworte, bworte, bworte, bworte, bworte, bworte, bworte, bworte, bworte, bworte, bworte},
	/*ERROR  */ {bworte, bworte, bworte, bworte, bworte, _absn, _absn, bworte, bworte, bworte, bworte},
	/*NULL   */ {bworte, bworte, bworte, bworte, bworte, bworte, bworte, bworte, bworte, bworte, _absn},
	/*ABSENT */ {_2___, bworte, bworte, _absn, bworte, _absn, _absn, bworte, bworte, _absn, _absn},
}

func BIF_bitwise_or(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return bitwise_or_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ----------------------------------------------------------------
// Bitwise XOR

func bitwise_xor_i_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(input1.AcquireIntValue() ^ input2.AcquireIntValue())
}

func bwxorte(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromTypeErrorBinary("^", input1, input2)
}

var bitwise_xor_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT               FLOAT  BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {bitwise_xor_i_ii, bwxorte, bwxorte, _void, bwxorte, _absn, _absn, bwxorte, bwxorte, bwxorte, _1___},
	/*FLOAT  */ {bwxorte, bwxorte, bwxorte, _void, bwxorte, _absn, _absn, bwxorte, bwxorte, bwxorte, bwxorte},
	/*BOOL   */ {bwxorte, bwxorte, bwxorte, bwxorte, bwxorte, _absn, _absn, bwxorte, bwxorte, bwxorte, bwxorte},
	/*VOID   */ {_void, _void, bwxorte, _void, bwxorte, _absn, _absn, bwxorte, bwxorte, bwxorte, _absn},
	/*STRING */ {bwxorte, bwxorte, bwxorte, bwxorte, bwxorte, _absn, _absn, bwxorte, bwxorte, bwxorte, bwxorte},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, bwxorte, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, bwxorte, _absn, _absn, _absn},
	/*FUNC   */ {bwxorte, bwxorte, bwxorte, bwxorte, bwxorte, bwxorte, bwxorte, bwxorte, bwxorte, bwxorte, bwxorte},
	/*ERROR  */ {bwxorte, bwxorte, bwxorte, bwxorte, bwxorte, _absn, _absn, bwxorte, bwxorte, bwxorte, bwxorte},
	/*NULL   */ {bwxorte, bwxorte, bwxorte, bwxorte, bwxorte, bwxorte, bwxorte, bwxorte, bwxorte, bwxorte, _absn},
	/*ABSENT */ {_2___, bwxorte, bwxorte, _absn, bwxorte, _absn, _absn, bwxorte, bwxorte, _absn, _absn},
}

func BIF_bitwise_xor(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return bitwise_xor_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ----------------------------------------------------------------
// Left shift

func lsh_i_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(input1.AcquireIntValue() << uint64(input2.AcquireIntValue()))
}

func lshfte(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromTypeErrorBinary("<<", input1, input2)
}

var left_shift_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT       FLOAT  BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {lsh_i_ii, lshfte, lshfte, _void, lshfte, _absn, _absn, lshfte, lshfte, lshfte, _1___},
	/*FLOAT  */ {lshfte, lshfte, lshfte, _void, lshfte, _absn, _absn, lshfte, lshfte, lshfte, lshfte},
	/*BOOL   */ {lshfte, lshfte, lshfte, lshfte, lshfte, _absn, _absn, lshfte, lshfte, lshfte, lshfte},
	/*VOID   */ {_void, _void, lshfte, _void, lshfte, _absn, _absn, lshfte, lshfte, lshfte, _absn},
	/*STRING */ {lshfte, lshfte, lshfte, lshfte, lshfte, _absn, _absn, lshfte, lshfte, lshfte, lshfte},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, lshfte, _absn, lshfte, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, lshfte, _absn, lshfte, _absn},
	/*FUNC   */ {lshfte, lshfte, lshfte, lshfte, lshfte, lshfte, lshfte, lshfte, lshfte, lshfte, lshfte},
	/*ERROR  */ {lshfte, lshfte, lshfte, lshfte, lshfte, _absn, _absn, lshfte, lshfte, lshfte, lshfte},
	/*NULL   */ {lshfte, lshfte, lshfte, lshfte, lshfte, lshfte, lshfte, lshfte, lshfte, lshfte, _absn},
	/*ABSENT */ {_2___, lshfte, lshfte, _absn, lshfte, _absn, _absn, lshfte, lshfte, _absn, _absn},
}

func BIF_left_shift(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return left_shift_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ----------------------------------------------------------------
// Signed right shift

func srsh_i_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(input1.AcquireIntValue() >> uint64(input2.AcquireIntValue()))
}

func srste(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromTypeErrorBinary(">>>", input1, input2)
}

var signed_right_shift_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT        FLOAT  BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {srsh_i_ii, srste, srste, _void, srste, _absn, _absn, srste, srste, srste, _1___},
	/*FLOAT  */ {srste, srste, srste, _void, srste, _absn, _absn, srste, srste, srste, srste},
	/*BOOL   */ {srste, srste, srste, srste, srste, _absn, _absn, srste, srste, srste, srste},
	/*VOID   */ {_void, _void, srste, _void, srste, _absn, _absn, srste, srste, srste, _absn},
	/*STRING */ {srste, srste, srste, srste, srste, _absn, _absn, srste, srste, srste, srste},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, srste, _absn, srste, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, srste, _absn, srste, _absn},
	/*FUNC   */ {srste, srste, srste, srste, srste, srste, srste, srste, srste, srste, srste},
	/*ERROR  */ {srste, srste, srste, srste, srste, _absn, _absn, srste, srste, srste, srste},
	/*NULL   */ {srste, srste, srste, srste, srste, srste, srste, srste, srste, srste, _absn},
	/*ABSENT */ {_2___, srste, srste, _absn, srste, _absn, _absn, srste, srste, _absn, _absn},
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

func rste(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromTypeErrorBinary(">>", input1, input2)
}

var unsigned_right_shift_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT        FLOAT  BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {ursh_i_ii, rste, rste, _void, rste, _absn, _absn, rste, rste, rste, _1___},
	/*FLOAT  */ {rste, rste, rste, _void, rste, _absn, _absn, rste, rste, rste, rste},
	/*BOOL   */ {rste, rste, rste, rste, rste, _absn, _absn, rste, rste, rste, rste},
	/*VOID   */ {_void, _void, rste, _void, rste, _absn, _absn, rste, rste, rste, _absn},
	/*STRING */ {rste, rste, rste, rste, rste, _absn, _absn, rste, rste, rste, rste},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, rste, _absn, rste, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, rste, _absn, rste, _absn},
	/*FUNC   */ {rste, rste, rste, rste, rste, rste, rste, rste, rste, rste, rste},
	/*ERROR  */ {rste, rste, rste, rste, rste, _absn, _absn, rste, rste, rste, rste},
	/*NULL   */ {rste, rste, rste, rste, rste, rste, rste, rste, rste, rste, _absn},
	/*ABSENT */ {_2___, rste, rste, _absn, rste, _absn, _absn, rste, rste, _absn, _absn},
}

func BIF_unsigned_right_shift(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return unsigned_right_shift_dispositions[input1.Type()][input2.Type()](input1, input2)
}
