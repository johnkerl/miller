#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include <time.h>
#include <regex.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/mlrregex.h"
#include "containers/mlrval.h"

// ================================================================
// See important notes at the top of mlrval.h.
// ================================================================

// For some Linux distros, in spite of including time.h:
char *strptime(const char *s, const char *format, struct tm *ptm);

typedef int mv_i_nn_comparator_func_t(mv_t* pa, mv_t* pb);
typedef int mv_i_cncn_comparator_func_t(const mv_t* pa, const mv_t* pb);

// ----------------------------------------------------------------
char* mt_describe_type(int type) {
	switch (type) {
	case MT_ERROR:  return "MT_ERROR";  break;
	case MT_ABSENT: return "MT_ABSENT"; break;
	case MT_UNINIT: return "MT_UNINIT"; break;
	case MT_VOID:   return "MT_VOID";  break;
	case MT_STRING: return "MT_STRING"; break;
	case MT_INT:    return "MT_INT";    break;
	case MT_FLOAT:  return "MT_FLOAT";  break;
	case MT_BOOL:   return "MT_BOOL";   break;
	default:        return "???";       break;
	}
}

// The caller should free the return value
char* mv_alloc_format_val(mv_t* pval) {
	switch(pval->type) {
	case MT_ABSENT:
	case MT_VOID:
	case MT_UNINIT:
		return mlr_strdup_or_die("");
		break;
	case MT_ERROR:
		return mlr_strdup_or_die("(error)");
		break;
	case MT_BOOL:
		return mlr_strdup_or_die(pval->u.boolv ? "true" : "false");
		break;
	case MT_FLOAT:
		return mlr_alloc_string_from_double(pval->u.fltv, MLR_GLOBALS.ofmt);
		break;
	case MT_INT:
		return mlr_alloc_string_from_ll(pval->u.intv);
		break;
	case MT_STRING:
		return mlr_strdup_or_die(pval->u.strv);
		break;
	default:
		return mlr_strdup_or_die("???");
		break;
	}
}

// The caller should free the return value if the free-flag is non-zero after return
char* mv_format_val(mv_t* pval, char* pfree_flags) {
	char* rv = NULL;
	switch(pval->type) {
	case MT_ABSENT:
	case MT_VOID:
	case MT_UNINIT:
		*pfree_flags = NO_FREE;
		rv = "";
		break;
	case MT_ERROR:
		*pfree_flags = NO_FREE;
		rv = "(error)";
		break;
	case MT_BOOL:
		*pfree_flags = NO_FREE;
		rv = pval->u.boolv ? "true" : "false";
		break;
	case MT_FLOAT:
		*pfree_flags = FREE_ENTRY_VALUE;
		rv = mlr_alloc_string_from_double(pval->u.fltv, MLR_GLOBALS.ofmt);
		break;
	case MT_INT:
		*pfree_flags = FREE_ENTRY_VALUE;
		rv = mlr_alloc_string_from_ll(pval->u.intv);
		break;
	case MT_STRING:
		// Ownership transfer to the caller
		*pfree_flags = pval->free_flags;;
		rv = pval->u.strv;
		*pval = mv_void();
		break;
	default:
		*pfree_flags = NO_FREE;
		rv = "???";
		break;
	}
	return rv;
}

char* mv_describe_val(mv_t val) {
	char* stype = mt_describe_type(val.type);
	char* strv  = mv_alloc_format_val(&val);
	char* desc  = mlr_malloc_or_die(strlen(stype) + strlen(strv) + 4);
	sprintf(desc, "[%s] %s", stype, strv);
	free(strv);
	return desc;
}

// ----------------------------------------------------------------
void mv_set_boolean_strict(mv_t* pval) {
	if (pval->type != MT_BOOL) {
		char* desc = mt_describe_type(pval->type);
		fprintf(stderr, "Expression does not evaluate to boolean: got %s.\n", desc);
		exit(1);
	}
}

// ----------------------------------------------------------------
void mv_set_float_strict(mv_t* pval) {
	double fltv = 0.0;
	mv_t nval = mv_error();
	switch (pval->type) {
	case MT_ABSENT:
	case MT_VOID:
	case MT_UNINIT:
		break;
	case MT_ERROR:
		break;
	case MT_FLOAT:
		break;
	case MT_STRING:
		if (!mlr_try_float_from_string(pval->u.strv, &fltv)) {
			// keep nval = mv_error()
		} else {
			nval = mv_from_float(fltv);
		}
		mv_free(pval);
		*pval = nval;
		break;
	case MT_INT:
		pval ->type = MT_FLOAT;
		pval->u.fltv = (double)pval->u.intv;
		break;
	case MT_BOOL:
		pval->type = MT_ERROR;
		pval->u.intv = 0LL;
		break;
	default:
		fprintf(stderr, "%s: internal coding error detected at file %s, line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		break;
	}
}

// ----------------------------------------------------------------
void mv_set_float_nullable(mv_t* pval) {
	double fltv = 0.0;
	mv_t nval = mv_error();
	switch (pval->type) {
	case MT_ABSENT:
	case MT_VOID:
	case MT_UNINIT:
		break;
	case MT_ERROR:
		break;
	case MT_FLOAT:
		break;
	case MT_INT:
		pval ->type = MT_FLOAT;
		pval->u.fltv = (double)pval->u.intv;
		break;
	case MT_BOOL:
		pval->type = MT_ERROR;
		pval->u.intv = 0;
		break;
	case MT_STRING:
		if (*pval->u.strv == '\0') {
			nval = mv_void();
		} else if (!mlr_try_float_from_string(pval->u.strv, &fltv)) {
			// keep nval = mv_error()
		} else {
			nval = mv_from_float(fltv);
		}
		mv_free(pval);
		*pval = nval;
		break;
	default:
		fprintf(stderr, "%s: internal coding error detected at file %s, line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		break;
	}
}

// ----------------------------------------------------------------
void mv_set_int_nullable(mv_t* pval) {
	long long intv = 0LL;
	mv_t nval = mv_error();
	switch (pval->type) {
	case MT_ABSENT:
	case MT_VOID:
	case MT_UNINIT:
		break;
	case MT_ERROR:
		break;
	case MT_INT:
		break;
	case MT_FLOAT:
		pval ->type = MT_INT;
		pval->u.intv = (long long)pval->u.fltv;
		break;
	case MT_BOOL:
		pval->type = MT_ERROR;
		pval->u.intv = 0;
		break;
	case MT_STRING:
		if (*pval->u.strv == '\0') {
			nval = mv_void();
		} else if (!mlr_try_int_from_string(pval->u.strv, &intv)) {
			// keep nval = mv_error()
		} else {
			nval = mv_from_int(intv);
		}
		mv_free(pval);
		*pval = nval;
		break;
	default:
		fprintf(stderr, "%s: internal coding error detected at file %s, line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		break;
	}
}

// ----------------------------------------------------------------
void mv_set_number_nullable(mv_t* pval) {
	mv_t nval = mv_void();
	switch (pval->type) {
	case MT_ABSENT:
	case MT_VOID:
	case MT_UNINIT:
		break;
	case MT_ERROR:
		break;
	case MT_INT:
		break;
	case MT_FLOAT:
		break;
	case MT_BOOL:
		pval->type = MT_ERROR;
		pval->u.intv = 0;
		break;
	case MT_STRING:
		nval = mv_scan_number_nullable(pval->u.strv);
		mv_free(pval);
		*pval = nval;
		break;
	default:
		fprintf(stderr, "%s: internal coding error detected at file %s, line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		break;
	}
}

mv_t mv_scan_number_nullable(char* string) {
	double fltv = 0.0;
	long long intv = 0LL;
	mv_t rv = mv_void();
	if (*string == '\0') {
		// keep rv = mv_void();
	} else if (mlr_try_int_from_string(string, &intv)) {
		rv = mv_from_int(intv);
	} else if (mlr_try_float_from_string(string, &fltv)) {
		rv = mv_from_float(fltv);
	} else {
		rv = mv_error();
	}
	return rv;
}

mv_t mv_scan_number_or_die(char* string) {
	mv_t rv = mv_scan_number_nullable(string);
	if (!mv_is_numeric(&rv)) {
		fprintf(stderr, "%s: couldn't parse \"%s\" as number.\n",
			MLR_GLOBALS.argv0, string);
	}
	return rv;
}

// ----------------------------------------------------------------
mv_t s_ss_dot_func(mv_t* pval1, mv_t* pval2) {
	int len1 = strlen(pval1->u.strv);
	int len2 = strlen(pval2->u.strv);
	int len3 = len1 + len2 + 1; // for the null-terminator byte
	char* string3 = mlr_malloc_or_die(len3);
	strcpy(&string3[0], pval1->u.strv);
	strcpy(&string3[len1], pval2->u.strv);

	mv_free(pval1);
	mv_free(pval2);

	return mv_from_string_with_free(string3);
}

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

	if (matched) {
		mv_free(pval1);
		mv_free(pval3);
		return mv_from_string_with_free(output);
	} else {
		mv_free(pval3);
		return mv_from_string_no_free(output);
	}
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
	unsigned char free_flags = NO_FREE;
	char* output     = regex_gsub(input, pregex, psb, pval3->u.strv, &matched, &all_captured, &free_flags);

	if (matched) {
		mv_free(pval1);
		mv_free(pval3);
		return mv_from_string(output, free_flags);
	} else {
		mv_free(pval3);
		return mv_from_string(output, free_flags);
	}
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
#define NZBUFLEN 63

// Precondition: psec is either int or float.
mv_t time_string_from_seconds(mv_t* psec, char* format) {
	time_t clock = 0;
	if (psec->type == MT_FLOAT) {
		if (isinf(psec->u.fltv) || isnan(psec->u.fltv)) {
			return mv_error();
		}
		clock = (time_t) psec->u.fltv;
	} else {
		clock = (time_t) psec->u.intv;
	}

	struct tm tm;
	struct tm *ptm = gmtime_r(&clock, &tm);
	if (ptm == NULL) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}
	char* string = mlr_malloc_or_die(NZBUFLEN + 1);
	int written_length = strftime(string, NZBUFLEN, format, ptm);
	if (written_length > NZBUFLEN || written_length == 0) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}

	return mv_from_string_with_free(string);
}

mv_t s_n_sec2gmt_func(mv_t* pval1) {
	return time_string_from_seconds(pval1, ISO8601_TIME_FORMAT);
}

mv_t s_ns_strftime_func(mv_t* pval1, mv_t* pval2) {
	mv_t rv = time_string_from_seconds(pval1, pval2->u.strv);
	mv_free(pval2);
	return rv;
}

// ----------------------------------------------------------------
static mv_t seconds_from_time_string(char* time, char* format) {
	if (*time == '\0') {
		return mv_void();
	} else {
		struct tm tm;
		memset(&tm, 0, sizeof(tm));
		char* retval = strptime(time, format, &tm);
		if (retval == NULL) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}
		if (*retval != 0) { // Parseable input followed by non-parseable
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d: unparseable trailer \"%s\".\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__, retval);
			exit(1);
		}
		time_t t = mlr_timegm(&tm);
		return mv_from_int((long long)t);
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
static mv_t plus_v_xx(mv_t* pa, mv_t* pb) {
	return mv_void();
}
static mv_t plus_e_xx(mv_t* pa, mv_t* pb) {
	return mv_error();
}
static mv_t plus_f_uf(mv_t* pa, mv_t* pb) {
	return *pb;
}
static mv_t plus_f_fu(mv_t* pa, mv_t* pb) {
	return *pa;
}
static mv_t plus_i_ui(mv_t* pa, mv_t* pb) {
	return *pb;
}
static mv_t plus_i_iu(mv_t* pa, mv_t* pb) {
	return *pa;
}
static mv_t plus_i_uu(mv_t* pa, mv_t* pb) {
	return mv_from_int(0LL);
}

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
	//         ERROR       ABSENT     UNINIT     VOID       STRING     INT        FLOAT      BOOL
	/*ERROR*/  {plus_e_xx, plus_e_xx, plus_e_xx, plus_e_xx, plus_e_xx, plus_e_xx, plus_e_xx, plus_e_xx},
	/*ABSENT*/ {plus_e_xx, plus_v_xx, plus_v_xx, plus_v_xx, plus_e_xx, plus_i_ui, plus_f_uf, plus_e_xx},
	/*UNINIT*/ {plus_e_xx, plus_v_xx, plus_i_uu, plus_v_xx, plus_e_xx, plus_i_ui, plus_f_uf, plus_e_xx},
	/*VOID*/   {plus_e_xx, plus_v_xx, plus_v_xx, plus_v_xx, plus_e_xx, plus_v_xx, plus_v_xx, plus_e_xx},
	/*STRING*/ {plus_e_xx, plus_e_xx, plus_e_xx, plus_e_xx, plus_e_xx, plus_e_xx, plus_e_xx, plus_e_xx},
	/*INT*/    {plus_e_xx, plus_i_iu, plus_i_iu, plus_v_xx, plus_e_xx, plus_n_ii, plus_f_if, plus_e_xx},
	/*FLOAT*/  {plus_e_xx, plus_f_fu, plus_f_fu, plus_v_xx, plus_e_xx, plus_f_fi, plus_f_ff, plus_e_xx},
	/*BOOL*/   {plus_e_xx, plus_e_xx, plus_e_xx, plus_e_xx, plus_e_xx, plus_e_xx, plus_e_xx, plus_e_xx},
};

mv_t x_xx_plus_func(mv_t* pval1, mv_t* pval2) { return (plus_dispositions[pval1->type][pval2->type])(pval1,pval2); }

// ----------------------------------------------------------------
static mv_t minus_v_xx(mv_t* pa, mv_t* pb) {
	return mv_void();
}
static mv_t minus_e_xx(mv_t* pa, mv_t* pb) {
	return mv_error();
}
static mv_t minus_f_uf(mv_t* pa, mv_t* pb) {
	return *pb;
}
static mv_t minus_f_fu(mv_t* pa, mv_t* pb) {
	return *pa;
}
static mv_t minus_i_ui(mv_t* pa, mv_t* pb) {
	return *pb;
}
static mv_t minus_i_iu(mv_t* pa, mv_t* pb) {
	return *pa;
}
static mv_t minus_i_uu(mv_t* pa, mv_t* pb) {
	return mv_from_int(0LL);
}

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
	//         ERROR        ABSENT      UNINIT      VOID        STRING      INT         FLOAT       BOOL
	/*ERROR*/  {minus_e_xx, minus_e_xx, minus_e_xx, minus_e_xx, minus_e_xx, minus_e_xx, minus_e_xx, minus_e_xx},
	/*ABSENT*/ {minus_e_xx, minus_v_xx, minus_v_xx, minus_v_xx, minus_e_xx, minus_i_ui, minus_f_uf, minus_e_xx},
	/*UNINIT*/ {minus_e_xx, minus_v_xx, minus_i_uu, minus_v_xx, minus_e_xx, minus_i_ui, minus_f_uf, minus_e_xx},
	/*VOID*/   {minus_e_xx, minus_v_xx, minus_v_xx, minus_v_xx, minus_e_xx, minus_v_xx, minus_v_xx, minus_e_xx},
	/*STRING*/ {minus_e_xx, minus_e_xx, minus_e_xx, minus_e_xx, minus_e_xx, minus_e_xx, minus_e_xx, minus_e_xx},
	/*INT*/    {minus_e_xx, minus_i_iu, minus_i_iu, minus_v_xx, minus_e_xx, minus_n_ii, minus_f_if, minus_e_xx},
	/*FLOAT*/  {minus_e_xx, minus_f_fu, minus_f_fu, minus_v_xx, minus_e_xx, minus_f_fi, minus_f_ff, minus_e_xx},
	/*BOOL*/   {minus_e_xx, minus_e_xx, minus_e_xx, minus_e_xx, minus_e_xx, minus_e_xx, minus_e_xx, minus_e_xx},
};

mv_t x_xx_minus_func(mv_t* pval1, mv_t* pval2) { return (minus_dispositions[pval1->type][pval2->type])(pval1,pval2); }

// ----------------------------------------------------------------
static mv_t times_v_xx(mv_t* pa, mv_t* pb) {
	return mv_void();
}
static mv_t times_e_xx(mv_t* pa, mv_t* pb) {
	return mv_error();
}
static mv_t times_f_uf(mv_t* pa, mv_t* pb) {
	return *pb;
}
static mv_t times_f_fu(mv_t* pa, mv_t* pb) {
	return *pa;
}
static mv_t times_i_ui(mv_t* pa, mv_t* pb) {
	return *pb;
}
static mv_t times_i_iu(mv_t* pa, mv_t* pb) {
	return *pa;
}
static mv_t times_i_uu(mv_t* pa, mv_t* pb) {
	return mv_from_int(1LL);
}

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
	//         ERROR        ABSENT      UNINIT      VOID        STRING      INT         FLOAT       BOOL
	/*ERROR*/  {times_e_xx, times_e_xx, times_e_xx, times_e_xx, times_e_xx, times_e_xx, times_e_xx, times_e_xx},
	/*ABSENT*/ {times_e_xx, times_v_xx, times_v_xx, times_v_xx, times_e_xx, times_i_ui, times_f_uf, times_e_xx},
	/*UNINIT*/ {times_e_xx, times_v_xx, times_i_uu, times_v_xx, times_e_xx, times_i_ui, times_f_uf, times_e_xx},
	/*VOID*/   {times_e_xx, times_v_xx, times_v_xx, times_v_xx, times_e_xx, times_v_xx, times_v_xx, times_e_xx},
	/*STRING*/ {times_e_xx, times_e_xx, times_e_xx, times_e_xx, times_e_xx, times_e_xx, times_e_xx, times_e_xx},
	/*INT*/    {times_e_xx, times_i_iu, times_i_iu, times_v_xx, times_e_xx, times_n_ii, times_f_if, times_e_xx},
	/*FLOAT*/  {times_e_xx, times_f_fu, times_f_fu, times_v_xx, times_e_xx, times_f_fi, times_f_ff, times_e_xx},
	/*BOOL*/   {times_e_xx, times_e_xx, times_e_xx, times_e_xx, times_e_xx, times_e_xx, times_e_xx, times_e_xx},
};

mv_t x_xx_times_func(mv_t* pval1, mv_t* pval2) { return (times_dispositions[pval1->type][pval2->type])(pval1,pval2); }

// ----------------------------------------------------------------
static mv_t divide_v_xx(mv_t* pa, mv_t* pb) {
	return mv_void();
}
static mv_t divide_e_xx(mv_t* pa, mv_t* pb) {
	return mv_error();
}
static mv_t divide_f_uf(mv_t* pa, mv_t* pb) {
	return mv_from_float(0.0);
}
static mv_t divide_f_fu(mv_t* pa, mv_t* pb) {
	return mv_void();
}
static mv_t divide_i_ui(mv_t* pa, mv_t* pb) {
	return mv_from_int(0LL);
}
static mv_t divide_i_iu(mv_t* pa, mv_t* pb) {
	return *pa;
}
static mv_t divide_i_uu(mv_t* pa, mv_t* pb) {
	return mv_void();
}

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
	//         ERROR         ABSENT       UNINIT       VOID         STRING       INT          FLOAT        BOOL
	/*ERROR*/  {divide_e_xx, divide_e_xx, divide_e_xx, divide_e_xx, divide_e_xx, divide_e_xx, divide_e_xx, divide_e_xx},
	/*ABSENT*/ {divide_e_xx, divide_v_xx, divide_v_xx, divide_v_xx, divide_e_xx, divide_i_ui, divide_f_uf, divide_e_xx},
	/*UNINIT*/ {divide_e_xx, divide_v_xx, divide_i_uu, divide_v_xx, divide_e_xx, divide_i_ui, divide_f_uf, divide_e_xx},
	/*VOID*/   {divide_e_xx, divide_v_xx, divide_v_xx, divide_v_xx, divide_e_xx, divide_e_xx, divide_e_xx, divide_e_xx},
	/*STRING*/ {divide_e_xx, divide_e_xx, divide_e_xx, divide_e_xx, divide_e_xx, divide_e_xx, divide_e_xx, divide_e_xx},
	/*INT*/    {divide_e_xx, divide_i_iu, divide_i_iu, divide_v_xx, divide_e_xx, divide_i_ii, divide_f_if, divide_e_xx},
	/*FLOAT*/  {divide_e_xx, divide_f_fu, divide_f_fu, divide_v_xx, divide_e_xx, divide_f_fi, divide_f_ff, divide_e_xx},
	/*BOOL*/   {divide_e_xx, divide_e_xx, divide_e_xx, divide_e_xx, divide_e_xx, divide_e_xx, divide_e_xx, divide_e_xx},
};

mv_t x_xx_divide_func(mv_t* pval1, mv_t* pval2) { return (divide_dispositions[pval1->type][pval2->type])(pval1,pval2); }

// ----------------------------------------------------------------
static mv_t idiv_v_xx(mv_t* pa, mv_t* pb) {
	return mv_void();
}
static mv_t idiv_e_xx(mv_t* pa, mv_t* pb) {
	return mv_error();
}
static mv_t idiv_f_uf(mv_t* pa, mv_t* pb) {
	return mv_from_float(0.0);
}
static mv_t idiv_f_fu(mv_t* pa, mv_t* pb) {
	return mv_void();
}
static mv_t idiv_i_ui(mv_t* pa, mv_t* pb) {
	return mv_from_int(0LL);
}
static mv_t idiv_i_iu(mv_t* pa, mv_t* pb) {
	return *pa;
}
static mv_t idiv_i_uu(mv_t* pa, mv_t* pb) {
	return mv_void();
}

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
	//         ERROR       ABSENT     UNINIT     VOID       STRING     INT        FLOAT      BOOL
	/*ERROR*/  {idiv_e_xx, idiv_e_xx, idiv_e_xx, idiv_e_xx, idiv_e_xx, idiv_e_xx, idiv_e_xx, idiv_e_xx},
	/*ABSENT*/ {idiv_e_xx, idiv_v_xx, idiv_v_xx, idiv_v_xx, idiv_e_xx, idiv_i_ui, idiv_f_uf, idiv_e_xx},
	/*UNINIT*/ {idiv_e_xx, idiv_v_xx, idiv_i_uu, idiv_v_xx, idiv_e_xx, idiv_i_ui, idiv_f_uf, idiv_e_xx},
	/*VOID*/   {idiv_e_xx, idiv_v_xx, idiv_v_xx, idiv_v_xx, idiv_e_xx, idiv_e_xx, idiv_e_xx, idiv_e_xx},
	/*STRING*/ {idiv_e_xx, idiv_e_xx, idiv_e_xx, idiv_e_xx, idiv_e_xx, idiv_e_xx, idiv_e_xx, idiv_e_xx},
	/*INT*/    {idiv_e_xx, idiv_i_iu, idiv_i_iu, idiv_v_xx, idiv_e_xx, idiv_i_ii, idiv_f_if, idiv_e_xx},
	/*FLOAT*/  {idiv_e_xx, idiv_f_fu, idiv_f_fu, idiv_v_xx, idiv_e_xx, idiv_f_fi, idiv_f_ff, idiv_e_xx},
	/*BOOL*/   {idiv_e_xx, idiv_e_xx, idiv_e_xx, idiv_e_xx, idiv_e_xx, idiv_e_xx, idiv_e_xx, idiv_e_xx},
};


mv_t x_xx_int_divide_func(mv_t* pval1, mv_t* pval2) {
	return (idiv_dispositions[pval1->type][pval2->type])(pval1,pval2);
}

// ----------------------------------------------------------------
static mv_t mod_v_xx(mv_t* pa, mv_t* pb) {
	return mv_void();
}
static mv_t mod_e_xx(mv_t* pa, mv_t* pb) {
	return mv_error();
}
static mv_t mod_f_uf(mv_t* pa, mv_t* pb) {
	return mv_from_float(0.0);
}
static mv_t mod_f_fu(mv_t* pa, mv_t* pb) {
	return mv_void();
}
static mv_t mod_i_ui(mv_t* pa, mv_t* pb) {
	return mv_from_int(0LL);
}
static mv_t mod_i_iu(mv_t* pa, mv_t* pb) {
	return *pa;
}
static mv_t mod_i_uu(mv_t* pa, mv_t* pb) {
	return mv_void();
}

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
	//         ERROR      ABSENT    UNINIT    VOID      STRING    INT       FLOAT     BOOL
	/*ERROR*/  {mod_e_xx, mod_e_xx, mod_e_xx, mod_e_xx, mod_e_xx, mod_e_xx, mod_e_xx, mod_e_xx},
	/*ABSENT*/ {mod_e_xx, mod_v_xx, mod_v_xx, mod_v_xx, mod_e_xx, mod_i_ui, mod_f_uf, mod_e_xx},
	/*UNINIT*/ {mod_e_xx, mod_v_xx, mod_i_uu, mod_v_xx, mod_e_xx, mod_i_ui, mod_f_uf, mod_e_xx},
	/*VOID*/   {mod_e_xx, mod_v_xx, mod_v_xx, mod_v_xx, mod_e_xx, mod_e_xx, mod_e_xx, mod_e_xx},
	/*STRING*/ {mod_e_xx, mod_e_xx, mod_e_xx, mod_e_xx, mod_e_xx, mod_e_xx, mod_e_xx, mod_e_xx},
	/*INT*/    {mod_e_xx, mod_i_iu, mod_i_iu, mod_v_xx, mod_e_xx, mod_i_ii, mod_f_if, mod_e_xx},
	/*FLOAT*/  {mod_e_xx, mod_f_fu, mod_f_fu, mod_v_xx, mod_e_xx, mod_f_fi, mod_f_ff, mod_e_xx},
	/*BOOL*/   {mod_e_xx, mod_e_xx, mod_e_xx, mod_e_xx, mod_e_xx, mod_e_xx, mod_e_xx, mod_e_xx},
};

mv_t x_xx_mod_func(mv_t* pval1, mv_t* pval2) {
	return (mod_dispositions[pval1->type][pval2->type])(pval1,pval2);
}

// ----------------------------------------------------------------
static mv_t upos_e_x(mv_t* pa) {
	return mv_error();
}
static mv_t upos_v_x(mv_t* pa) {
	return mv_void();
}
static mv_t upos_i_x(mv_t* pa) {
	return mv_from_int(0LL);
}
static mv_t upos_i_i(mv_t* pa) {
	return mv_from_int(pa->u.intv);
}
static mv_t upos_f_f(mv_t* pa) {
	return mv_from_float(pa->u.fltv);
}

static mv_unary_func_t* upos_dispositions[MT_DIM] = {
	/*ERROR*/  upos_e_x,
	/*ABSENT*/ upos_v_x,
	/*UNINIT*/ upos_i_x,
	/*VOID*/   upos_v_x,
	/*STRING*/ upos_e_x,
	/*INT*/    upos_i_i,
	/*FLOAT*/  upos_f_f,
	/*BOOL*/   upos_e_x,
};

mv_t x_x_upos_func(mv_t* pval1) { return (upos_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t uneg_e_x(mv_t* pa) {
	return mv_error();
}
static mv_t uneg_v_x(mv_t* pa) {
	return mv_void();
}
static mv_t uneg_i_x(mv_t* pa) {
	return mv_from_int(0LL);
}
static mv_t uneg_i_i(mv_t* pa) {
	return mv_from_int(-pa->u.intv);
}
static mv_t uneg_f_f(mv_t* pa) {
	return mv_from_float(-pa->u.fltv);
}

static mv_unary_func_t* uneg_disnegitions[MT_DIM] = {
	/*ERROR*/  uneg_e_x,
	/*ABSENT*/ uneg_v_x,
	/*UNINIT*/ uneg_i_x,
	/*VOID*/   uneg_v_x,
	/*STRING*/ uneg_e_x,
	/*INT*/    uneg_i_i,
	/*FLOAT*/  uneg_f_f,
	/*BOOL*/   uneg_e_x,
};

mv_t x_x_uneg_func(mv_t* pval1) { return (uneg_disnegitions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t abs_e_x(mv_t* pa) {
	return mv_error();
}
static mv_t abs_v_x(mv_t* pa) {
	return mv_void();
}
static mv_t abs_n_f(mv_t* pa) {
	return mv_from_float(fabs(pa->u.fltv));
}
static mv_t abs_n_i(mv_t* pa) {
	return mv_from_int(pa->u.intv < 0LL ? -pa->u.intv : pa->u.intv);
}

static mv_unary_func_t* abs_dispositions[MT_DIM] = {
	/*ERROR*/  abs_e_x,
	/*ABSENT*/ abs_v_x,
	/*UNINIT*/ abs_v_x,
	/*VOID*/   abs_v_x,
	/*STRING*/ abs_e_x,
	/*INT*/    abs_n_i,
	/*FLOAT*/  abs_n_f,
	/*BOOL*/   abs_e_x,
};

mv_t x_x_abs_func(mv_t* pval1) { return (abs_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t ceil_e_x(mv_t* pa) {
	return mv_error();
}
static mv_t ceil_v_x(mv_t* pa) {
	return mv_void();
}
static mv_t ceil_n_f(mv_t* pa) {
	return mv_from_float(ceil(pa->u.fltv));
}
static mv_t ceil_n_i(mv_t* pa) {
	return mv_from_int(pa->u.intv);
}

static mv_unary_func_t* ceil_dispositions[MT_DIM] = {
	/*ERROR*/  ceil_e_x,
	/*ABSENT*/ ceil_v_x,
	/*UNINIT*/ ceil_v_x,
	/*VOID*/   ceil_v_x,
	/*STRING*/ ceil_e_x,
	/*INT*/    ceil_n_i,
	/*FLOAT*/  ceil_n_f,
	/*BOOL*/   ceil_e_x,
};

mv_t x_x_ceil_func(mv_t* pval1) { return (ceil_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t floor_e_x(mv_t* pa) {
	return mv_error();
}
static mv_t floor_v_x(mv_t* pa) {
	return mv_void();
}
static mv_t floor_n_f(mv_t* pa) {
	return mv_from_float(floor(pa->u.fltv));
}
static mv_t floor_n_i(mv_t* pa) {
	return mv_from_int(pa->u.intv);
}

static mv_unary_func_t* floor_dispositions[MT_DIM] = {
	/*ERROR*/  floor_e_x,
	/*ABSENT*/ floor_v_x,
	/*UNINIT*/ floor_v_x,
	/*VOID*/   floor_v_x,
	/*STRING*/ floor_e_x,
	/*INT*/    floor_n_i,
	/*FLOAT*/  floor_n_f,
	/*BOOL*/   floor_e_x,
};

mv_t x_x_floor_func(mv_t* pval1) { return (floor_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t round_e_x(mv_t* pa) {
	return mv_error();
}
static mv_t round_v_x(mv_t* pa) {
	return mv_void();
}
static mv_t round_n_f(mv_t* pa) {
	return mv_from_float(round(pa->u.fltv));
}
static mv_t round_n_i(mv_t* pa) {
	return mv_from_int(pa->u.intv);
}

static mv_unary_func_t* round_dispositions[MT_DIM] = {
	/*ERROR*/  round_e_x,
	/*ABSENT*/ round_v_x,
	/*UNINIT*/ round_v_x,
	/*VOID*/   round_v_x,
	/*STRING*/ round_e_x,
	/*INT*/    round_n_i,
	/*FLOAT*/  round_n_f,
	/*BOOL*/   round_e_x,
};

mv_t x_x_round_func(mv_t* pval1) { return (round_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t roundm_e_xx(mv_t* pa, mv_t* pb) {
	return mv_error();
}
static mv_t roundm_v_xx(mv_t* pa, mv_t* pb) {
	return mv_void();
}
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
	//         ERROR         ABSENT       UNINIT       VOID         STRING       INT          FLOAT        BOOL
	/*ERROR*/  {roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx},
	/*ABSENT*/ {roundm_e_xx, roundm_v_xx, roundm_v_xx, roundm_v_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx},
	/*UNINIT*/ {roundm_e_xx, roundm_v_xx, roundm_v_xx, roundm_v_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx},
	/*VOID*/   {roundm_e_xx, roundm_v_xx, roundm_v_xx, roundm_v_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx},
	/*STRING*/ {roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx},
	/*INT*/    {roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_i_ii, roundm_f_if, roundm_e_xx},
	/*FLOAT*/  {roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_f_fi, roundm_f_ff, roundm_e_xx},
	/*BOOL*/   {roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx, roundm_e_xx},
};

mv_t x_xx_roundm_func(mv_t* pval1, mv_t* pval2) { return (roundm_dispositions[pval1->type][pval2->type])(pval1,pval2); }

// ----------------------------------------------------------------
static mv_t min_e_xx(mv_t* pa, mv_t* pb) {
	return mv_error();
}

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

static mv_t min_f_fv(mv_t* pa, mv_t* pb) {
	return mv_from_float(pa->u.fltv);
}

static mv_t min_f_if(mv_t* pa, mv_t* pb) {
	double a = (double)pa->u.intv;
	double b = pb->u.fltv;
	return mv_from_float(fmin(a, b));
}

static mv_t min_f_vf(mv_t* pa, mv_t* pb) {
	return mv_from_float(pb->u.fltv);
}

static mv_t min_i_ii(mv_t* pa, mv_t* pb) {
	long long a = pa->u.intv;
	long long b = pb->u.intv;
	return mv_from_int(a < b ? a : b);
}

static mv_t min_i_iv(mv_t* pa, mv_t* pb) {
	return mv_from_int(pa->u.intv);
}

static mv_t min_i_vi(mv_t* pa, mv_t* pb) {
	return mv_from_int(pb->u.intv);
}

static mv_t min_v_xx(mv_t* pa, mv_t* pb) {
	return mv_void();
}

static mv_binary_func_t* min_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR      ABSENT    UNINIT    VOID      STRING    INT       FLOAT     BOOL
	/*ERROR*/  {min_e_xx, min_e_xx, min_e_xx, min_e_xx, min_e_xx, min_e_xx, min_e_xx, min_e_xx},
	/*ABSENT*/ {min_e_xx, min_v_xx, min_v_xx, min_v_xx, min_e_xx, min_i_vi, min_f_vf, min_e_xx},
	/*UNINIT*/ {min_e_xx, min_v_xx, min_v_xx, min_v_xx, min_e_xx, min_i_vi, min_f_vf, min_e_xx},
	/*VOID*/   {min_e_xx, min_v_xx, min_v_xx, min_v_xx, min_e_xx, min_i_vi, min_f_vf, min_e_xx},
	/*STRING*/ {min_e_xx, min_e_xx, min_e_xx, min_e_xx, min_e_xx, min_e_xx, min_e_xx, min_e_xx},
	/*INT*/    {min_e_xx, min_i_iv, min_i_iv, min_i_iv, min_e_xx, min_i_ii, min_f_if, min_e_xx},
	/*FLOAT*/  {min_e_xx, min_f_fv, min_f_fv, min_f_fv, min_e_xx, min_f_fi, min_f_ff, min_e_xx},
	/*BOOL*/   {min_e_xx, min_e_xx, min_e_xx, min_e_xx, min_e_xx, min_e_xx, min_e_xx, min_e_xx},
};

mv_t x_xx_min_func(mv_t* pval1, mv_t* pval2) { return (min_dispositions[pval1->type][pval2->type])(pval1,pval2); }

// ----------------------------------------------------------------
static mv_t max_e_xx(mv_t* pa, mv_t* pb) {
	return mv_error();
}

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

static mv_t max_f_fv(mv_t* pa, mv_t* pb) {
	return mv_from_float(pa->u.fltv);
}

static mv_t max_f_if(mv_t* pa, mv_t* pb) {
	double a = (double)pa->u.intv;
	double b = pb->u.fltv;
	return mv_from_float(fmax(a, b));
}

static mv_t max_f_vf(mv_t* pa, mv_t* pb) {
	return mv_from_float(pb->u.fltv);
}

static mv_t max_i_ii(mv_t* pa, mv_t* pb) {
	long long a = pa->u.intv;
	long long b = pb->u.intv;
	return mv_from_int(a > b ? a : b);
}

static mv_t max_i_iv(mv_t* pa, mv_t* pb) {
	return mv_from_int(pa->u.intv);
}

static mv_t max_i_vi(mv_t* pa, mv_t* pb) {
	return mv_from_int(pb->u.intv);
}

static mv_t max_v_xx(mv_t* pa, mv_t* pb) {
	return mv_void();
}

static mv_binary_func_t* max_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR      ABSENT    UNINIT    VOID      STRING    INT       FLOAT     BOOL
	/*ERROR*/  {max_e_xx, max_e_xx, max_e_xx, max_e_xx, max_e_xx, max_e_xx, max_e_xx, max_e_xx},
	/*ABSENT*/ {max_e_xx, max_v_xx, max_v_xx, max_v_xx, max_e_xx, max_i_vi, max_f_vf, max_e_xx},
	/*UNINIT*/ {max_e_xx, max_v_xx, max_v_xx, max_v_xx, max_e_xx, max_i_vi, max_f_vf, max_e_xx},
	/*VOID*/   {max_e_xx, max_v_xx, max_v_xx, max_v_xx, max_e_xx, max_i_vi, max_f_vf, max_e_xx},
	/*STRING*/ {max_e_xx, max_e_xx, max_e_xx, max_e_xx, max_e_xx, max_e_xx, max_e_xx, max_e_xx},
	/*INT*/    {max_e_xx, max_i_iv, max_i_iv, max_i_iv, max_e_xx, max_i_ii, max_f_if, max_e_xx},
	/*FLOAT*/  {max_e_xx, max_f_fv, max_f_fv, max_f_fv, max_e_xx, max_f_fi, max_f_ff, max_e_xx},
	/*BOOL*/   {max_e_xx, max_e_xx, max_e_xx, max_e_xx, max_e_xx, max_e_xx, max_e_xx, max_e_xx},
};

mv_t x_xx_max_func(mv_t* pval1, mv_t* pval2) { return (max_dispositions[pval1->type][pval2->type])(pval1,pval2); }

// ----------------------------------------------------------------
static mv_t sgn_e_x(mv_t* pa) {
	return mv_error();
}
static mv_t sgn_v_x(mv_t* pa) {
	return mv_void();
}
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
	/*ERROR*/  sgn_e_x,
	/*ABSENT*/ sgn_v_x,
	/*UNINIT*/ sgn_v_x,
	/*VOID*/   sgn_v_x,
	/*STRING*/ sgn_e_x,
	/*INT*/    sgn_n_i,
	/*FLOAT*/  sgn_n_f,
	/*BOOL*/   sgn_e_x,
};

mv_t x_x_sgn_func(mv_t* pval1) { return (sgn_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t int_v_x(mv_t* pa) { return mv_void(); }
static mv_t int_e_x(mv_t* pa) { return mv_error(); }
static mv_t int_i_b(mv_t* pa) { return mv_from_int(pa->u.boolv ? 1 : 0); }
static mv_t int_i_f(mv_t* pa) { return mv_from_int((long long)round(pa->u.fltv)); }
static mv_t int_i_i(mv_t* pa) { return mv_from_int(pa->u.intv); }
static mv_t int_i_s(mv_t* pa) {
	if (*pa->u.strv == '\0') {
		mv_free(pa);
		return mv_void();
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
	/*ERROR*/  int_e_x,
	/*ABSENT*/ int_v_x,
	/*UNINIT*/ int_v_x,
	/*VOID*/   int_v_x,
	/*STRING*/ int_i_s,
	/*INT*/    int_i_i,
	/*FLOAT*/  int_i_f,
	/*BOOL*/   int_i_b,
};

mv_t i_x_int_func(mv_t* pval1) { return (int_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t float_v_x(mv_t* pa) { return mv_void(); }
static mv_t float_e_x(mv_t* pa) { return mv_error(); }
static mv_t float_f_b(mv_t* pa) { return mv_from_float(pa->u.boolv ? 1.0 : 0.0); }
static mv_t float_f_f(mv_t* pa) { return mv_from_float(pa->u.fltv); }
static mv_t float_f_i(mv_t* pa) { return mv_from_float((double)pa->u.intv); }
static mv_t float_f_s(mv_t* pa) {
	if (*pa->u.strv == '\0') {
		mv_free(pa);
		return mv_void();
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
	/*ERROR*/  float_e_x,
	/*ABSENT*/ float_v_x,
	/*UNINIT*/ float_v_x,
	/*VOID*/   float_v_x,
	/*STRING*/ float_f_s,
	/*INT*/    float_f_i,
	/*FLOAT*/  float_f_f,
	/*BOOL*/   float_f_b,
};

mv_t f_x_float_func(mv_t* pval1) { return (float_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
mv_t b_x_isnull_func(mv_t* pval1) {
	mv_t rv = mv_from_bool(mv_is_null(pval1));
	mv_free(pval1);
	return rv;
}
mv_t b_x_isnotnull_func(mv_t* pval1) {
	mv_t rv = mv_from_bool(mv_is_non_null(pval1));
	mv_free(pval1);
	return rv;
}

// ----------------------------------------------------------------
mv_t b_x_isabsent_func(mv_t* pval1) {
	mv_t rv = mv_from_bool(mv_is_absent(pval1));
	mv_free(pval1);
	return rv;
}
mv_t b_x_ispresent_func(mv_t* pval1) {
	mv_t rv = mv_from_bool(mv_is_present(pval1));
	mv_free(pval1);
	return rv;
}

// ----------------------------------------------------------------
static mv_t boolean_v_x(mv_t* pa) { return mv_void(); }
static mv_t boolean_e_x(mv_t* pa) { return mv_error(); }
static mv_t boolean_b_b(mv_t* pa) { return mv_from_bool(pa->u.boolv); }
static mv_t boolean_b_f(mv_t* pa) { return mv_from_bool((pa->u.fltv == 0.0) ? FALSE : TRUE); }
static mv_t boolean_b_i(mv_t* pa) { return mv_from_bool((pa->u.intv == 0LL) ? FALSE : TRUE); }
static mv_t boolean_b_s(mv_t* pa) { return mv_from_bool((streq(pa->u.strv, "true") || streq(pa->u.strv, "TRUE")) ? TRUE : FALSE);}

static mv_unary_func_t* boolean_dispositions[MT_DIM] = {
	/*ERROR*/  boolean_e_x,
	/*ABSENT*/ boolean_v_x,
	/*UNINIT*/ boolean_v_x,
	/*VOID*/   boolean_v_x,
	/*STRING*/ boolean_b_s,
	/*INT*/    boolean_b_i,
	/*FLOAT*/  boolean_b_f,
	/*BOOL*/   boolean_b_b,
};

mv_t b_x_boolean_func(mv_t* pval1) { return (boolean_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t string_v_x(mv_t* pa) { return mv_void(); }
static mv_t string_e_x(mv_t* pa) { return mv_error(); }
static mv_t string_s_b(mv_t* pa) { return mv_from_string_no_free(pa->u.boolv?"true":"false"); }
static mv_t string_s_f(mv_t* pa) { return mv_from_string_with_free(mlr_alloc_string_from_double(pa->u.fltv, MLR_GLOBALS.ofmt)); }
static mv_t string_s_i(mv_t* pa) { return mv_from_string_with_free(mlr_alloc_string_from_ll(pa->u.intv)); }
static mv_t string_s_s(mv_t* pa) {
	unsigned char free_flags = pa->free_flags;
	pa->free_flags = NO_FREE;
	return mv_from_string(pa->u.strv, free_flags);
}

static mv_unary_func_t* string_dispositions[MT_DIM] = {
	/*ERROR*/  string_e_x,
	/*ABSENT*/ string_v_x,
	/*UNINIT*/ string_v_x,
	/*VOID*/   string_v_x,
	/*STRING*/ string_s_s,
	/*INT*/    string_s_i,
	/*FLOAT*/  string_s_f,
	/*BOOL*/   string_s_b,
};

mv_t s_x_string_func(mv_t* pval1) { return (string_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t hexfmt_v_x(mv_t* pa) { return mv_void(); }
static mv_t hexfmt_e_x(mv_t* pa) { return mv_error(); }
static mv_t hexfmt_s_b(mv_t* pa) { return mv_from_string_no_free(pa->u.boolv?"0x1":"0x0"); }
static mv_t hexfmt_s_f(mv_t* pa) { return mv_from_string_with_free(mlr_alloc_hexfmt_from_ll((long long)pa->u.fltv)); }
static mv_t hexfmt_s_i(mv_t* pa) { return mv_from_string_with_free(mlr_alloc_hexfmt_from_ll(pa->u.intv)); }
static mv_t hexfmt_s_s(mv_t* pa) {
	unsigned char free_flags = pa->free_flags;
	pa->free_flags = NO_FREE;
	return mv_from_string(pa->u.strv, free_flags);
}

static mv_unary_func_t* hexfmt_dispositions[MT_DIM] = {
	/*ERROR*/  hexfmt_e_x,
	/*ABSENT*/ hexfmt_v_x,
	/*UNINIT*/ hexfmt_v_x,
	/*VOID*/   hexfmt_v_x,
	/*STRING*/ hexfmt_s_s,
	/*INT*/    hexfmt_s_i,
	/*FLOAT*/  hexfmt_s_f,
	/*BOOL*/   hexfmt_s_b,
};

mv_t s_x_hexfmt_func(mv_t* pval1) { return (hexfmt_dispositions[pval1->type])(pval1); }

// ----------------------------------------------------------------
static mv_t fmtnum_v_xx(mv_t* pa, mv_t* pfmt) { return mv_void(); }
static mv_t fmtnum_e_xx(mv_t* pa, mv_t* pfmt) { return mv_error(); }
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
static mv_t fmtnum_s_ss(mv_t* pa, mv_t* pfmt) { return mv_error(); }

static mv_binary_func_t* fmtnum_dispositions[MT_DIM] = {
	/*ERROR*/  fmtnum_e_xx,
	/*ABSENT*/ fmtnum_v_xx,
	/*UNINIT*/ fmtnum_v_xx,
	/*VOID*/   fmtnum_v_xx,
	/*STRING*/ fmtnum_s_ss,
	/*INT*/    fmtnum_s_is,
	/*FLOAT*/  fmtnum_s_ds,
	/*BOOL*/   fmtnum_s_bs,
};

mv_t s_xs_fmtnum_func(mv_t* pval1, mv_t* pval2) { return (fmtnum_dispositions[pval1->type])(pval1, pval2); }

// ----------------------------------------------------------------
static mv_t op_v_xx(mv_t* pa, mv_t* pb) { return mv_void(); }
static mv_t op_e_xx(mv_t* pa, mv_t* pb) { return mv_error(); }

static  mv_t eq_b_ii(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv == pb->u.intv); }
static  mv_t ne_b_ii(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv != pb->u.intv); }
static  mv_t gt_b_ii(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv >  pb->u.intv); }
static  mv_t ge_b_ii(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv >= pb->u.intv); }
static  mv_t lt_b_ii(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv <  pb->u.intv); }
static  mv_t le_b_ii(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv <= pb->u.intv); }

static  mv_t eq_b_ff(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv == pb->u.fltv); }
static  mv_t ne_b_ff(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv != pb->u.fltv); }
static  mv_t gt_b_ff(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv >  pb->u.fltv); }
static  mv_t ge_b_ff(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv >= pb->u.fltv); }
static  mv_t lt_b_ff(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv <  pb->u.fltv); }
static  mv_t le_b_ff(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv <= pb->u.fltv); }

static  mv_t eq_b_fi(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv == pb->u.intv); }
static  mv_t ne_b_fi(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv != pb->u.intv); }
static  mv_t gt_b_fi(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv >  pb->u.intv); }
static  mv_t ge_b_fi(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv >= pb->u.intv); }
static  mv_t lt_b_fi(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv <  pb->u.intv); }
static  mv_t le_b_fi(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv <= pb->u.intv); }

static  mv_t eq_b_if(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv == pb->u.fltv); }
static  mv_t ne_b_if(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv != pb->u.fltv); }
static  mv_t gt_b_if(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv >  pb->u.fltv); }
static  mv_t ge_b_if(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv >= pb->u.fltv); }
static  mv_t lt_b_if(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv <  pb->u.fltv); }
static  mv_t le_b_if(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv <= pb->u.fltv); }

static  mv_t eq_b_xs(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sa = mv_format_val(pa, &free_flags);
	mv_t rv = mv_from_bool(strcmp(sa, pb->u.strv) == 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sa);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static  mv_t ne_b_xs(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sa = mv_format_val(pa, &free_flags);
	mv_t rv = mv_from_bool(strcmp(sa, pb->u.strv) != 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sa);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static  mv_t gt_b_xs(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sa = mv_format_val(pa, &free_flags);
	mv_t rv = mv_from_bool(strcmp(sa, pb->u.strv) > 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sa);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static  mv_t ge_b_xs(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sa = mv_format_val(pa, &free_flags);
	mv_t rv = mv_from_bool(strcmp(sa, pb->u.strv) >= 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sa);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static  mv_t lt_b_xs(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sa = mv_format_val(pa, &free_flags);
	mv_t rv = mv_from_bool(strcmp(sa, pb->u.strv) < 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sa);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static  mv_t le_b_xs(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sa = mv_format_val(pa, &free_flags);
	mv_t rv = mv_from_bool(strcmp(sa, pb->u.strv) <= 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sa);
	mv_free(pa);
	mv_free(pb);
	return rv;
}

static  mv_t eq_b_sx(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sb = mv_format_val(pb, &free_flags);
	mv_t rv = mv_from_bool(strcmp(pa->u.strv, sb) == 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sb);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static  mv_t ne_b_sx(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sb = mv_format_val(pb, &free_flags);
	mv_t rv = mv_from_bool(strcmp(pa->u.strv, sb) != 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sb);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static  mv_t gt_b_sx(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sb = mv_format_val(pb, &free_flags);
	mv_t rv = mv_from_bool(strcmp(pa->u.strv, sb) > 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sb);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static  mv_t ge_b_sx(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sb = mv_format_val(pb, &free_flags);
	mv_t rv = mv_from_bool(strcmp(pa->u.strv, sb) >= 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sb);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static  mv_t lt_b_sx(mv_t* pa, mv_t* pb) {
	char free_flags;
	char* sb = mv_format_val(pb, &free_flags);
	mv_t rv = mv_from_bool(strcmp(pa->u.strv, sb) < 0);
	if (free_flags & FREE_ENTRY_VALUE)
		free(sb);
	mv_free(pa);
	mv_free(pb);
	return rv;
}
static  mv_t le_b_sx(mv_t* pa, mv_t* pb) {
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
	//         ERROR     ABSENT   UNINIT   VOID     STRING   INT      FLOAT    BOOL
	/*ERROR*/  {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
	/*ABSENT*/ {op_e_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_e_xx},
	/*UNINIT*/ {op_e_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_e_xx},
	/*VOID*/   {op_e_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_e_xx},
	/*STRING*/ {op_e_xx, op_v_xx, op_v_xx, op_v_xx, eq_b_ss, eq_b_sx, eq_b_sx, op_e_xx},
	/*INT*/    {op_e_xx, op_v_xx, op_v_xx, op_v_xx, eq_b_xs, eq_b_ii, eq_b_if, op_e_xx},
	/*FLOAT*/  {op_e_xx, op_v_xx, op_v_xx, op_v_xx, eq_b_xs, eq_b_fi, eq_b_ff, op_e_xx},
	/*BOOL*/   {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
};

static mv_binary_func_t* ne_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR     ABSENT   UNINIT   VOID     STRING   INT      FLOAT    BOOL
	/*ERROR*/  {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
	/*ABSENT*/ {op_e_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_e_xx},
	/*UNINIT*/ {op_e_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_e_xx},
	/*VOID*/   {op_e_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_e_xx},
	/*STRING*/ {op_e_xx, op_v_xx, op_v_xx, op_v_xx, ne_b_ss, ne_b_sx, ne_b_sx, op_e_xx},
	/*INT*/    {op_e_xx, op_v_xx, op_v_xx, op_v_xx, ne_b_xs, ne_b_ii, ne_b_if, op_e_xx},
	/*FLOAT*/  {op_e_xx, op_v_xx, op_v_xx, op_v_xx, ne_b_xs, ne_b_fi, ne_b_ff, op_e_xx},
	/*BOOL*/   {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
};

static mv_binary_func_t* gt_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR     ABSENT   UNINIT   VOID     STRING   INT      FLOAT    BOOL
	/*ERROR*/  {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
	/*ABSENT*/ {op_e_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_e_xx},
	/*UNINIT*/ {op_e_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_e_xx},
	/*VOID*/   {op_e_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_e_xx},
	/*STRING*/ {op_e_xx, op_e_xx, op_v_xx, op_e_xx, gt_b_ss, gt_b_sx, gt_b_sx, op_e_xx},
	/*INT*/    {op_e_xx, op_e_xx, op_v_xx, op_e_xx, gt_b_xs, gt_b_ii, gt_b_if, op_e_xx},
	/*FLOAT*/  {op_e_xx, op_e_xx, op_v_xx, op_e_xx, gt_b_xs, gt_b_fi, gt_b_ff, op_e_xx},
	/*BOOL*/   {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
};

static mv_binary_func_t* ge_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR     ABSENT   UNINIT   VOID     STRING   INT      FLOAT    BOOL
	/*ERROR*/  {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
	/*ABSENT*/ {op_e_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_e_xx},
	/*UNINIT*/ {op_e_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_e_xx},
	/*VOID*/   {op_e_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_e_xx},
	/*STRING*/ {op_e_xx, op_v_xx, op_v_xx, op_v_xx, ge_b_ss, ge_b_sx, ge_b_sx, op_e_xx},
	/*INT*/    {op_e_xx, op_v_xx, op_v_xx, op_v_xx, ge_b_xs, ge_b_ii, ge_b_if, op_e_xx},
	/*FLOAT*/  {op_e_xx, op_v_xx, op_v_xx, op_v_xx, ge_b_xs, ge_b_fi, ge_b_ff, op_e_xx},
	/*BOOL*/   {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
};

static mv_binary_func_t* lt_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR     ABSENT   UNINIT   VOID     STRING   INT      FLOAT    BOOL
	/*ERROR*/  {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
	/*ABSENT*/ {op_e_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_e_xx},
	/*UNINIT*/ {op_e_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_e_xx},
	/*VOID*/   {op_e_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_e_xx},
	/*STRING*/ {op_e_xx, op_v_xx, op_v_xx, op_v_xx, lt_b_ss, lt_b_sx, lt_b_sx, op_e_xx},
	/*INT*/    {op_e_xx, op_v_xx, op_v_xx, op_v_xx, lt_b_xs, lt_b_ii, lt_b_if, op_e_xx},
	/*FLOAT*/  {op_e_xx, op_v_xx, op_v_xx, op_v_xx, lt_b_xs, lt_b_fi, lt_b_ff, op_e_xx},
	/*BOOL*/   {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
};

static mv_binary_func_t* le_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR     ABSENT   UNINIT   VOID     STRING   INT      FLOAT    BOOL
	/*ERROR*/  {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
	/*ABSENT*/ {op_e_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_e_xx},
	/*UNINIT*/ {op_e_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_e_xx},
	/*VOID*/   {op_e_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_v_xx, op_e_xx},
	/*STRING*/ {op_e_xx, op_v_xx, op_v_xx, op_v_xx, le_b_ss, le_b_sx, le_b_sx, op_e_xx},
	/*INT*/    {op_e_xx, op_v_xx, op_v_xx, op_v_xx, le_b_xs, le_b_ii, le_b_if, op_e_xx},
	/*FLOAT*/  {op_e_xx, op_v_xx, op_v_xx, op_v_xx, le_b_xs, le_b_fi, le_b_ff, op_e_xx},
	/*BOOL*/   {op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx, op_e_xx},
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
	//         ERROR  ABSENT UNINIT VOID  STRING INT      FLOAT    BOOL
	/*ERROR*/  {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*ABSENT*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*UNINIT*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*VOID*/   {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*STRING*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*INT*/    {NULL, NULL,  NULL,  NULL, NULL,  eq_i_ii, eq_i_if, NULL},
	/*FLOAT*/  {NULL, NULL,  NULL,  NULL, NULL,  eq_i_fi, eq_i_ff, NULL},
	/*BOOL*/   {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
};

static mv_i_nn_comparator_func_t* ine_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT UNINIT VOID  STRING INT      FLOAT    BOOL
	/*ERROR*/  {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*ABSENT*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*UNINIT*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*VOID*/   {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*STRING*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*INT*/    {NULL, NULL,  NULL,  NULL, NULL,  ne_i_ii, ne_i_if, NULL},
	/*FLOAT*/  {NULL, NULL,  NULL,  NULL, NULL,  ne_i_fi, ne_i_ff, NULL},
	/*BOOL*/   {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
};

static mv_i_nn_comparator_func_t* igt_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT UNINIT VOID  STRING INT      FLOAT    BOOL
	/*ERROR*/  {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*ABSENT*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*UNINIT*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*VOID*/   {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*STRING*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*INT*/    {NULL, NULL,  NULL,  NULL, NULL,  gt_i_ii, gt_i_if, NULL},
	/*FLOAT*/  {NULL, NULL,  NULL,  NULL, NULL,  gt_i_fi, gt_i_ff, NULL},
	/*BOOL*/   {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
};

static mv_i_nn_comparator_func_t* ige_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT UNINIT VOID  STRING INT      FLOAT    BOOL
	/*ERROR*/  {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*ABSENT*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*UNINIT*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*VOID*/   {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*STRING*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*INT*/    {NULL, NULL,  NULL,  NULL, NULL,  ge_i_ii, ge_i_if, NULL},
	/*FLOAT*/  {NULL, NULL,  NULL,  NULL, NULL,  ge_i_fi, ge_i_ff, NULL},
	/*BOOL*/   {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
};

static mv_i_nn_comparator_func_t* ilt_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT UNINIT VOID  STRING INT      FLOAT    BOOL
	/*ERROR*/  {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*ABSENT*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*UNINIT*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*VOID*/   {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*STRING*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*INT*/    {NULL, NULL,  NULL,  NULL, NULL,  lt_i_ii, lt_i_if, NULL},
	/*FLOAT*/  {NULL, NULL,  NULL,  NULL, NULL,  lt_i_fi, lt_i_ff, NULL},
	/*BOOL*/   {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
};

static mv_i_nn_comparator_func_t* ile_dispositions[MT_DIM][MT_DIM] = {
	//         ERROR  ABSENT UNINIT VOID  STRING INT      FLOAT    BOOL
	/*ERROR*/  {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*ABSENT*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*UNINIT*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*VOID*/   {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*STRING*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
	/*INT*/    {NULL, NULL,  NULL,  NULL, NULL,  le_i_ii, le_i_if, NULL},
	/*FLOAT*/  {NULL, NULL,  NULL,  NULL, NULL,  le_i_fi, le_i_ff, NULL},
	/*BOOL*/   {NULL, NULL,  NULL,  NULL, NULL,  NULL,    NULL,    NULL},
};

int mv_i_nn_eq(mv_t* pval1, mv_t* pval2) { return (ieq_dispositions[pval1->type][pval2->type])(pval1, pval2); }
int mv_i_nn_ne(mv_t* pval1, mv_t* pval2) { return (ine_dispositions[pval1->type][pval2->type])(pval1, pval2); }
int mv_i_nn_gt(mv_t* pval1, mv_t* pval2) { return (igt_dispositions[pval1->type][pval2->type])(pval1, pval2); }
int mv_i_nn_ge(mv_t* pval1, mv_t* pval2) { return (ige_dispositions[pval1->type][pval2->type])(pval1, pval2); }
int mv_i_nn_lt(mv_t* pval1, mv_t* pval2) { return (ilt_dispositions[pval1->type][pval2->type])(pval1, pval2); }
int mv_i_nn_le(mv_t* pval1, mv_t* pval2) { return (ile_dispositions[pval1->type][pval2->type])(pval1, pval2); }

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
		if (ppregex_captures != NULL)
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
	//         ERROR  ABSENT UNINIT VOID  STRING INT               FLOAT             BOOL
	/*ERROR*/  {NULL, NULL,  NULL,  NULL, NULL,  NULL,             NULL,             NULL},
	/*ABSENT*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,             NULL,             NULL},
	/*UNINIT*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,             NULL,             NULL},
	/*VOID*/   {NULL, NULL,  NULL,  NULL, NULL,  NULL,             NULL,             NULL},
	/*STRING*/ {NULL, NULL,  NULL,  NULL, NULL,  NULL,             NULL,             NULL},
	/*INT*/    {NULL, NULL,  NULL,  NULL, NULL,  mv_ii_comparator, mv_if_comparator, NULL},
	/*FLOAT*/  {NULL, NULL,  NULL,  NULL, NULL,  mv_fi_comparator, mv_ff_comparator, NULL},
	/*BOOL*/   {NULL, NULL,  NULL,  NULL, NULL,  NULL,             NULL,             NULL},
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
