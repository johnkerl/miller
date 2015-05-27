#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include "lib/mlr_globals.h"
#include "mapping/mlr_val.h"

// ----------------------------------------------------------------
mlr_val_t MV_NULL = {
	.type = MT_NULL,
	.u.int_val = 0
};
mlr_val_t MV_ERROR = {
	.type = MT_ERROR,
	.u.int_val = 0
};

// ----------------------------------------------------------------
char* mt_describe_type(int type) {
	switch (type) {
	case MT_NULL:   return "T_NULL";   break;
	case MT_ERROR:  return "T_ERROR";  break;
	case MT_BOOL:   return "T_BOOL";   break;
	case MT_DOUBLE: return "T_DOUBLE"; break;
	case MT_INT:    return "T_INT";    break;
	case MT_STRING: return "T_STRING"; break;
	default:        return "???";      break;
	}
}

// xxx cmt mem-mgt
// xxx put "alloc" in the name
char* mt_format_val(mlr_val_t* pval) {
	char* string = NULL;
	switch(pval->type) {
	case MT_NULL:
		return strdup("(null)");
		break;
	case MT_ERROR:
		return strdup("(error)");
		break;
	case MT_BOOL:
		return strdup(pval->u.bool_val ? "true" : "false");
		break;
	case MT_DOUBLE:
		// xxx what is worst-case here ...
		string = mlr_malloc_or_die(32);
		sprintf(string, MLR_GLOBALS.ofmt, pval->u.double_val);
		return string;
		break;
	case MT_INT:
		// log10(2**64) is < 20 so this is plenty.
		string = mlr_malloc_or_die(32);
		sprintf(string, "%lld", pval->u.int_val);
		return string;
		break;
	case MT_STRING:
		return strdup(pval->u.string_val);
		break;
	default:
		return strdup("???");
		break;
	}
}

char* mt_describe_val(mlr_val_t val) {
	char* stype = mt_describe_type(val.type);
	char* sval  = mt_format_val(&val);
	char* desc  = mlr_malloc_or_die(strlen(stype) + strlen(sval) + 4);
	strcpy(desc, "[");
	strcat(desc, stype);
	strcat(desc, "] ");
	strcat(desc, sval);
	return desc;
}

// ----------------------------------------------------------------
int mt_get_boolean_strict(mlr_val_t* pval) {
	if (pval->type != MT_BOOL) {
		char* desc = mt_describe_type(pval->type);
		fprintf(stderr, "Expression does not evaluate to boolean: got %s.\n", desc);
		free(desc);
		exit(1);
	}
	return pval->u.bool_val;
}

// ----------------------------------------------------------------
// xxx check for semantics comparable to mt_get_boolean_strict
void mt_get_double_strict(mlr_val_t* pval) {
	if (pval->type == MT_NULL)
		return;
	if (pval->type == MT_ERROR)
		return;
	if (pval->type == MT_DOUBLE)
		return;
	if (pval->type == MT_STRING) {
		double double_val;
		if (!mlr_try_double_from_string(pval->u.string_val, &double_val)) {
			pval->type = MT_ERROR;
			pval->u.int_val = 0;
		} else {
			pval->type = MT_DOUBLE;
			pval->u.double_val = double_val;
		}
	} else if (pval->type == MT_INT) {
		pval ->type = MT_DOUBLE;
		pval->u.double_val = (double)pval->u.int_val;
	} else if (pval->type == MT_BOOL) {
		pval->type = MT_ERROR;
		pval->u.int_val = 0;
	}
	// xxx else panic
}

// ----------------------------------------------------------------
mlr_val_t s_ss_dot_func(mlr_val_t* pval1, mlr_val_t* pval2) {
	int len1 = strlen(pval1->u.string_val);
	int len2 = strlen(pval1->u.string_val);
	int len3 = len1 + len2 + 1; // for the null-terminator byte
	char* string3 = mlr_malloc_or_die(len3);
	strcpy(&string3[0], pval1->u.string_val);
	strcpy(&string3[len1], pval2->u.string_val);

	// xxx encapsulate this:
	free(pval1->u.string_val);
	free(pval2->u.string_val);
	pval1->u.string_val = NULL;
	pval2->u.string_val = NULL;

	mlr_val_t rv = {.type = MT_STRING, .u.string_val = string3};
	return rv;
}

// ----------------------------------------------------------------
mlr_val_t s_sss_sub_func(mlr_val_t* pval1, mlr_val_t* pval2, mlr_val_t* pval3) {
	char* substr = strstr(pval1->u.string_val, pval2->u.string_val);
	if (substr == NULL) {
		return *pval1;
	} else {
		int  len1 = substr - pval1->u.string_val;
		int olen2 = strlen(pval2->u.string_val);
		int nlen2 = strlen(pval3->u.string_val);
		int  len3 = strlen(&pval1->u.string_val[len1 + olen2]);
		int  len4 = len1 + nlen2 + len3;

		char* string4 = mlr_malloc_or_die(len4);
		strncpy(&string4[0],    pval1->u.string_val, len1);
		strncpy(&string4[len1], pval3->u.string_val, nlen2);
		strncpy(&string4[len1+nlen2], &pval1->u.string_val[len1+olen2], len3);

		free(pval1->u.string_val);
		free(pval2->u.string_val);
		free(pval3->u.string_val);
		pval1->u.string_val = NULL;
		pval2->u.string_val = NULL;
		pval3->u.string_val = NULL;

		mlr_val_t rv = {.type = MT_STRING, .u.string_val = string4};
		return rv;
	}
}

// ----------------------------------------------------------------
// xxx cmt mem-mgt & contract. similar to lrec-mapper contract.
mlr_val_t s_s_tolower_func(mlr_val_t* pval1) {
	char* string = strdup(pval1->u.string_val);
	for (char* c = string; *c; c++)
		*c = tolower(*c);
	// xxx encapsulate this:
	free(pval1->u.string_val);
	pval1->u.string_val = NULL;

	mlr_val_t rv = {.type = MT_STRING, .u.string_val = string};
	return rv;
}

// xxx cmt mem-mgt & contract. similar to lrec-mapper contract.
mlr_val_t s_s_toupper_func(mlr_val_t* pval1) {
	char* string = strdup(pval1->u.string_val);
	for (char* c = string; *c; c++)
		*c = toupper(*c);
	// xxx encapsulate this:
	free(pval1->u.string_val);
	pval1->u.string_val = NULL;

	mlr_val_t rv = {.type = MT_STRING, .u.string_val = string};
	return rv;
}

// ----------------------------------------------------------------
mlr_val_t s_f_sec2gmt_func(mlr_val_t* pval1) {
	NULL_OR_ERROR_OUT(*pval1);
	mt_get_double_strict(pval1);
	if (pval1->type != MT_DOUBLE)
		return MV_ERROR;
	time_t clock = (time_t) pval1->u.double_val;
	struct tm tm;
	struct tm *ptm = gmtime_r(&clock, &tm);
	// xxx use retval which is size_t
	// xxx error-check all of this ...
	char* string = mlr_malloc_or_die(32);
	(void)strftime(string, 32, "%Y-%m-%dT%H:%M:%SZ", ptm);

	mlr_val_t rv = {.type = MT_STRING, .u.string_val = string};
	return rv;
}

mlr_val_t f_s_gmt2sec_func(mlr_val_t* pval1) {
	struct tm tm;
	strptime(pval1->u.string_val, "%Y-%m-%dT%H:%M:%SZ", &tm);
	time_t t = timegm(&tm);

	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = (double)t};
	return rv;
}

// ----------------------------------------------------------------
mlr_val_t f_s_strlen_func(mlr_val_t* pval1) {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = strlen(pval1->u.string_val)};
	return rv;
}

// ----------------------------------------------------------------
// xxx cmt us!!!!

static mlr_val_t op_n_xx(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_NULL, .u.int_val = 0}; }
static mlr_val_t op_e_xx(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_ERROR, .u.int_val = 0}; }

static  mlr_val_t eq_b_ii(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.int_val == pb->u.int_val}; }
static  mlr_val_t ne_b_ii(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.int_val != pb->u.int_val}; }
static  mlr_val_t gt_b_ii(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.int_val >  pb->u.int_val}; }
static  mlr_val_t ge_b_ii(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.int_val >= pb->u.int_val}; }
static  mlr_val_t lt_b_ii(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.int_val <  pb->u.int_val}; }
static  mlr_val_t le_b_ii(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.int_val <= pb->u.int_val}; }

static  mlr_val_t eq_b_ff(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.double_val == pb->u.double_val}; }
static  mlr_val_t ne_b_ff(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.double_val != pb->u.double_val}; }
static  mlr_val_t gt_b_ff(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.double_val >  pb->u.double_val}; }
static  mlr_val_t ge_b_ff(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.double_val >= pb->u.double_val}; }
static  mlr_val_t lt_b_ff(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.double_val <  pb->u.double_val}; }
static  mlr_val_t le_b_ff(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.double_val <= pb->u.double_val}; }

static  mlr_val_t eq_b_fi(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.double_val == pb->u.int_val}; }
static  mlr_val_t ne_b_fi(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.double_val != pb->u.int_val}; }
static  mlr_val_t gt_b_fi(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.double_val >  pb->u.int_val}; }
static  mlr_val_t ge_b_fi(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.double_val >= pb->u.int_val}; }
static  mlr_val_t lt_b_fi(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.double_val <  pb->u.int_val}; }
static  mlr_val_t le_b_fi(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.double_val <= pb->u.int_val}; }

static  mlr_val_t eq_b_if(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.int_val == pb->u.double_val}; }
static  mlr_val_t ne_b_if(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.int_val != pb->u.double_val}; }
static  mlr_val_t gt_b_if(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.int_val >  pb->u.double_val}; }
static  mlr_val_t ge_b_if(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.int_val >= pb->u.double_val}; }
static  mlr_val_t lt_b_if(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.int_val <  pb->u.double_val}; }
static  mlr_val_t le_b_if(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = pa->u.int_val <= pb->u.double_val}; }

static  mlr_val_t eq_b_xs(mlr_val_t* pa, mlr_val_t* pb) {
	char* a = mt_format_val(pa);
	mlr_val_t rv = {.type = MT_BOOL, .u.bool_val = strcmp(a, pb->u.string_val) == 0};
	free(a);
	return rv;
}
static  mlr_val_t ne_b_xs(mlr_val_t* pa, mlr_val_t* pb) {
	char* a = mt_format_val(pa);
	mlr_val_t rv = {.type = MT_BOOL, .u.bool_val = strcmp(a, pb->u.string_val) != 0};
	free(a);
	return rv;
}
static  mlr_val_t gt_b_xs(mlr_val_t* pa, mlr_val_t* pb) {
	char* a = mt_format_val(pa);
	mlr_val_t rv = {.type = MT_BOOL, .u.bool_val = strcmp(a, pb->u.string_val) >  0};
	free(a);
	return rv;
}
static  mlr_val_t ge_b_xs(mlr_val_t* pa, mlr_val_t* pb) {
	char* a = mt_format_val(pa);
	mlr_val_t rv = {.type = MT_BOOL, .u.bool_val = strcmp(a, pb->u.string_val) >= 0};
	free(a);
	return rv;
}
static  mlr_val_t lt_b_xs(mlr_val_t* pa, mlr_val_t* pb) {
	char* a = mt_format_val(pa);
	mlr_val_t rv = {.type = MT_BOOL, .u.bool_val = strcmp(a, pb->u.string_val) <  0};
	free(a);
	return rv;
}
static  mlr_val_t le_b_xs(mlr_val_t* pa, mlr_val_t* pb) {
	char* a = mt_format_val(pa);
	mlr_val_t rv = {.type = MT_BOOL, .u.bool_val = strcmp(a, pb->u.string_val) <= 0};
	free(a);
	return rv;
}

static  mlr_val_t eq_b_sx(mlr_val_t* pa, mlr_val_t* pb) {
	char* b = mt_format_val(pb);
	mlr_val_t rv = {.type = MT_BOOL, .u.bool_val = strcmp(pa->u.string_val, b) == 0};
	free(b);
	return rv;
}
static  mlr_val_t ne_b_sx(mlr_val_t* pa, mlr_val_t* pb) {
	char* b = mt_format_val(pb);
	mlr_val_t rv = {.type = MT_BOOL, .u.bool_val = strcmp(pa->u.string_val, b) != 0};
	free(b);
	return rv;
}
static  mlr_val_t gt_b_sx(mlr_val_t* pa, mlr_val_t* pb) {
	char* b = mt_format_val(pb);
	mlr_val_t rv = {.type = MT_BOOL, .u.bool_val = strcmp(pa->u.string_val, b) >  0};
	free(b);
	return rv;
}
static  mlr_val_t ge_b_sx(mlr_val_t* pa, mlr_val_t* pb) {
	char* b = mt_format_val(pb);
	mlr_val_t rv = {.type = MT_BOOL, .u.bool_val = strcmp(pa->u.string_val, b) >= 0};
	free(b);
	return rv;
}
static  mlr_val_t lt_b_sx(mlr_val_t* pa, mlr_val_t* pb) {
	char* b = mt_format_val(pb);
	mlr_val_t rv = {.type = MT_BOOL, .u.bool_val = strcmp(pa->u.string_val, b) <  0};
	free(b);
	return rv;
}
static  mlr_val_t le_b_sx(mlr_val_t* pa, mlr_val_t* pb) {
	char* b = mt_format_val(pb);
	mlr_val_t rv = {.type = MT_BOOL, .u.bool_val = strcmp(pa->u.string_val, b) <= 0};
	free(b);
	return rv;
}

static  mlr_val_t eq_b_ss(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = strcmp(pa->u.string_val, pb->u.string_val) == 0}; }
static  mlr_val_t ne_b_ss(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = strcmp(pa->u.string_val, pb->u.string_val) != 0}; }
static  mlr_val_t gt_b_ss(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = strcmp(pa->u.string_val, pb->u.string_val) >  0}; }
static  mlr_val_t ge_b_ss(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = strcmp(pa->u.string_val, pb->u.string_val) >= 0}; }
static  mlr_val_t lt_b_ss(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = strcmp(pa->u.string_val, pb->u.string_val) <  0}; }
static  mlr_val_t le_b_ss(mlr_val_t* pa, mlr_val_t* pb) { return (mlr_val_t) {.type = MT_BOOL, .u.bool_val = strcmp(pa->u.string_val, pb->u.string_val) <= 0}; }

static mv_binary_func_t* eq_dispositions[MT_MAX][MT_MAX] = {
    //         NULL      ERROR    BOOL     DOUBLE   INT      STRING
    /*NULL*/   {op_n_xx, op_e_xx, op_e_xx, op_n_xx, op_n_xx, op_n_xx},
    /*ERROR*/  {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
    /*BOOL*/   {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
    /*DOUBLE*/ {op_n_xx, op_e_xx, op_e_xx, eq_b_ff, eq_b_fi, eq_b_xs},
    /*INT*/    {op_n_xx, op_e_xx, op_e_xx, eq_b_if, eq_b_ii, eq_b_xs},
    /*STRING*/ {op_n_xx, op_e_xx, op_e_xx, eq_b_sx, eq_b_sx, eq_b_ss},
};

static mv_binary_func_t* ne_dispositions[MT_MAX][MT_MAX] = {
    //         NULL      ERROR    BOOL     DOUBLE   INT      STRING
    /*NULL*/   {op_n_xx, op_e_xx, op_e_xx, op_n_xx, op_n_xx, op_n_xx},
    /*ERROR*/  {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
    /*BOOL*/   {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
    /*DOUBLE*/ {op_n_xx, op_e_xx, op_e_xx, ne_b_ff, ne_b_fi, ne_b_xs},
    /*INT*/    {op_n_xx, op_e_xx, op_e_xx, ne_b_if, ne_b_ii, ne_b_xs},
    /*STRING*/ {op_n_xx, op_e_xx, op_e_xx, ne_b_sx, ne_b_sx, ne_b_ss},
};

static mv_binary_func_t* gt_dispositions[MT_MAX][MT_MAX] = {
    //         NULL      ERROR    BOOL     DOUBLE   INT      STRING
    /*NULL*/   {op_n_xx, op_e_xx, op_e_xx, op_n_xx, op_n_xx, op_n_xx},
    /*ERROR*/  {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
    /*BOOL*/   {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
    /*DOUBLE*/ {op_n_xx, op_e_xx, op_e_xx, gt_b_ff, gt_b_fi, gt_b_xs},
    /*INT*/    {op_n_xx, op_e_xx, op_e_xx, gt_b_if, gt_b_ii, gt_b_xs},
    /*STRING*/ {op_n_xx, op_e_xx, op_e_xx, gt_b_sx, gt_b_sx, gt_b_ss},
};

static mv_binary_func_t* ge_dispositions[MT_MAX][MT_MAX] = {
    //         NULL      ERROR    BOOL     DOUBLE   INT      STRING
    /*NULL*/   {op_n_xx, op_e_xx, op_e_xx, op_n_xx, op_n_xx, op_n_xx},
    /*ERROR*/  {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
    /*BOOL*/   {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
    /*DOUBLE*/ {op_n_xx, op_e_xx, op_e_xx, ge_b_ff, ge_b_fi, ge_b_xs},
    /*INT*/    {op_n_xx, op_e_xx, op_e_xx, ge_b_if, ge_b_ii, ge_b_xs},
    /*STRING*/ {op_n_xx, op_e_xx, op_e_xx, ge_b_sx, ge_b_sx, ge_b_ss},
};

static mv_binary_func_t* lt_dispositions[MT_MAX][MT_MAX] = {
    //         NULL      ERROR    BOOL     DOUBLE   INT      STRING
    /*NULL*/   {op_n_xx, op_e_xx, op_e_xx, op_n_xx, op_n_xx, op_n_xx},
    /*ERROR*/  {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
    /*BOOL*/   {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
    /*DOUBLE*/ {op_n_xx, op_e_xx, op_e_xx, lt_b_ff, lt_b_fi, lt_b_xs},
    /*INT*/    {op_n_xx, op_e_xx, op_e_xx, lt_b_if, lt_b_ii, lt_b_xs},
    /*STRING*/ {op_n_xx, op_e_xx, op_e_xx, lt_b_sx, lt_b_sx, lt_b_ss},
};

static mv_binary_func_t* le_dispositions[MT_MAX][MT_MAX] = {
    //         NULL      ERROR    BOOL     DOUBLE   INT      STRING
    /*NULL*/   {op_n_xx, op_e_xx, op_e_xx, op_n_xx, op_n_xx, op_n_xx},
    /*ERROR*/  {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
    /*BOOL*/   {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
    /*DOUBLE*/ {op_n_xx, op_e_xx, op_e_xx, le_b_ff, le_b_fi, le_b_xs},
    /*INT*/    {op_n_xx, op_e_xx, op_e_xx, le_b_if, le_b_ii, le_b_xs},
    /*STRING*/ {op_n_xx, op_e_xx, op_e_xx, le_b_sx, le_b_sx, le_b_ss},
};

mlr_val_t eq_op_func(mlr_val_t* pval1, mlr_val_t* pval2) { return (eq_dispositions[pval1->type][pval2->type])(pval1, pval2); }
mlr_val_t ne_op_func(mlr_val_t* pval1, mlr_val_t* pval2) { return (ne_dispositions[pval1->type][pval2->type])(pval1, pval2); }
mlr_val_t gt_op_func(mlr_val_t* pval1, mlr_val_t* pval2) { return (gt_dispositions[pval1->type][pval2->type])(pval1, pval2); }
mlr_val_t ge_op_func(mlr_val_t* pval1, mlr_val_t* pval2) { return (ge_dispositions[pval1->type][pval2->type])(pval1, pval2); }
mlr_val_t lt_op_func(mlr_val_t* pval1, mlr_val_t* pval2) { return (lt_dispositions[pval1->type][pval2->type])(pval1, pval2); }
mlr_val_t le_op_func(mlr_val_t* pval1, mlr_val_t* pval2) { return (le_dispositions[pval1->type][pval2->type])(pval1, pval2); }
