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

// ----------------------------------------------------------------
// xxx to do

//mlhmmv_xvalue_t b_x_ispresent_xfunc(mlhmmv_xvalue_t* pval1);
//mlhmmv_xvalue_t b_x_isabsent_xfunc(mlhmmv_xvalue_t* pval1);
//mlhmmv_xvalue_t b_x_ismap_xfunc(mlhmmv_xvalue_t* pval1);
//mlhmmv_xvalue_t b_x_isscalar_xfunc(mlhmmv_xvalue_t* pval1);

// isnull
// isnotnull
// isabsent
// ispresent
// isempty
// isnotempty
// isnumeric
// isint
// isfloat
// isbool
// isboolean
// isstring

//mlhmmv_xvalue_t b_m_haskey_xfunc(mlhmmv_xvalue_t* pval1);
//mlhmmv_xvalue_t s_x_typeof_xfunc(mlhmmv_xvalue_t* pval1);
//mlhmmv_xvalue_t i_m_length_xfunc(mlhmmv_xvalue_t* pval1);
//mlhmmv_xvalue_t i_m_depth_xfunc(mlhmmv_xvalue_t* pval1);
//mlhmmv_xvalue_t i_m_deepcount_xfunc(mlhmmv_xvalue_t* pval1);
//mlhmmv_xvalue_t m_mm_mapsum_xfunc(mlhmmv_xvalue_t* pval1, mlhmmv_xvalue_t* pval2);
//  - binary m,s->s JOIN
//  - binary s,s->m SPLIT
//mlhmmv_xvalue_t m_mm_mapdiff_xfunc(mlhmmv_xvalue_t* pval1, mlhmmv_xvalue_t* pval2);

//static inline mlhmmv_xvalue_t b_b_not_func(mlhmmv_xvalue_t* pval1) {
//	return mv_from_bool(!pval1->u.boolv);
//}

#endif // XVFUNCS_H
