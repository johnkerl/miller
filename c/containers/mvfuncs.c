#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/mlrdatetime.h"
#include "lib/mlrregex.h"
#include "containers/mvfuncs.h"

// ================================================================
// See important notes at the top of mlrval.h.
// ================================================================

typedef int mv_i_nn_comparator_func_t(mv_t* pa, mv_t* pb);
typedef int mv_i_cncn_comparator_func_t(const mv_t* pa, const mv_t* pb);

// ----------------------------------------------------------------
// Keystroke-savers for disposition matrices:

static mv_t _a(mv_t* pa, mv_t* pb) {
	return mv_absent();
}
static mv_t _emt(mv_t* pa, mv_t* pb) {
	return mv_empty();
}
static mv_t _err(mv_t* pa, mv_t* pb) {
	return mv_error();
}
static mv_t _0(mv_t* pa) {
	return *pa;
}
static mv_t _1(mv_t* pa, mv_t* pb) {
	return *pa;
}
static mv_t _2(mv_t* pa, mv_t* pb) {
	return *pb;
}
static mv_t _s1(mv_t* pa, mv_t* pb) {
	return s_x_string_func(pa);
}
static mv_t _s2(mv_t* pa, mv_t* pb) {
	return s_x_string_func(pb);
}
static mv_t _f0(mv_t* pa, mv_t* pb) {
	return mv_from_float(0.0);
}
static mv_t _i0(mv_t* pa, mv_t* pb) {
	return mv_from_int(0LL);
}
static mv_t _a1(mv_t* pa) {
	return mv_absent();
}
static mv_t _emt1(mv_t* pa) {
	return mv_empty();
}
static mv_t _err1(mv_t* pa) {
	return mv_error();
}

// ----------------------------------------------------------------
static mv_t dot_strings(char* string1, char* string2) {
	int len1 = strlen(string1);
	int len2 = strlen(string2);
	int len3 = len1 + len2 + 1; // for the null-terminator byte
	char* string3 = mlr_malloc_or_die(len3);
	strcpy(&string3[0], string1);
	strcpy(&string3[len1], string2);
	return mv_from_string_with_free(string3);
}

mv_t dot_s_ss(mv_t* pval1, mv_t* pval2) {
	mv_t rv = dot_strings(pval1->u.strv, pval2->u.strv);
	mv_free(pval1);
	mv_free(pval2);
	return rv;
}

mv_t dot_s_xs(mv_t* pval1, mv_t* pval2) {
	mv_t sval1 = s_x_string_func(pval1);
	mv_free(pval1);
	mv_t rv = dot_strings(sval1.u.strv, pval2->u.strv);
	mv_free(&sval1);
	mv_free(pval2);
	return rv;
}

mv_t dot_s_sx(mv_t* pval1, mv_t* pval2) {
	mv_t sval2 = s_x_string_func(pval2);
	mv_free(pval2);
	mv_t rv = dot_strings(pval1->u.strv, sval2.u.strv);
	mv_free(pval1);
	mv_free(&sval2);
	return rv;
}

mv_t dot_s_xx(mv_t* pval1, mv_t* pval2) {
	mv_t sval1 = s_x_string_func(pval1);
	mv_t sval2 = s_x_string_func(pval2);
	mv_t rv = dot_strings(sval1.u.strv, sval2.u.strv);
	mv_free(&sval1);
	mv_free(&sval2);
	return rv;
}

static mv_binary_func_t* dot_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY STRING     INT        FLOAT      BOOL
	/*ERROR*/  {_err, _err,  _err, _err,      _err,      _err,      _err},
	/*ABSENT*/ {_err, _a,    _emt, _2,        _s2,       _s2,       _s2},
	/*EMPTY*/  {_err, _emt,  _emt, _2,        _s2,       _s2,       _s2},
	/*STRING*/ {_err, _1,    _1,   dot_s_ss,  dot_s_sx,  dot_s_sx,  dot_s_sx},
	/*INT*/    {_err, _s1,   _s1,  dot_s_xs,  dot_s_xx,  dot_s_xx,  dot_s_xx},
	/*FLOAT*/  {_err, _s1,   _s1,  dot_s_xs,  dot_s_xx,  dot_s_xx,  dot_s_xx},
	/*BOOL*/   {_err, _s1,   _s1,  dot_s_xs,  dot_s_xx,  dot_s_xx,  dot_s_xx},
};

mv_t s_xx_dot_func(mv_t* pval1, mv_t* pval2) { return (dot_dispositions[pval1->type][pval2->type])(pval1,pval2); }

// ----------------------------------------------------------------
mv_t sub_no_precomp_func(mv_t* pval1, mv_t* pval2, mv_t* pval3) {
	regex_t regex;
	string_builder_t *psb = sb_alloc(MV_SB_ALLOC_LENGTH);
	mv_t rv = sub_precomp_func(pval1, regcomp_or_die(&regex, pval2->u.strv, 0), psb, pval3);
	sb_free(psb);
	regfree(&regex);
	mv_free(pval2);
	return rv;
}

// ----------------------------------------------------------------
// Example:
// * pval1->u.strv = "hello"
// * regex = "l+"
// * pval3->u.strv = "yyy"
//
// *  len1 = 2 = length of "he"
// * olen2 = 2 = length of "ll"
// * nlen2 = 3 = length of "yyy"
// *  len3 = 1 = length of "o"
// *  len4 = 6 = 2+3+1

mv_t sub_precomp_func(mv_t* pval1, regex_t* pregex, string_builder_t* psb, mv_t* pval3) {
	int matched      = FALSE;
	int all_captured = FALSE;
	char* input      = pval1->u.strv;
	char* output     = regex_sub(input, pregex, psb, pval3->u.strv, &matched, &all_captured);

	mv_free(pval1);
	mv_free(pval3);
	return mv_from_string_with_free(output);
}

// ----------------------------------------------------------------
// Example:
// * pval1->u.strv = "hello"
// * regex = "l+"
// * pval3->u.strv = "yyy"
//
// *  len1 = 2 = length of "he"
// * olen2 = 2 = length of "ll"
// * nlen2 = 3 = length of "yyy"
// *  len3 = 1 = length of "o"
// *  len4 = 6 = 2+3+1

mv_t gsub_no_precomp_func(mv_t* pval1, mv_t* pval2, mv_t* pval3) {
	regex_t regex;
	string_builder_t *psb = sb_alloc(MV_SB_ALLOC_LENGTH);
	mv_t rv = gsub_precomp_func(pval1, regcomp_or_die(&regex, pval2->u.strv, 0), psb, pval3);
	sb_free(psb);
	regfree(&regex);
	mv_free(pval2);
	return rv;
}

mv_t gsub_precomp_func(mv_t* pval1, regex_t* pregex, string_builder_t* psb, mv_t* pval3) {
	int matched      = FALSE;
	int all_captured = FALSE;
	char* input      = pval1->u.strv;
	char free_flags  = NO_FREE;
	char* output     = regex_gsub(input, pregex, psb, pval3->u.strv, &matched, &all_captured, &free_flags);

	mv_free(pval1);
	mv_free(pval3);
	return mv_from_string(output, free_flags);
}

// ----------------------------------------------------------------
mv_t i_iii_modadd_func(mv_t* pval1, mv_t* pval2, mv_t* pval3) {
	long long m = pval3->u.intv;
	if (m <= 0LL)
		return mv_error();
	long long a = pval1->u.intv % m;
	if (a < 0LL)
		a += m; // crazy C-language mod operator
	long long b = pval2->u.intv % m;
	if (b < 0LL)
		b += m;
	long long c = (a + b) % m;
	return mv_from_int(c);
}

mv_t i_iii_modsub_func(mv_t* pval1, mv_t* pval2, mv_t* pval3) {
	long long m = pval3->u.intv;
	if (m <= 0LL)
		return mv_error();
	long long a = pval1->u.intv % m;
	if (a < 0LL)
		a += m; // crazy C-language mod operator
	long long b = pval2->u.intv % m;
	if (b < 0LL)
		b += m;
	long long c = (a - b) % m;
	if (c < 0LL)
		c += m;
	return mv_from_int(c);
}

mv_t i_iii_modmul_func(mv_t* pval1, mv_t* pval2, mv_t* pval3) {
	long long m = pval3->u.intv;
	if (m <= 0LL)
		return mv_error();
	long long a = pval1->u.intv % m;
	if (a < 0LL)
		a += m; // crazy C-language mod operator
	long long b = pval2->u.intv % m;
	if (b < 0LL)
		b += m;
	long long c = (a * b) % m;
	return mv_from_int(c);
}

mv_t i_iii_modexp_func(mv_t* pval1, mv_t* pval2, mv_t* pval3) {
	long long m = pval3->u.intv;
	if (m <= 0LL)
		return mv_error();
	long long a = pval1->u.intv % m;
	if (a < 0LL)
		a += m; // crazy C-language mod operator

	long long e = pval2->u.intv;

	long long c = 1LL;
	if (e == 1LL) {
		c = a;
	} else if (e == 0LL) {
		c = 1LL;
	} else if (e > 0) {
		long long ap = a;
		c = 1LL;
		unsigned long long u = (unsigned long long)e;

		// repeated-squaring algorithm
		while (u != 0) {
			if ((u & 1LL) == 1LL) {
				c = (c * ap) % m;
			}
			u >>= 1;
			ap = (ap * ap) % m;
		}
	} else {
		return mv_error();
	}

	return mv_from_int(c);
}

// ----------------------------------------------------------------
mv_t s_s_tolower_func(mv_t* pval1) {
	char* string = mlr_strdup_or_die(pval1->u.strv);
	for (char* c = string; *c; c++)
		*c = tolower((unsigned char)*c);
	mv_free(pval1);
	pval1->u.strv = NULL;

	return mv_from_string_with_free(string);
}

mv_t s_s_toupper_func(mv_t* pval1) {
	char* string = mlr_strdup_or_die(pval1->u.strv);
	for (char* c = string; *c; c++)
		*c = toupper((unsigned char)*c);
	mv_free(pval1);
	pval1->u.strv = NULL;

	return mv_from_string_with_free(string);
}

mv_t i_s_strlen_func(mv_t* pval1) {
	mv_t rv = mv_from_int(strlen_for_utf8_display(pval1->u.strv));
	mv_free(pval1);
	return rv;
}

mv_t s_x_typeof_func(mv_t* pval1) {
	mv_t rv = mv_from_string(mt_describe_type(pval1->type), NO_FREE);
	mv_free(pval1);
	return rv;
}

// ----------------------------------------------------------------
// Precondition: psec is either int or float.
mv_t time_string_from_seconds(mv_t* psec, char* format) {
	double seconds_since_the_epoch = 0.0;
	if (psec->type == MT_FLOAT) {
		if (isinf(psec->u.fltv) || isnan(psec->u.fltv)) {
			return mv_error();
		}
		seconds_since_the_epoch = psec->u.fltv;
	} else {
		seconds_since_the_epoch = psec->u.intv;
	}

	char* string = mlr_alloc_time_string_from_seconds(seconds_since_the_epoch, format);

	return mv_from_string_with_free(string);
}

// ----------------------------------------------------------------
static mv_t sec2gmt_s_n(mv_t* pa) { return time_string_from_seconds(pa, ISO8601_TIME_FORMAT); }

static mv_unary_func_t* sec2gmt_dispositions[MT_DIM] = {
	/*ERROR*/  _err1,
	/*ABSENT*/ _a1,
	/*EMPTY*/  _emt1,
	/*STRING*/ _0,
	/*INT*/    sec2gmt_s_n,
	/*FLOAT*/  sec2gmt_s_n,
	/*BOOL*/   _0,
};

mv_t s_x_sec2gmt_func(mv_t* pval1) { return (sec2gmt_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t sec2gmtdate_s_n(mv_t* pa) { return time_string_from_seconds(pa, ISO8601_DATE_FORMAT); }

static mv_unary_func_t* sec2gmtdate_dispositions[MT_DIM] = {
	/*ERROR*/  _err1,
	/*ABSENT*/ _a1,
	/*EMPTY*/  _emt1,
	/*STRING*/ _0,
	/*INT*/    sec2gmtdate_s_n,
	/*FLOAT*/  sec2gmtdate_s_n,
	/*BOOL*/   _0,
};

mv_t s_x_sec2gmtdate_func(mv_t* pval1) { return (sec2gmtdate_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
mv_t s_ns_strftime_func(mv_t* pval1, mv_t* pval2) {
	mv_t rv = time_string_from_seconds(pval1, pval2->u.strv);
	mv_free(pval2);
	return rv;
}

// ----------------------------------------------------------------
static mv_t seconds_from_time_string(char* string, char* format) {
	if (*string == '\0') {
		return mv_empty();
	} else {
		time_t seconds = mlr_seconds_from_time_string(string, format);
		return mv_from_int((long long)seconds);
	}
}

mv_t i_s_gmt2sec_func(mv_t* pval1) {
	mv_t rv = seconds_from_time_string(pval1->u.strv, ISO8601_TIME_FORMAT);
	mv_free(pval1);
	return rv;
}

mv_t i_ss_strptime_func(mv_t* pval1, mv_t* pval2) {
	mv_t rv = seconds_from_time_string(pval1->u.strv, pval2->u.strv);
	mv_free(pval1);
	mv_free(pval2);
	return rv;
}

// ----------------------------------------------------------------
static void split_ull_to_hms(long long u, long long* ph, long long* pm, long long* ps) {
	long long h = 0LL, m = 0LL, s = 0LL;
	long long sign = 1LL;
	if (u < 0LL) {
		u = -u;
		sign = -1LL;
	}
	s = u % 60LL;
	u = u / 60LL;
	if (u == 0LL) {
		s = s * sign;
	} else {
		m = u % 60LL;
		u = u / 60LL;
		if (u == 0LL) {
			m = m * sign;
		} else {
			h = u * sign;
		}
	}
	*ph = h;
	*pm = m;
	*ps = s;
}

static void split_ull_to_dhms(long long u, long long* pd, long long* ph, long long* pm, long long* ps) {
	long long d = 0LL, h = 0LL, m = 0LL, s = 0LL;
	long long sign = 1LL;
	if (u < 0LL) {
		u = -u;
		sign = -1LL;
	}
	s = u % 60LL;
	u = u / 60LL;
	if (u == 0LL) {
		s = s * sign;
	} else {
		m = u % 60LL;
		u = u / 60LL;
		if (u == 0LL) {
			m = m * sign;
		} else {
			h = u % 24LL;
			u = u / 24LL;
			if (u == 0LL) {
				h = h * sign;
			} else {
				d = u * sign;
			}
		}
	}
	*pd = d;
	*ph = h;
	*pm = m;
	*ps = s;
}

mv_t s_i_sec2hms_func(mv_t* pval1) {
	long long u = pval1->u.intv;
	long long h, m, s;
	char* fmt = "%02lld:%02lld:%02lld";
	if (u < 0) {
		u = -u;
		fmt = "-%02lld:%02lld:%02lld";
	}
	split_ull_to_hms(u, &h, &m, &s);
	int n = snprintf(NULL, 0, fmt, h, m, s);
	char* string = mlr_malloc_or_die(n+1);
	sprintf(string, fmt, h, m, s);
	return mv_from_string_with_free(string);
}

mv_t s_f_fsec2hms_func(mv_t* pval1) {
	double v = fabs(pval1->u.fltv);
	long long h, m, s;
	char* fmt = "%lld:%02lld:%09.6lf";
	long long u = (long long)trunc(v);
	double f = v - u;
	if (pval1->u.fltv < 0.0) {
		fmt = "-%02lld:%02lld:%09.6lf";
	}
	split_ull_to_hms(u, &h, &m, &s);
	int n = snprintf(NULL, 0, fmt, h, m, s+f);
	char* string = mlr_malloc_or_die(n+1);
	sprintf(string, fmt, h, m, s+f);
	return mv_from_string_with_free(string);
}

mv_t s_i_sec2dhms_func(mv_t* pval1) {
	long long u = pval1->u.intv;
	long long d, h, m, s;
	split_ull_to_dhms(u, &d, &h, &m, &s);
	if (d != 0.0) {
		char* fmt = "%lldd%02lldh%02lldm%02llds";
		int n = snprintf(NULL, 0, fmt, d, h, m, s);
		char* string = mlr_malloc_or_die(n+1);
		sprintf(string, fmt, d, h, m, s);
		return mv_from_string_with_free(string);
	} else if (h != 0.0) {
		char* fmt = "%lldh%02lldm%02llds";
		int n = snprintf(NULL, 0, fmt, h, m, s);
		char* string = mlr_malloc_or_die(n+1);
		sprintf(string, fmt, h, m, s);
		return mv_from_string_with_free(string);
	} else if (m != 0.0) {
		char* fmt = "%lldm%02llds";
		int n = snprintf(NULL, 0, fmt, m, s);
		char* string = mlr_malloc_or_die(n+1);
		sprintf(string, fmt, m, s);
		return mv_from_string_with_free(string);
	} else {
		char* fmt = "%llds";
		int n = snprintf(NULL, 0, fmt, s);
		char* string = mlr_malloc_or_die(n+1);
		sprintf(string, fmt, s);
		return mv_from_string_with_free(string);
	}
}

mv_t s_f_fsec2dhms_func(mv_t* pval1) {
	double v = fabs(pval1->u.fltv);
	long long sign = pval1->u.fltv < 0.0 ? -1LL : 1LL;
	long long d, h, m, s;
	long long u = (long long)trunc(v);
	double f = v - u;
	split_ull_to_dhms(u, &d, &h, &m, &s);
	if (d != 0.0) {
		d = sign * d;
		char* fmt = "%lldd%02lldh%02lldm%09.6lfs";
		int n = snprintf(NULL, 0, fmt, d, h, m, s+f);
		char* string = mlr_malloc_or_die(n+1);
		sprintf(string, fmt, d, h, m, s+f);
		return mv_from_string_with_free(string);
	} else if (h != 0.0) {
		h = sign * h;
		char* fmt = "%lldh%02lldm%09.6lfs";
		int n = snprintf(NULL, 0, fmt, h, m, s+f);
		char* string = mlr_malloc_or_die(n+1);
		sprintf(string, fmt, h, m, s+f);
		return mv_from_string_with_free(string);
	} else if (m != 0.0) {
		m = sign * m;
		char* fmt = "%lldm%09.6lfs";
		int n = snprintf(NULL, 0, fmt, m, s+f);
		char* string = mlr_malloc_or_die(n+1);
		sprintf(string, fmt, m, s+f);
		return mv_from_string_with_free(string);
	} else {
		s = sign * s;
		f = sign * f;
		char* fmt = "%.6lfs";
		int n = snprintf(NULL, 0, fmt, s+f);
		char* string = mlr_malloc_or_die(n+1);
		sprintf(string, fmt, s+f);
		return mv_from_string_with_free(string);
	}
}

// ----------------------------------------------------------------
mv_t i_s_hms2sec_func(mv_t* pval1) {
	long long h = 0LL, m = 0LL, s = 0LL;
	long long sec = 0LL;
	char* p = pval1->u.strv;
	long long sign = 1LL;
	if (*p == '-') {
		p++;
		sign = -1LL;
	}
	if (sscanf(p, "%lld:%lld:%lld", &h, &m, &s) == 3) {
		if (h >= 0LL)
			sec = 3600LL*h + 60LL*m + s;
		else
			sec = -(-3600LL*h + 60LL*m + s);
	} else if (sscanf(p, "%lld:%lld", &m, &s) == 2) {
		if (m >= 0LL)
			sec = 60LL*m + s;
		else
			sec = -(-60LL*m + s);
	} else if (sscanf(p, "%lld", &s) == 1) {
		sec = s;
	} else {
		mv_free(pval1);
		return mv_error();
	}
	mv_free(pval1);
	return mv_from_int(sec * sign);
}

mv_t f_s_hms2fsec_func(mv_t* pval1) {
	long long h = 0LL, m = 0LL;
	double s = 0.0;
	double sec = 0.0;
	char* p = pval1->u.strv;
	double sign = 1.0;
	if (*p == '-') {
		p++;
		sign = -1.0;
	}
	if (sscanf(p, "%lld:%lld:%lf", &h, &m, &s) == 3) {
		sec = 3600*h + 60*m + s;
	} else if (sscanf(p, "%lld:%lf", &m, &s) == 2) {
		sec = 60*m + s;
	} else if (sscanf(p, "%lf", &s) == 2) {
		sec = s;
	} else {
		mv_free(pval1);
		return mv_error();
	}
	mv_free(pval1);
	return mv_from_float(sec * sign);
}

mv_t i_s_dhms2sec_func(mv_t* pval1) {
	long long d = 0LL, h = 0LL, m = 0LL, s = 0LL;
	long long sec = 0LL;
	char* p = pval1->u.strv;
	long long sign = 1LL;
	if (*p == '-') {
		p++;
		sign = -1LL;
	}
	if (sscanf(p, "%lldd%lldh%lldm%llds", &d, &h, &m, &s) == 4) {
		sec = 86400*d + 3600*h + 60*m + s;
	} else if (sscanf(p, "%lldh%lldm%llds", &h, &m, &s) == 3) {
		sec = 3600*h + 60*m + s;
	} else if (sscanf(p, "%lldm%llds", &m, &s) == 2) {
		sec = 60*m + s;
	} else if (sscanf(p, "%llds", &s) == 1) {
		sec = s;
	} else {
		mv_free(pval1);
		return mv_error();
	}
	mv_free(pval1);
	return mv_from_int(sec * sign);
}

mv_t f_s_dhms2fsec_func(mv_t* pval1) {
	long long d = 0LL, h = 0LL, m = 0LL;
	double s = 0.0;
	double sec = 0.0;
	char* p = pval1->u.strv;
	long long sign = 1.0;
	if (*p == '-') {
		p++;
		sign = -1.0;
	}
	if (sscanf(p, "%lldd%lldh%lldm%lfs", &d, &h, &m, &s) == 4) {
		sec = 86400*d + 3600*h + 60*m + s;
	} else if (sscanf(p, "%lldh%lldm%lfs", &h, &m, &s) == 3) {
		sec = 3600*h + 60*m + s;
	} else if (sscanf(p, "%lldm%lfs", &m, &s) == 2) {
		sec = 60*m + s;
	} else if (sscanf(p, "%lfs", &s) == 1) {
		sec = s;
	} else {
		mv_free(pval1);
		return mv_error();
	}
	mv_free(pval1);
	return mv_from_float(sec * sign);
}

// ----------------------------------------------------------------
static mv_t plus_f_ff(mv_t* pa, mv_t* pb) {
	double a = pa->u.fltv;
	double b = pb->u.fltv;
	return mv_from_float(a + b);
}
static mv_t plus_f_fi(mv_t* pa, mv_t* pb) {
	double a = pa->u.fltv;
	double b = (double)pb->u.intv;
	return mv_from_float(a + b);
}
static mv_t plus_f_if(mv_t* pa, mv_t* pb) {
	double a = (double)pa->u.intv;
	double b = pb->u.fltv;
	return mv_from_float(a + b);
}
// Adds & subtracts overflow by at most one bit so it suffices to check
// sign-changes.
static mv_t plus_n_ii(mv_t* pa, mv_t* pb) {
	long long a = pa->u.intv;
	long long b = pb->u.intv;
	long long c = a + b;

	int overflowed = FALSE;
	if (a > 0LL) {
		if (b > 0LL && c < 0LL)
			overflowed = TRUE;
	} else if (a < 0LL) {
		if (b < 0LL && c > 0LL)
			overflowed = TRUE;
	}

	if (overflowed) {
		return mv_from_float((double)a + (double)b);
	} else {
		return mv_from_int(c);
	}
}

static mv_binary_func_t* plus_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY STRING INT        FLOAT      BOOL
	/*ERROR*/  {_err, _err,  _err, _err,  _err,      _err,      _err},
	/*ABSENT*/ {_err, _a,    _a,   _err,  _2,        _2,        _err},
	/*EMPTY*/  {_err, _a,    _emt, _err,  _emt,      _emt,      _err},
	/*STRING*/ {_err, _err,  _err, _err,  _err,      _err,      _err},
	/*INT*/    {_err, _1,    _emt, _err,  plus_n_ii, plus_f_if, _err},
	/*FLOAT*/  {_err, _1,    _emt, _err,  plus_f_fi, plus_f_ff, _err},
	/*BOOL*/   {_err, _err,  _err, _err,  _err,      _err,      _err},
};

mv_t x_xx_plus_func(mv_t* pval1, mv_t* pval2) { return (plus_dispositions[pval1->type][pval2->type])(pval1,pval2); }

// ----------------------------------------------------------------
static mv_t minus_f_ff(mv_t* pa, mv_t* pb) {
	double a = pa->u.fltv;
	double b = pb->u.fltv;
	return mv_from_float(a - b);
}
static mv_t minus_f_fi(mv_t* pa, mv_t* pb) {
	double a = pa->u.fltv;
	double b = (double)pb->u.intv;
	return mv_from_float(a - b);
}
static mv_t minus_f_if(mv_t* pa, mv_t* pb) {
	double a = (double)pa->u.intv;
	double b = pb->u.fltv;
	return mv_from_float(a - b);
}
// Adds & subtracts overflow by at most one bit so it suffices to check
// sign-changes.
static mv_t minus_n_ii(mv_t* pa, mv_t* pb) {
	long long a = pa->u.intv;
	long long b = pb->u.intv;
	long long c = a - b;

	int overflowed = FALSE;
	if (a > 0LL) {
		if (b < 0LL && c < 0LL)
			overflowed = TRUE;
	} else if (a < 0LL) {
		if (b > 0LL && c > 0LL)
			overflowed = TRUE;
	}

	if (overflowed) {
		return mv_from_float((double)a + (double)b);
	} else {
		return mv_from_int(c);
	}
}

static mv_binary_func_t* minus_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY STRING INT         FLOAT       BOOL
	/*ERROR*/  {_err, _err,  _err, _err,  _err,       _err,       _err},
	/*ABSENT*/ {_err, _a,    _a,   _err,  _2,         _2,         _err},
	/*EMPTY*/  {_err, _a,    _emt, _err,  _emt,       _emt,       _err},
	/*STRING*/ {_err, _err,  _err, _err,  _err,       _err,       _err},
	/*INT*/    {_err, _1,    _emt, _err,  minus_n_ii, minus_f_if, _err},
	/*FLOAT*/  {_err, _1,    _emt, _err,  minus_f_fi, minus_f_ff, _err},
	/*BOOL*/   {_err, _err,  _err, _err,  _err,       _err,       _err},
};

mv_t x_xx_minus_func(mv_t* pval1, mv_t* pval2) { return (minus_dispositions[pval1->type][pval2->type])(pval1,pval2); }

// ----------------------------------------------------------------
static mv_t times_f_ff(mv_t* pa, mv_t* pb) {
	double a = pa->u.fltv;
	double b = pb->u.fltv;
	return mv_from_float(a * b);
}
static mv_t times_f_fi(mv_t* pa, mv_t* pb) {
	double a = pa->u.fltv;
	double b = (double)pb->u.intv;
	return mv_from_float(a * b);
}
static mv_t times_f_if(mv_t* pa, mv_t* pb) {
	double a = (double)pa->u.intv;
	double b = pb->u.fltv;
	return mv_from_float(a * b);
}
// Unlike adds & subtracts which overflow by at most one bit, multiplies can
// overflow by a word size. Thus detecting sign-changes does not suffice to
// detect overflow. Instead we test whether the floating-point product exceeds
// the representable integer range. Now 64-bit integers have 64-bit precision
// while IEEE-doubles have only 52-bit mantissas -- so, 53 bits including
// implicit leading one.
//
// The following experiment explicitly demonstrates the resolution at this range:
//
//    64-bit integer     64-bit integer     Casted to double           Back to 64-bit
//        in hex           in decimal                                    integer
// 0x7ffffffffffff9ff 9223372036854774271 9223372036854773760.000000 0x7ffffffffffff800
// 0x7ffffffffffffa00 9223372036854774272 9223372036854773760.000000 0x7ffffffffffff800
// 0x7ffffffffffffbff 9223372036854774783 9223372036854774784.000000 0x7ffffffffffffc00
// 0x7ffffffffffffc00 9223372036854774784 9223372036854774784.000000 0x7ffffffffffffc00
// 0x7ffffffffffffdff 9223372036854775295 9223372036854774784.000000 0x7ffffffffffffc00
// 0x7ffffffffffffe00 9223372036854775296 9223372036854775808.000000 0x8000000000000000
// 0x7ffffffffffffffe 9223372036854775806 9223372036854775808.000000 0x8000000000000000
// 0x7fffffffffffffff 9223372036854775807 9223372036854775808.000000 0x8000000000000000
//
// That is, we cannot check an integer product to see if it is greater than
// 2**63-1 (or is less than -2**63) using integer arithmetic (it may have
// already overflowed) *or* using double-precision (granularity). Instead we
// check if the absolute value of the product exceeds the largest representable
// double less than 2**63. (An alterative would be to do all integer multipies
// using handcrafted multi-word 128-bit arithmetic).
static mv_t times_n_ii(mv_t* pa, mv_t* pb) {
	long long a = pa->u.intv;
	long long b = pb->u.intv;

	double d = (double)a * (double)b;
	if (fabs(d) > 9223372036854774784.0) {
		return mv_from_float(d);
	} else {
		return mv_from_int(a * b);
	}
}

static mv_binary_func_t* times_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY STRING INT         FLOAT       BOOL
	/*ERROR*/  {_err, _err,  _err, _err,  _err,       _err,       _err},
	/*ABSENT*/ {_err, _a,    _a,   _err,  _2,         _2,         _err},
	/*EMPTY*/  {_err, _a,    _emt, _err,  _emt,       _emt,       _err},
	/*STRING*/ {_err, _err,  _err, _err,  _err,       _err,       _err},
	/*INT*/    {_err, _1,    _emt, _err,  times_n_ii, times_f_if, _err},
	/*FLOAT*/  {_err, _1,    _emt, _err,  times_f_fi, times_f_ff, _err},
	/*BOOL*/   {_err, _err,  _err, _err,  _err,       _err,       _err},
};

mv_t x_xx_times_func(mv_t* pval1, mv_t* pval2) { return (times_dispositions[pval1->type][pval2->type])(pval1,pval2); }

// ----------------------------------------------------------------
static mv_t divide_f_ff(mv_t* pa, mv_t* pb) {
	double a = pa->u.fltv;
	double b = pb->u.fltv;
	return mv_from_float(a / b);
}
static mv_t divide_f_fi(mv_t* pa, mv_t* pb) {
	double a = pa->u.fltv;
	double b = (double)pb->u.intv;
	return mv_from_float(a / b);
}
static mv_t divide_f_if(mv_t* pa, mv_t* pb) {
	double a = (double)pa->u.intv;
	double b = pb->u.fltv;
	return mv_from_float(a / b);
}
static mv_t divide_i_ii(mv_t* pa, mv_t* pb) {
	long long a = pa->u.intv;
	long long b = pb->u.intv;
	long long r = a % b;
	// Pythonic division, not C division.
	if (r == 0LL) {
		return mv_from_int(a / b);
	} else {
		return mv_from_float((double)a / (double)b);
	}
}

static mv_binary_func_t* divide_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY STRING INT          FLOAT        BOOL
	/*ERROR*/  {_err, _err,  _err, _err,  _err,        _err,        _err},
	/*ABSENT*/ {_err, _a,    _a,   _err,  _i0,         _f0,         _err},
	/*EMPTY*/  {_err, _a,    _emt, _err,  _emt,        _emt,        _err},
	/*STRING*/ {_err, _err,  _err, _err,  _err,        _err,        _err},
	/*INT*/    {_err, _1,    _emt, _err,  divide_i_ii, divide_f_if, _err},
	/*FLOAT*/  {_err, _1,    _emt, _err,  divide_f_fi, divide_f_ff, _err},
	/*BOOL*/   {_err, _err,  _err, _err,  _err,        _err,        _err},
};

mv_t x_xx_divide_func(mv_t* pval1, mv_t* pval2) { return (divide_dispositions[pval1->type][pval2->type])(pval1,pval2); }

// ----------------------------------------------------------------
static mv_t idiv_f_ff(mv_t* pa, mv_t* pb) {
	double a = pa->u.fltv;
	double b = pb->u.fltv;
	return mv_from_float(floor(a / b));
}
static mv_t idiv_f_fi(mv_t* pa, mv_t* pb) {
	double a = pa->u.fltv;
	double b = (double)pb->u.intv;
	return mv_from_float(floor(a / b));
}
static mv_t idiv_f_if(mv_t* pa, mv_t* pb) {
	double a = (double)pa->u.intv;
	double b = pb->u.fltv;
	return mv_from_float(floor(a / b));
}
static mv_t idiv_i_ii(mv_t* pa, mv_t* pb) {
	long long a = pa->u.intv;
	long long b = pb->u.intv;
	// Pythonic division, not C division.
	long long q = a / b;
	long long r = a % b;
	if (a < 0) {
		if (b > 0) {
			if (r != 0)
				q--;
		}
	} else {
		if (b < 0) {
			if (r != 0)
				q--;
		}
	}
	return mv_from_int(q);
}

static mv_binary_func_t* idiv_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY STRING INT        FLOAT      BOOL
	/*ERROR*/  {_err, _err,  _err, _err,  _err,      _err,      _err},
	/*ABSENT*/ {_err, _a,    _a,   _err,  _i0,       _f0,       _err},
	/*EMPTY*/  {_err, _a,    _emt, _err,  _emt,      _emt,      _err},
	/*STRING*/ {_err, _err,  _err, _err,  _err,      _err,      _err},
	/*INT*/    {_err, _1,    _emt, _err,  idiv_i_ii, idiv_f_if, _err},
	/*FLOAT*/  {_err, _1,    _emt, _err,  idiv_f_fi, idiv_f_ff, _err},
	/*BOOL*/   {_err, _err,  _err, _err,  _err,      _err,      _err},
};

mv_t x_xx_int_divide_func(mv_t* pval1, mv_t* pval2) {
	return (idiv_dispositions[pval1->type][pval2->type])(pval1,pval2);
}

// ----------------------------------------------------------------
static mv_t mod_f_ff(mv_t* pa, mv_t* pb) {
	double a = pa->u.fltv;
	double b = pb->u.fltv;
	return mv_from_float(a - b * floor(a / b));
}
static mv_t mod_f_fi(mv_t* pa, mv_t* pb) {
	double a = pa->u.fltv;
	double b = (double)pb->u.intv;
	return mv_from_float(a - b * floor(a / b));
}
static mv_t mod_f_if(mv_t* pa, mv_t* pb) {
	double a = (double)pa->u.intv;
	double b = pb->u.fltv;
	return mv_from_float(a - b * floor(a / b));
}
static mv_t mod_i_ii(mv_t* pa, mv_t* pb) {
	long long a = pa->u.intv;
	long long b = pb->u.intv;
	long long u = a % b;
	// Pythonic division, not C division.
	if (a >= 0LL) {
		if (b < 0LL) {
			u += b;
		}
	} else {
		if (b >= 0LL) {
			u += b;
		}
	}
	return mv_from_int(u);
}

static mv_binary_func_t* mod_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY STRING INT       FLOAT     BOOL
	/*ERROR*/  {_err, _err,  _err, _err,  _err,     _err,     _err},
	/*ABSENT*/ {_err, _a,    _a,   _err,  _i0,      _f0,      _err},
	/*EMPTY*/  {_err, _a,    _emt, _err,  _emt,     _emt,     _err},
	/*STRING*/ {_err, _err,  _err, _err,  _err,     _err,     _err},
	/*INT*/    {_err, _1,    _emt, _err,  mod_i_ii, mod_f_if, _err},
	/*FLOAT*/  {_err, _1,    _emt, _err,  mod_f_fi, mod_f_ff, _err},
	/*BOOL*/   {_err, _err,  _err, _err,  _err,     _err,     _err},
};

mv_t x_xx_mod_func(mv_t* pval1, mv_t* pval2) {
	return (mod_dispositions[pval1->type][pval2->type])(pval1,pval2);
}

// ----------------------------------------------------------------
static mv_t upos_i_i(mv_t* pa) {
	return *pa;
}
static mv_t upos_f_f(mv_t* pa) {
	return *pa;
}

static mv_unary_func_t* upos_dispositions[MT_DIM] = {
	/*ERROR*/  _err1,
	/*ABSENT*/ _a1,
	/*EMPTY*/  _emt1,
	/*STRING*/ _err1,
	/*INT*/    upos_i_i,
	/*FLOAT*/  upos_f_f,
	/*BOOL*/   _err1,
};

mv_t x_x_upos_func(mv_t* pval1) { return (upos_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t uneg_i_i(mv_t* pa) {
	return mv_from_int(-pa->u.intv);
}
static mv_t uneg_f_f(mv_t* pa) {
	return mv_from_float(-pa->u.fltv);
}

static mv_unary_func_t* uneg_disnegitions[MT_DIM] = {
	/*ERROR*/  _err1,
	/*ABSENT*/ _a1,
	/*EMPTY*/  _emt1,
	/*STRING*/ _err1,
	/*INT*/    uneg_i_i,
	/*FLOAT*/  uneg_f_f,
	/*BOOL*/   _err1,
};

mv_t x_x_uneg_func(mv_t* pval1) { return (uneg_disnegitions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t abs_n_f(mv_t* pa) {
	return mv_from_float(fabs(pa->u.fltv));
}
static mv_t abs_n_i(mv_t* pa) {
	return mv_from_int(pa->u.intv < 0LL ? -pa->u.intv : pa->u.intv);
}

static mv_unary_func_t* abs_dispositions[MT_DIM] = {
	/*ERROR*/  _err1,
	/*ABSENT*/ _a1,
	/*EMPTY*/  _emt1,
	/*STRING*/ _err1,
	/*INT*/    abs_n_i,
	/*FLOAT*/  abs_n_f,
	/*BOOL*/   _err1,
};

mv_t x_x_abs_func(mv_t* pval1) { return (abs_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t sgn_n_f(mv_t* pa) {
	if (pa->u.fltv > 0.0)
		return mv_from_float(1.0);
	if (pa->u.fltv < 0.0)
		return mv_from_float(-1.0);
	return mv_from_float(0.0);
}
static mv_t sgn_n_i(mv_t* pa) {
	if (pa->u.intv > 0LL)
		return mv_from_int(1LL);
	if (pa->u.intv < 0LL)
		return mv_from_int(-1LL);
	return mv_from_int(0LL);
}

static mv_unary_func_t* sgn_dispositions[MT_DIM] = {
	/*ERROR*/  _err1,
	/*ABSENT*/ _a1,
	/*EMPTY*/  _emt1,
	/*STRING*/ _err1,
	/*INT*/    sgn_n_i,
	/*FLOAT*/  sgn_n_f,
	/*BOOL*/   _err1,
};

mv_t x_x_sgn_func(mv_t* pval1) { return (sgn_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t ceil_n_f(mv_t* pa) {
	return mv_from_float(ceil(pa->u.fltv));
}
static mv_t ceil_n_i(mv_t* pa) {
	return mv_from_int(pa->u.intv);
}

static mv_unary_func_t* ceil_dispositions[MT_DIM] = {
	/*ERROR*/  _err1,
	/*ABSENT*/ _a1,
	/*EMPTY*/  _emt1,
	/*STRING*/ _err1,
	/*INT*/    ceil_n_i,
	/*FLOAT*/  ceil_n_f,
	/*BOOL*/   _err1,
};

mv_t x_x_ceil_func(mv_t* pval1) { return (ceil_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t floor_n_f(mv_t* pa) {
	return mv_from_float(floor(pa->u.fltv));
}
static mv_t floor_n_i(mv_t* pa) {
	return mv_from_int(pa->u.intv);
}

static mv_unary_func_t* floor_dispositions[MT_DIM] = {
	/*ERROR*/  _err1,
	/*ABSENT*/ _a1,
	/*EMPTY*/  _emt1,
	/*STRING*/ _err1,
	/*INT*/    floor_n_i,
	/*FLOAT*/  floor_n_f,
	/*BOOL*/   _err1,
};

mv_t x_x_floor_func(mv_t* pval1) { return (floor_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t round_n_f(mv_t* pa) {
	return mv_from_float(round(pa->u.fltv));
}
static mv_t round_n_i(mv_t* pa) {
	return mv_from_int(pa->u.intv);
}

static mv_unary_func_t* round_dispositions[MT_DIM] = {
	/*ERROR*/  _err1,
	/*ABSENT*/ _a1,
	/*EMPTY*/  _emt1,
	/*STRING*/ _err1,
	/*INT*/    round_n_i,
	/*FLOAT*/  round_n_f,
	/*BOOL*/   _err1,
};

mv_t x_x_round_func(mv_t* pval1) { return (round_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t roundm_f_ff(mv_t* pa, mv_t* pb) {
	double x = pa->u.fltv;
	double m = pb->u.fltv;
	return mv_from_float(round(x / m) * m);
}
static mv_t roundm_f_fi(mv_t* pa, mv_t* pb) {
	double x = pa->u.fltv;
	double m = (double)pb->u.intv;
	return mv_from_float(round(x / m) * m);
}
static mv_t roundm_f_if(mv_t* pa, mv_t* pb) {
	double x = (double)pa->u.intv;
	double m = pb->u.fltv;
	return mv_from_float(round(x / m) * m);
}
static mv_t roundm_i_ii(mv_t* pa, mv_t* pb) {
	long long x = pa->u.intv;
	long long m = pb->u.intv;
	return mv_from_int((x / m) * m);
}

static mv_binary_func_t* roundm_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY STRING INT          FLOAT        BOOL
	/*ERROR*/  {_err, _err,  _err, _err,  _err,        _err,        _err},
	/*ABSENT*/ {_err, _a,    _a,   _err,  _a,          _err,        _err},
	/*EMPTY*/  {_err, _a,    _emt, _err,  _err,        _err,        _err},
	/*STRING*/ {_err, _err,  _err, _err,  _err,        _err,        _err},
	/*INT*/    {_err, _a,    _err, _err,  roundm_i_ii, roundm_f_if, _err},
	/*FLOAT*/  {_err, _err,  _err, _err,  roundm_f_fi, roundm_f_ff, _err},
	/*BOOL*/   {_err, _err,  _err, _err,  _err,        _err,        _err},
};

mv_t x_xx_roundm_func(mv_t* pval1, mv_t* pval2) { return (roundm_dispositions[pval1->type][pval2->type])(pval1,pval2); }

// ----------------------------------------------------------------
static mv_t min_f_ff(mv_t* pa, mv_t* pb) {
	double a = pa->u.fltv;
	double b = pb->u.fltv;
	return mv_from_float(fmin(a, b));
}

static mv_t min_f_fi(mv_t* pa, mv_t* pb) {
	double a = pa->u.fltv;
	double b = (double)pb->u.intv;
	return mv_from_float(fmin(a, b));
}

static mv_t min_f_if(mv_t* pa, mv_t* pb) {
	double a = (double)pa->u.intv;
	double b = pb->u.fltv;
	return mv_from_float(fmin(a, b));
}

static mv_t min_i_ii(mv_t* pa, mv_t* pb) {
	long long a = pa->u.intv;
	long long b = pb->u.intv;
	return mv_from_int(a < b ? a : b);
}

static mv_binary_func_t* min_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY STRING INT       FLOAT     BOOL
	/*ERROR*/  {_err, _err,  _err, _err,  _err,     _err,     _err},
	/*ABSENT*/ {_err, _a,    _a,   _err,  _2,       _2,       _err},
	/*EMPTY*/  {_err, _a,    _emt, _err,  _2,       _2,       _err},
	/*STRING*/ {_err, _err,  _err, _err,  _err,     _err,     _err},
	/*INT*/    {_err, _1,    _1,   _err,  min_i_ii, min_f_if, _err},
	/*FLOAT*/  {_err, _1,    _1,   _err,  min_f_fi, min_f_ff, _err},
	/*BOOL*/   {_err, _err,  _err, _err,  _err,     _err,     _err},
};

mv_t x_xx_min_func(mv_t* pval1, mv_t* pval2) { return (min_dispositions[pval1->type][pval2->type])(pval1,pval2); }

mv_t variadic_min_func(mv_t* pvals, int nvals) {
	mv_t rv = mv_empty();
	for (int i = 0; i < nvals; i++) {
		rv = x_xx_min_func(&rv, &pvals[i]);
		mv_free(&pvals[i]);
	}
	return rv;
}

// ----------------------------------------------------------------
static mv_t max_f_ff(mv_t* pa, mv_t* pb) {
	double a = pa->u.fltv;
	double b = pb->u.fltv;
	return mv_from_float(fmax(a, b));
}

static mv_t max_f_fi(mv_t* pa, mv_t* pb) {
	double a = pa->u.fltv;
	double b = (double)pb->u.intv;
	return mv_from_float(fmax(a, b));
}

static mv_t max_f_if(mv_t* pa, mv_t* pb) {
	double a = (double)pa->u.intv;
	double b = pb->u.fltv;
	return mv_from_float(fmax(a, b));
}

static mv_t max_i_ii(mv_t* pa, mv_t* pb) {
	long long a = pa->u.intv;
	long long b = pb->u.intv;
	return mv_from_int(a > b ? a : b);
}

static mv_binary_func_t* max_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY STRING INT       FLOAT     BOOL
	/*ERROR*/  {_err, _err,  _err, _err,  _err,     _err,     _err},
	/*ABSENT*/ {_err, _a,    _a,   _err,  _2,       _2,       _err},
	/*EMPTY*/  {_err, _a,    _emt, _err,  _2,       _2,       _err},
	/*STRING*/ {_err, _err,  _err, _err,  _err,     _err,     _err},
	/*INT*/    {_err, _1,    _1,   _err,  max_i_ii, max_f_if, _err},
	/*FLOAT*/  {_err, _1,    _1,   _err,  max_f_fi, max_f_ff, _err},
	/*BOOL*/   {_err, _err,  _err, _err,  _err,     _err,     _err},
};

mv_t x_xx_max_func(mv_t* pval1, mv_t* pval2) { return (max_dispositions[pval1->type][pval2->type])(pval1,pval2); }

mv_t variadic_max_func(mv_t* pvals, int nvals) {
	mv_t rv = mv_empty();
	for (int i = 0; i < nvals; i++) {
		rv = x_xx_max_func(&rv, &pvals[i]);
		mv_free(&pvals[i]);
	}
	return rv;
}

// ----------------------------------------------------------------
static mv_t int_i_b(mv_t* pa) { return mv_from_int(pa->u.boolv ? 1 : 0); }
static mv_t int_i_f(mv_t* pa) { return mv_from_int((long long)round(pa->u.fltv)); }
static mv_t int_i_i(mv_t* pa) { return mv_from_int(pa->u.intv); }
static mv_t int_i_s(mv_t* pa) {
	if (*pa->u.strv == '\0') {
		mv_free(pa);
		return mv_empty();
	}
	mv_t retval = mv_from_int(0LL);
	if (!mlr_try_int_from_string(pa->u.strv, &retval.u.intv)) {
		mv_free(pa);
		return mv_error();
	}
	mv_free(pa);
	return retval;
}

static mv_unary_func_t* int_dispositions[MT_DIM] = {
	/*ERROR*/  _err1,
	/*ABSENT*/ _a1,
	/*EMPTY*/  _emt1,
	/*STRING*/ int_i_s,
	/*INT*/    int_i_i,
	/*FLOAT*/  int_i_f,
	/*BOOL*/   int_i_b,
};

mv_t i_x_int_func(mv_t* pval1) { return (int_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t float_f_b(mv_t* pa) { return mv_from_float(pa->u.boolv ? 1.0 : 0.0); }
static mv_t float_f_f(mv_t* pa) { return mv_from_float(pa->u.fltv); }
static mv_t float_f_i(mv_t* pa) { return mv_from_float((double)pa->u.intv); }
static mv_t float_f_s(mv_t* pa) {
	if (*pa->u.strv == '\0') {
		mv_free(pa);
		return mv_empty();
	}
	mv_t retval = mv_from_float(0.0);
	if (!mlr_try_float_from_string(pa->u.strv, &retval.u.fltv)) {
		mv_free(pa);
		return mv_error();
	}
	mv_free(pa);
	return retval;
}

static mv_unary_func_t* float_dispositions[MT_DIM] = {
	/*ERROR*/  _err1,
	/*ABSENT*/ _a1,
	/*EMPTY*/  _emt1,
	/*STRING*/ float_f_s,
	/*INT*/    float_f_i,
	/*FLOAT*/  float_f_f,
	/*BOOL*/   float_f_b,
};

mv_t f_x_float_func(mv_t* pval1) { return (float_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t band_i_ii(mv_t* pa, mv_t* pb) {
	return mv_from_int(pa->u.intv & pb->u.intv);
}

static mv_binary_func_t* band_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY STRING INT        FLOAT BOOL
	/*ERROR*/  {_err, _err,  _err, _err,  _err,      _err, _err},
	/*ABSENT*/ {_err, _a,    _a,   _err,  _2,        _err, _err},
	/*EMPTY*/  {_err, _a,    _emt, _err,  _emt,      _emt, _err},
	/*STRING*/ {_err, _err,  _err, _err,  _err,      _err, _err},
	/*INT*/    {_err, _1,    _emt, _err,  band_i_ii, _err, _err},
	/*FLOAT*/  {_err, _err,  _emt, _err,  _err,      _err, _err},
	/*BOOL*/   {_err, _err,  _err, _err,  _err,      _err, _err},
};

mv_t x_xx_band_func(mv_t* pval1, mv_t* pval2) { return (band_dispositions[pval1->type][pval2->type])(pval1,pval2); }

// ----------------------------------------------------------------
static mv_t bor_i_ii(mv_t* pa, mv_t* pb) {
	return mv_from_int(pa->u.intv | pb->u.intv);
}

static mv_binary_func_t* bor_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY STRING INT       FLOAT BOOL
	/*ERROR*/  {_err, _err,  _err, _err,  _err,     _err, _err},
	/*ABSENT*/ {_err, _a,    _a,   _err,  _2,       _err, _err},
	/*EMPTY*/  {_err, _a,    _emt, _err,  _emt,     _emt, _err},
	/*STRING*/ {_err, _err,  _err, _err,  _err,     _err, _err},
	/*INT*/    {_err, _1,    _emt, _err,  bor_i_ii, _err, _err},
	/*FLOAT*/  {_err, _err,  _emt, _err,  _err,     _err, _err},
	/*BOOL*/   {_err, _err,  _err, _err,  _err,     _err, _err},
};

mv_t x_xx_bor_func(mv_t* pval1, mv_t* pval2) { return (bor_dispositions[pval1->type][pval2->type])(pval1,pval2); }

// ----------------------------------------------------------------
static mv_t bxor_i_ii(mv_t* pa, mv_t* pb) {
	return mv_from_int(pa->u.intv ^ pb->u.intv);
}

static mv_binary_func_t* bxor_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY STRING INT        FLOAT BOOL
	/*ERROR*/  {_err, _err,  _err, _err,  _err,      _err, _err},
	/*ABSENT*/ {_err, _a,    _a,   _err,  _2,        _err, _err},
	/*EMPTY*/  {_err, _a,    _emt, _err,  _emt,      _emt, _err},
	/*STRING*/ {_err, _err,  _err, _err,  _err,      _err, _err},
	/*INT*/    {_err, _1,    _emt, _err,  bxor_i_ii, _err, _err},
	/*FLOAT*/  {_err, _err,  _emt, _err,  _err,      _err, _err},
	/*BOOL*/   {_err, _err,  _err, _err,  _err,      _err, _err},
};

mv_t x_xx_bxor_func(mv_t* pval1, mv_t* pval2) { return (bxor_dispositions[pval1->type][pval2->type])(pval1,pval2); }

// ----------------------------------------------------------------
static mv_t boolean_b_b(mv_t* pa) { return mv_from_bool(pa->u.boolv); }
static mv_t boolean_b_f(mv_t* pa) { return mv_from_bool((pa->u.fltv == 0.0) ? FALSE : TRUE); }
static mv_t boolean_b_i(mv_t* pa) { return mv_from_bool((pa->u.intv == 0LL) ? FALSE : TRUE); }
static mv_t boolean_b_s(mv_t* pa) { return mv_from_bool((streq(pa->u.strv, "true") || streq(pa->u.strv, "TRUE")) ? TRUE : FALSE);}

static mv_unary_func_t* boolean_dispositions[MT_DIM] = {
	/*ERROR*/  _err1,
	/*ABSENT*/ _a1,
	/*EMPTY*/  _emt1,
	/*STRING*/ boolean_b_s,
	/*INT*/    boolean_b_i,
	/*FLOAT*/  boolean_b_f,
	/*BOOL*/   boolean_b_b,
};

mv_t b_x_boolean_func(mv_t* pval1) { return (boolean_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t string_s_b(mv_t* pa) { return mv_from_string_no_free(pa->u.boolv?"true":"false"); }
static mv_t string_s_f(mv_t* pa) { return mv_from_string_with_free(mlr_alloc_string_from_double(pa->u.fltv, MLR_GLOBALS.ofmt)); }
static mv_t string_s_i(mv_t* pa) { return mv_from_string_with_free(mlr_alloc_string_from_ll(pa->u.intv)); }
static mv_t string_s_s(mv_t* pa) {
	char free_flags = pa->free_flags;
	pa->free_flags = NO_FREE;
	return mv_from_string(pa->u.strv, free_flags);
}

static mv_unary_func_t* string_dispositions[MT_DIM] = {
	/*ERROR*/  _err1,
	/*ABSENT*/ _a1,
	/*EMPTY*/  _emt1,
	/*STRING*/ string_s_s,
	/*INT*/    string_s_i,
	/*FLOAT*/  string_s_f,
	/*BOOL*/   string_s_b,
};

mv_t s_x_string_func(mv_t* pval1) { return (string_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
mv_t s_sii_substr_func(mv_t* pval1, mv_t* pval2, mv_t* pval3) {
	int m = pval2->u.intv; // inclusive lower; -len..-1 alias to 0..len-1
	int n = pval3->u.intv; // inclusive upper; -len..-1 alias to 0..len-1
	int len = strlen(pval1->u.strv);
	mv_t rv;

	if (m < 0)
		m = len + m;
	if (n < 0)
		n = len + n;

	if (m < 0 || m >= len || n < 0 || n >= len || n < m) {
		rv = mv_from_string("", 0);
	} else {
		int olen = n - m + 1;
		char* p = mlr_malloc_or_die(olen + 1);
		strncpy(p, &pval1->u.strv[m], olen);
		p[olen] = 0;
		rv = mv_from_string(p, FREE_ENTRY_VALUE);
	}
	mv_free(pval1);
	mv_free(pval2);
	mv_free(pval3);
	return rv;
}

// ----------------------------------------------------------------
static mv_t hexfmt_s_b(mv_t* pa) { return mv_from_string_no_free(pa->u.boolv?"0x1":"0x0"); }
static mv_t hexfmt_s_f(mv_t* pa) { return mv_from_string_with_free(mlr_alloc_hexfmt_from_ll((long long)pa->u.fltv)); }
static mv_t hexfmt_s_i(mv_t* pa) { return mv_from_string_with_free(mlr_alloc_hexfmt_from_ll(pa->u.intv)); }
static mv_t hexfmt_s_s(mv_t* pa) {
	char free_flags = pa->free_flags;
	pa->free_flags = NO_FREE;
	return mv_from_string(pa->u.strv, free_flags);
}

static mv_unary_func_t* hexfmt_dispositions[MT_DIM] = {
	/*ERROR*/  _err1,
	/*ABSENT*/ _a1,
	/*EMPTY*/  _emt1,
	/*STRING*/ hexfmt_s_s,
	/*INT*/    hexfmt_s_i,
	/*FLOAT*/  hexfmt_s_f,
	/*BOOL*/   hexfmt_s_b,
};

mv_t s_x_hexfmt_func(mv_t* pval1) { return (hexfmt_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t fmtnum_s_bs(mv_t* pa, mv_t* pfmt) { return mv_from_string_no_free(pa->u.boolv?"0x1":"0x0"); }
static mv_t fmtnum_s_ds(mv_t* pa, mv_t* pfmt) {
	mv_t rv = mv_from_string_with_free(mlr_alloc_string_from_double(pa->u.fltv, pfmt->u.strv));
	mv_free(pfmt);
	return rv;
}
static mv_t fmtnum_s_is(mv_t* pa, mv_t* pfmt) {
	mv_t rv = mv_from_string_with_free(mlr_alloc_string_from_ll_and_format(pa->u.intv, pfmt->u.strv));
	mv_free(pfmt);
	return rv;
}
static mv_binary_func_t* fmtnum_dispositions[MT_DIM] = {
	/*ERROR*/  _err,
	/*ABSENT*/ _a,
	/*EMPTY*/  _emt,
	/*STRING*/ _err,
	/*INT*/    fmtnum_s_is,
	/*FLOAT*/  fmtnum_s_ds,
	/*BOOL*/   fmtnum_s_bs,
};

mv_t s_xs_fmtnum_func(mv_t* pval1, mv_t* pval2) { return (fmtnum_dispositions[pval1->type])(pval1, pval2); }

// ----------------------------------------------------------------
static mv_t eq_b_ii(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv == pb->u.intv); }
static mv_t ne_b_ii(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv != pb->u.intv); }
static mv_t gt_b_ii(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv >  pb->u.intv); }
static mv_t ge_b_ii(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv >= pb->u.intv); }
static mv_t lt_b_ii(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv <  pb->u.intv); }
static mv_t le_b_ii(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv <= pb->u.intv); }

static mv_t eq_b_ff(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv == pb->u.fltv); }
static mv_t ne_b_ff(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv != pb->u.fltv); }
static mv_t gt_b_ff(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv >  pb->u.fltv); }
static mv_t ge_b_ff(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv >= pb->u.fltv); }
static mv_t lt_b_ff(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv <  pb->u.fltv); }
static mv_t le_b_ff(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv <= pb->u.fltv); }

static mv_t eq_b_fi(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv == pb->u.intv); }
static mv_t ne_b_fi(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv != pb->u.intv); }
static mv_t gt_b_fi(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv >  pb->u.intv); }
static mv_t ge_b_fi(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv >= pb->u.intv); }
static mv_t lt_b_fi(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv <  pb->u.intv); }
static mv_t le_b_fi(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv <= pb->u.intv); }

static mv_t eq_b_if(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv == pb->u.fltv); }
static mv_t ne_b_if(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv != pb->u.fltv); }
static mv_t gt_b_if(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv >  pb->u.fltv); }
static mv_t ge_b_if(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv >= pb->u.fltv); }
static mv_t lt_b_if(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv <  pb->u.fltv); }
static mv_t le_b_if(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv <= pb->u.fltv); }

static mv_t eq_b_xs(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sa = mv_format_val(pa, &free_flags);
	mv_t rv = mv_from_bool(strcmp(sa, pb->u.strv) == 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sa);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static mv_t ne_b_xs(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sa = mv_format_val(pa, &free_flags);
	mv_t rv = mv_from_bool(strcmp(sa, pb->u.strv) != 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sa);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static mv_t gt_b_xs(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sa = mv_format_val(pa, &free_flags);
	mv_t rv = mv_from_bool(strcmp(sa, pb->u.strv) > 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sa);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static mv_t ge_b_xs(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sa = mv_format_val(pa, &free_flags);
	mv_t rv = mv_from_bool(strcmp(sa, pb->u.strv) >= 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sa);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static mv_t lt_b_xs(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sa = mv_format_val(pa, &free_flags);
	mv_t rv = mv_from_bool(strcmp(sa, pb->u.strv) < 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sa);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static mv_t le_b_xs(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sa = mv_format_val(pa, &free_flags);
	mv_t rv = mv_from_bool(strcmp(sa, pb->u.strv) <= 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sa);
	mv_free(pa);
	mv_free(pb);
	return rv;
}

static mv_t eq_b_sx(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sb = mv_format_val(pb, &free_flags);
	mv_t rv = mv_from_bool(strcmp(pa->u.strv, sb) == 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sb);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static mv_t ne_b_sx(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sb = mv_format_val(pb, &free_flags);
	mv_t rv = mv_from_bool(strcmp(pa->u.strv, sb) != 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sb);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static mv_t gt_b_sx(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sb = mv_format_val(pb, &free_flags);
	mv_t rv = mv_from_bool(strcmp(pa->u.strv, sb) > 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sb);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static mv_t ge_b_sx(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sb = mv_format_val(pb, &free_flags);
	mv_t rv = mv_from_bool(strcmp(pa->u.strv, sb) >= 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sb);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static mv_t lt_b_sx(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sb = mv_format_val(pb, &free_flags);
	mv_t rv = mv_from_bool(strcmp(pa->u.strv, sb) < 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sb);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static mv_t le_b_sx(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sb = mv_format_val(pb, &free_flags);
	mv_t rv = mv_from_bool(strcmp(pa->u.strv, sb) <= 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sb);
	mv_free(pa);
	mv_free(pb);
	return rv;
}

static mv_t eq_b_ss(mv_t*pa, mv_t*pb) {
	mv_t rv = mv_from_bool(strcmp(pa->u.strv, pb->u.strv) == 0);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static mv_t ne_b_ss(mv_t*pa, mv_t*pb) {
	mv_t rv = mv_from_bool(strcmp(pa->u.strv, pb->u.strv) != 0);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static mv_t gt_b_ss(mv_t*pa, mv_t*pb) {
	mv_t rv = mv_from_bool(strcmp(pa->u.strv, pb->u.strv) >  0);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static mv_t ge_b_ss(mv_t*pa, mv_t*pb) {
	mv_t rv = mv_from_bool(strcmp(pa->u.strv, pb->u.strv) >= 0);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static mv_t lt_b_ss(mv_t*pa, mv_t*pb) {
	mv_t rv = mv_from_bool(strcmp(pa->u.strv, pb->u.strv) <  0);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static mv_t le_b_ss(mv_t*pa, mv_t*pb) {
	mv_t rv = mv_from_bool(strcmp(pa->u.strv, pb->u.strv) <= 0);
	mv_free(pa);
	mv_free(pb);
	return rv;
}

static mv_binary_func_t* eq_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY    STRING   INT      FLOAT    BOOL
	/*ERROR*/  {_err, _err,  _err,    _err,    _err,    _err,    _err},
	/*ABSENT*/ {_err, _a,    _a,      _a,      _a,      _a,      _a},
	/*EMPTY*/  {_err, _a,    eq_b_ss, eq_b_ss, eq_b_sx, eq_b_sx, _err},
	/*STRING*/ {_err, _a,    eq_b_ss, eq_b_ss, eq_b_sx, eq_b_sx, _err},
	/*INT*/    {_err, _a,    eq_b_xs, eq_b_xs, eq_b_ii, eq_b_if, _err},
	/*FLOAT*/  {_err, _a,    eq_b_xs, eq_b_xs, eq_b_fi, eq_b_ff, _err},
	/*BOOL*/   {_err, _err,  _a,      _err,    _err,    _err,    _err},
};

static mv_binary_func_t* ne_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY    STRING   INT      FLOAT    BOOL
	/*ERROR*/  {_err, _err,  _err,    _err,    _err,    _err,    _err},
	/*ABSENT*/ {_err, _a,    _a,      _a,      _a,      _a,      _a},
	/*EMPTY*/  {_err, _a,    ne_b_ss, ne_b_ss, ne_b_sx, ne_b_sx, _err},
	/*STRING*/ {_err, _a,    ne_b_ss, ne_b_ss, ne_b_sx, ne_b_sx, _err},
	/*INT*/    {_err, _a,    ne_b_xs, ne_b_xs, ne_b_ii, ne_b_if, _err},
	/*FLOAT*/  {_err, _a,    ne_b_xs, ne_b_xs, ne_b_fi, ne_b_ff, _err},
	/*BOOL*/   {_err, _err,  _a,      _err,    _err,    _err,    _err},
};

static mv_binary_func_t* gt_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY    STRING   INT      FLOAT    BOOL
	/*ERROR*/  {_err, _err,  _err,    _err,    _err,    _err,    _err},
	/*ABSENT*/ {_err, _a,    _a,      _a,      _a,      _a,      _a},
	/*EMPTY*/  {_err, _a,    gt_b_ss, gt_b_ss, gt_b_sx, gt_b_sx, _err},
	/*STRING*/ {_err, _a,    gt_b_ss, gt_b_ss, gt_b_sx, gt_b_sx, _err},
	/*INT*/    {_err, _a,    gt_b_xs, gt_b_xs, gt_b_ii, gt_b_if, _err},
	/*FLOAT*/  {_err, _a,    gt_b_xs, gt_b_xs, gt_b_fi, gt_b_ff, _err},
	/*BOOL*/   {_err, _err,  _a,      _err,    _err,    _err,    _err},
};

static mv_binary_func_t* ge_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY    STRING   INT      FLOAT    BOOL
	/*ERROR*/  {_err, _err,  _err,    _err,    _err,    _err,    _err},
	/*ABSENT*/ {_err, _a,    _a,      _a,      _a,      _a,      _a},
	/*EMPTY*/  {_err, _a,    ge_b_ss, ge_b_ss, ge_b_sx, ge_b_sx, _err},
	/*STRING*/ {_err, _a,    ge_b_ss, ge_b_ss, ge_b_sx, ge_b_sx, _err},
	/*INT*/    {_err, _a,    ge_b_xs, ge_b_xs, ge_b_ii, ge_b_if, _err},
	/*FLOAT*/  {_err, _a,    ge_b_xs, ge_b_xs, ge_b_fi, ge_b_ff, _err},
	/*BOOL*/   {_err, _err,  _a,      _err,    _err,    _err,    _err},
};

static mv_binary_func_t* lt_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY    STRING   INT      FLOAT    BOOL
	/*ERROR*/  {_err, _err,  _err,    _err,    _err,    _err,    _err},
	/*ABSENT*/ {_err, _a,    _a,      _a,      _a,      _a,      _a},
	/*EMPTY*/  {_err, _a,    lt_b_ss, lt_b_ss, lt_b_sx, lt_b_sx, _err},
	/*STRING*/ {_err, _a,    lt_b_ss, lt_b_ss, lt_b_sx, lt_b_sx, _err},
	/*INT*/    {_err, _a,    lt_b_xs, lt_b_xs, lt_b_ii, lt_b_if, _err},
	/*FLOAT*/  {_err, _a,    lt_b_xs, lt_b_xs, lt_b_fi, lt_b_ff, _err},
	/*BOOL*/   {_err, _err,  _a,      _err,    _err,    _err,    _err},
};

static mv_binary_func_t* le_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY    STRING   INT      FLOAT    BOOL
	/*ERROR*/  {_err, _err,  _err,    _err,    _err,    _err,    _err},
	/*ABSENT*/ {_err, _a,    _a,      _a,      _a,      _a,      _a},
	/*EMPTY*/  {_err, _a,    le_b_ss, le_b_ss, le_b_sx, le_b_sx, _err},
	/*STRING*/ {_err, _a,    le_b_ss, le_b_ss, le_b_sx, le_b_sx, _err},
	/*INT*/    {_err, _a,    le_b_xs, le_b_xs, le_b_ii, le_b_if, _err},
	/*FLOAT*/  {_err, _a,    le_b_xs, le_b_xs, le_b_fi, le_b_ff, _err},
	/*BOOL*/   {_err, _err,  _a,      _err,    _err,    _err,    _err},
};

mv_t eq_op_func(mv_t* pval1, mv_t* pval2) { return (eq_dispositions[pval1->type][pval2->type])(pval1, pval2); }
mv_t ne_op_func(mv_t* pval1, mv_t* pval2) { return (ne_dispositions[pval1->type][pval2->type])(pval1, pval2); }
mv_t gt_op_func(mv_t* pval1, mv_t* pval2) { return (gt_dispositions[pval1->type][pval2->type])(pval1, pval2); }
mv_t ge_op_func(mv_t* pval1, mv_t* pval2) { return (ge_dispositions[pval1->type][pval2->type])(pval1, pval2); }
mv_t lt_op_func(mv_t* pval1, mv_t* pval2) { return (lt_dispositions[pval1->type][pval2->type])(pval1, pval2); }
mv_t le_op_func(mv_t* pval1, mv_t* pval2) { return (le_dispositions[pval1->type][pval2->type])(pval1, pval2); }

// ----------------------------------------------------------------
int mv_equals_si(mv_t* pa, mv_t* pb) {
	if (pa->type == MT_INT) {
		return (pb->type == MT_INT) ? pa->u.intv == pb->u.intv : FALSE;
	} else {
		return (pb->type == MT_STRING) ? streq(pa->u.strv, pb->u.strv) : FALSE;
	}
}

// ----------------------------------------------------------------
static int eq_i_ii(mv_t* pa, mv_t* pb) { return  pa->u.intv == pb->u.intv; }
static int ne_i_ii(mv_t* pa, mv_t* pb) { return  pa->u.intv != pb->u.intv; }
static int gt_i_ii(mv_t* pa, mv_t* pb) { return  pa->u.intv >  pb->u.intv; }
static int ge_i_ii(mv_t* pa, mv_t* pb) { return  pa->u.intv >= pb->u.intv; }
static int lt_i_ii(mv_t* pa, mv_t* pb) { return  pa->u.intv <  pb->u.intv; }
static int le_i_ii(mv_t* pa, mv_t* pb) { return  pa->u.intv <= pb->u.intv; }

static int eq_i_ff(mv_t* pa, mv_t* pb) { return  pa->u.fltv == pb->u.fltv; }
static int ne_i_ff(mv_t* pa, mv_t* pb) { return  pa->u.fltv != pb->u.fltv; }
static int gt_i_ff(mv_t* pa, mv_t* pb) { return  pa->u.fltv >  pb->u.fltv; }
static int ge_i_ff(mv_t* pa, mv_t* pb) { return  pa->u.fltv >= pb->u.fltv; }
static int lt_i_ff(mv_t* pa, mv_t* pb) { return  pa->u.fltv <  pb->u.fltv; }
static int le_i_ff(mv_t* pa, mv_t* pb) { return  pa->u.fltv <= pb->u.fltv; }

static int eq_i_fi(mv_t* pa, mv_t* pb) { return  pa->u.fltv == pb->u.intv; }
static int ne_i_fi(mv_t* pa, mv_t* pb) { return  pa->u.fltv != pb->u.intv; }
static int gt_i_fi(mv_t* pa, mv_t* pb) { return  pa->u.fltv >  pb->u.intv; }
static int ge_i_fi(mv_t* pa, mv_t* pb) { return  pa->u.fltv >= pb->u.intv; }
static int lt_i_fi(mv_t* pa, mv_t* pb) { return  pa->u.fltv <  pb->u.intv; }
static int le_i_fi(mv_t* pa, mv_t* pb) { return  pa->u.fltv <= pb->u.intv; }

static int eq_i_if(mv_t* pa, mv_t* pb) { return  pa->u.intv == pb->u.fltv; }
static int ne_i_if(mv_t* pa, mv_t* pb) { return  pa->u.intv != pb->u.fltv; }
static int gt_i_if(mv_t* pa, mv_t* pb) { return  pa->u.intv >  pb->u.fltv; }
static int ge_i_if(mv_t* pa, mv_t* pb) { return  pa->u.intv >= pb->u.fltv; }
static int lt_i_if(mv_t* pa, mv_t* pb) { return  pa->u.intv <  pb->u.fltv; }
static int le_i_if(mv_t* pa, mv_t* pb) { return  pa->u.intv <= pb->u.fltv; }

static mv_i_nn_comparator_func_t* ieq_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY STRING INT      FLOAT    BOOL
	/*ERROR*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*ABSENT*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*EMPTY*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*STRING*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*INT*/    {NULL, NULL,  NULL, NULL,  eq_i_ii, eq_i_if, NULL},
	/*FLOAT*/  {NULL, NULL,  NULL, NULL,  eq_i_fi, eq_i_ff, NULL},
	/*BOOL*/   {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
};

static mv_i_nn_comparator_func_t* ine_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY STRING INT      FLOAT    BOOL
	/*ERROR*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*ABSENT*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*EMPTY*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*STRING*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*INT*/    {NULL, NULL,  NULL, NULL,  ne_i_ii, ne_i_if, NULL},
	/*FLOAT*/  {NULL, NULL,  NULL, NULL,  ne_i_fi, ne_i_ff, NULL},
	/*BOOL*/   {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
};

static mv_i_nn_comparator_func_t* igt_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY STRING INT      FLOAT    BOOL
	/*ERROR*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*ABSENT*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*EMPTY*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*STRING*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*INT*/    {NULL, NULL,  NULL, NULL,  gt_i_ii, gt_i_if, NULL},
	/*FLOAT*/  {NULL, NULL,  NULL, NULL,  gt_i_fi, gt_i_ff, NULL},
	/*BOOL*/   {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
};

static mv_i_nn_comparator_func_t* ige_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY STRING INT      FLOAT    BOOL
	/*ERROR*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*ABSENT*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*EMPTY*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*STRING*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*INT*/    {NULL, NULL,  NULL, NULL,  ge_i_ii, ge_i_if, NULL},
	/*FLOAT*/  {NULL, NULL,  NULL, NULL,  ge_i_fi, ge_i_ff, NULL},
	/*BOOL*/   {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
};

static mv_i_nn_comparator_func_t* ilt_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY STRING INT      FLOAT    BOOL
	/*ERROR*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*ABSENT*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*EMPTY*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*STRING*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*INT*/    {NULL, NULL,  NULL, NULL,  lt_i_ii, lt_i_if, NULL},
	/*FLOAT*/  {NULL, NULL,  NULL, NULL,  lt_i_fi, lt_i_ff, NULL},
	/*BOOL*/   {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
};

static mv_i_nn_comparator_func_t* ile_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY STRING INT      FLOAT    BOOL
	/*ERROR*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*ABSENT*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*EMPTY*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*STRING*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*INT*/    {NULL, NULL,  NULL, NULL,  le_i_ii, le_i_if, NULL},
	/*FLOAT*/  {NULL, NULL,  NULL, NULL,  le_i_fi, le_i_ff, NULL},
	/*BOOL*/   {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
};

int mv_i_nn_eq(mv_t* pval1, mv_t* pval2) { return (ieq_dispositions[pval1->type][pval2->type])(pval1, pval2); }
int mv_i_nn_ne(mv_t* pval1, mv_t* pval2) { return (ine_dispositions[pval1->type][pval2->type])(pval1, pval2); }
int mv_i_nn_gt(mv_t* pval1, mv_t* pval2) { return (igt_dispositions[pval1->type][pval2->type])(pval1, pval2); }
int mv_i_nn_ge(mv_t* pval1, mv_t* pval2) { return (ige_dispositions[pval1->type][pval2->type])(pval1, pval2); }
int mv_i_nn_lt(mv_t* pval1, mv_t* pval2) { return (ilt_dispositions[pval1->type][pval2->type])(pval1, pval2); }
int mv_i_nn_le(mv_t* pval1, mv_t* pval2) { return (ile_dispositions[pval1->type][pval2->type])(pval1, pval2); }

// ----------------------------------------------------------------
// For unit-test keystroke-saving

int mveq(mv_t* pval1, mv_t* pval2) {
	mv_t cmp = eq_op_func(pval1, pval2);
	MLR_INTERNAL_CODING_ERROR_UNLESS(cmp.type == MT_BOOLEAN);
	return cmp.u.boolv;
}

int mvne(mv_t* pval1, mv_t* pval2) {
	return !mveq(pval1, pval2);
}

int mveqcopy(mv_t* pval1, mv_t* pval2) {
	mv_t c1 = mv_copy(pval1);
	mv_t c2 = mv_copy(pval2);
	return mveq(&c1, &c2);
}

int mvnecopy(mv_t* pval1, mv_t* pval2) {
	mv_t c1 = mv_copy(pval1);
	mv_t c2 = mv_copy(pval2);
	return mvne(&c1, &c2);
}

// ----------------------------------------------------------------
// arg2 evaluates to string via compound expression; regexes compiled on each call.
mv_t matches_no_precomp_func(mv_t* pval1, mv_t* pval2, string_array_t** ppregex_captures) {
	char* s1 = pval1->u.strv;
	char* s2 = pval2->u.strv;

	regex_t regex;
	char* sstr   = s1;
	char* sregex = s2;

	regcomp_or_die(&regex, sregex, REG_NOSUB);

	const size_t nmatchmax = 10; // Capture-groups \1 through \9 supported, along with entire-string match
	regmatch_t matches[nmatchmax];
	if (regmatch_or_die(&regex, sstr, nmatchmax, matches)) {
		if (ppregex_captures != NULL && *ppregex_captures != NULL)
			save_regex_captures(ppregex_captures, pval1->u.strv, matches, nmatchmax);
		regfree(&regex);
		mv_free(pval1);
		mv_free(pval2);
		return mv_from_true();
	} else {
		regfree(&regex);
		mv_free(pval1);
		mv_free(pval2);
		return mv_from_false();
	}
}

mv_t does_not_match_no_precomp_func(mv_t* pval1, mv_t* pval2, string_array_t** ppregex_captures) {
	mv_t rv = matches_no_precomp_func(pval1, pval2, ppregex_captures);
	rv.u.boolv = !rv.u.boolv;
	return rv;
}

// ----------------------------------------------------------------
// arg2 is a string, compiled to regex only once at alloc time
mv_t matches_precomp_func(mv_t* pval1, regex_t* pregex, string_builder_t* psb, string_array_t** ppregex_captures) {
	const size_t nmatchmax = 10; // Capture-groups \1 through \9 supported, along with entire-string match
	regmatch_t matches[nmatchmax];
	if (regmatch_or_die(pregex, pval1->u.strv, nmatchmax, matches)) {
		if (ppregex_captures != NULL)
			save_regex_captures(ppregex_captures, pval1->u.strv, matches, nmatchmax);
		mv_free(pval1);
		return mv_from_true();
	} else {
		// See comments in mapper_put.c. Setting this array to length 0 (i.e. zero matches) signals to the
		// lrec-evaluator's from-literal function that we *are* in a regex-match context but there are *no* matches to
		// be interpolated.
		if (ppregex_captures != NULL) {
			if (*ppregex_captures != NULL)
				string_array_realloc(*ppregex_captures, 0);
			else
				*ppregex_captures = string_array_alloc(0);
		}
		mv_free(pval1);
		return mv_from_false();
	}
}

mv_t does_not_match_precomp_func(mv_t* pval1, regex_t* pregex, string_builder_t* psb, string_array_t** ppregex_captures) {
	mv_t rv = matches_precomp_func(pval1, pregex, psb, ppregex_captures);
	rv.u.boolv = !rv.u.boolv;
	return rv;
}

// ----------------------------------------------------------------
static int mv_ff_comparator(const mv_t* pa, const mv_t* pb) {
	double d = pa->u.fltv - pb->u.fltv;
	return (d < 0) ? -1 : (d > 0) ? 1 : 0;
}
static int mv_fi_comparator(const mv_t* pa, const mv_t* pb) {
	double d = pa->u.fltv - pb->u.intv;
	return (d < 0) ? -1 : (d > 0) ? 1 : 0;
}
static int mv_if_comparator(const mv_t* pa, const mv_t* pb) {
	double d = pa->u.intv - pb->u.fltv;
	return (d < 0) ? -1 : (d > 0) ? 1 : 0;
}
static int mv_ii_comparator(const mv_t* pa, const mv_t* pb) {
	long long d = pa->u.intv - pb->u.intv;
	return (d < 0) ? -1 : (d > 0) ? 1 : 0;
}
// We assume mv_t's coming into percentile keeper are int or double -- in particular, non-null.
static mv_i_cncn_comparator_func_t* mv_comparator_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT EMPTY STRING INT               FLOAT             BOOL
	/*ERROR*/  {NULL, NULL,  NULL, NULL,  NULL,             NULL,             NULL},
	/*ABSENT*/ {NULL, NULL,  NULL, NULL,  NULL,             NULL,             NULL},
	/*EMPTY*/  {NULL, NULL,  NULL, NULL,  NULL,             NULL,             NULL},
	/*STRING*/ {NULL, NULL,  NULL, NULL,  NULL,             NULL,             NULL},
	/*INT*/    {NULL, NULL,  NULL, NULL,  mv_ii_comparator, mv_if_comparator, NULL},
	/*FLOAT*/  {NULL, NULL,  NULL, NULL,  mv_fi_comparator, mv_ff_comparator, NULL},
	/*BOOL*/   {NULL, NULL,  NULL, NULL,  NULL,             NULL,             NULL},
};

int mv_nn_comparator(const void* pva, const void* pvb) {
	const mv_t* pa = pva;
	const mv_t* pb = pvb;
	return mv_comparator_dispositions[pa->type][pb->type](pa, pb);
}

// ----------------------------------------------------------------
int mlr_bsearch_mv_n_for_insert(mv_t* array, int size, mv_t* pvalue) {
	int lo = 0;
	int hi = size-1;
	int mid = (hi+lo)/2;
	int newmid;

	if (size == 0)
		return 0;
	if (mv_i_nn_gt(pvalue, &array[0]))
		return 0;
	if (mv_i_nn_lt(pvalue, &array[hi]))
		return size;

	while (lo < hi) {
		mv_t* pa = &array[mid];
		if (mv_i_nn_eq(pvalue, pa)) {
			return mid;
		}
		else if (mv_i_nn_gt(pvalue, pa)) {
			hi = mid;
			newmid = (hi+lo)/2;
		}
		else {
			lo = mid;
			newmid = (hi+lo)/2;
		}
		if (mid == newmid) {
			if (mv_i_nn_ge(pvalue, &array[lo]))
				return lo;
			else if (mv_i_nn_ge(pvalue, &array[hi]))
				return hi;
			else
				return hi+1;
		}
		mid = newmid;
	}

	return lo;
}
