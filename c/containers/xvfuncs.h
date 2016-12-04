#ifndef XVFUNCS_H
#define XVFUNCS_H

// ================================================================
// Functions on extended values, namely, mlrvals/hashmaps.
// ================================================================

#include "../lib/mlrutil.h"
#include "../containers/mlhmmv.h"

// ----------------------------------------------------------------
typedef mlhmmv_xvalue_t mv_variadic_func_t(
	mlhmmv_xvalue_t* pxvals,
	int              nxvals);

typedef mlhmmv_xvalue_t mv_zary_func_t();

typedef mlhmmv_xvalue_t mv_unary_func_t(
	mlhmmv_xvalue_t* pval1);

typedef mlhmmv_xvalue_t mv_binary_func_t(
	mlhmmv_xvalue_t* pval1,
	mlhmmv_xvalue_t* pval2);

typedef mlhmmv_xvalue_t mv_ternary_func_t(
	mlhmmv_xvalue_t* pval1,
	mlhmmv_xvalue_t* pval2,
	mlhmmv_xvalue_t* pval3);

//// ----------------------------------------------------------------
//static inline mlhmmv_xvalue_t b_b_not_func(mlhmmv_xvalue_t* pval1) {
//	return mv_from_bool(!pval1->u.boolv);
//}
//
//static inline mlhmmv_xvalue_t b_bb_or_func(mlhmmv_xvalue_t* pval1, mlhmmv_xvalue_t* pval2) {
//	return mv_from_bool(pval1->u.boolv || pval2->u.boolv);
//}
//static inline mlhmmv_xvalue_t b_bb_and_func(mlhmmv_xvalue_t* pval1, mlhmmv_xvalue_t* pval2) {
//	return mv_from_bool(pval1->u.boolv && pval2->u.boolv);
//}
//static inline mlhmmv_xvalue_t b_bb_xor_func(mlhmmv_xvalue_t* pval1, mlhmmv_xvalue_t* pval2) {
//	return mv_from_bool(pval1->u.boolv ^ pval2->u.boolv);
//}

//  - unary m->b ISPRESENT ISABSENT ISMAP ISSCALAR
//  - unary m->b HASKEY
//  - unary m->i LENGTH DEPTH DEEPCOUNT
//  - unary x->s TYPEOF
//  - binary m,s->s JOIN
//
//  - binary s,s->m SPLIT
//  - binary/variadic m,m->m MAPSUM MAPDIFF

#endif // XVFUNCS_H
