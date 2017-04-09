#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <sys/time.h>
#include <sys/stat.h>
#include "lib/mlr_globals.h"
#include "lib/mlr_arch.h"
#include "lib/mlrutil.h"
#include "lib/mlrdatetime.h"
#include "lib/string_builder.h"

// ----------------------------------------------------------------
// seconds since the epoch
double get_systime() {
	struct timeval tv = { .tv_sec = 0, .tv_usec = 0 };
	(void)gettimeofday(&tv, NULL);
	return (double)tv.tv_sec + (double)tv.tv_usec * 1e-6;
}

// ----------------------------------------------------------------
#define NZBUFLEN 63
// The essential idea is that we use the library function gmtime to get a struct tm, then strftime
// to produce a formatted string. The only complication is that we support "%1S" through "%9S" for
// formatting the seconds with a desired number of decimal places.

char* mlr_alloc_time_string_from_seconds(double seconds_since_the_epoch, char* format_string) {

	// 1. Split out the integer seconds since the epoch, which the stdlib can handle, and
	//    the fractional part, which it cannot.
	time_t iseconds = (time_t) seconds_since_the_epoch;
	double fracsec = seconds_since_the_epoch - iseconds;

	struct tm tm = *gmtime(&iseconds); // No gmtime_r on Windows so just use gmtime.

	// 2. See if "%nS" (for n in 1..9) is a substring of the format string.
	char* middle_nS_format = NULL;
	char* right_subformat = NULL;
	for (char* p = format_string; *p; p++) {
		// We can't use strstr since we're searching for a pattern, and regexes are overkill.
		// Here we rely on left-to-right evaluation of the boolean expressions, with non-evaluation
		// of a subexpression if a subexpression to its left is false (this keeps us from reading
		// past end of input string).
		if (p[0] == '%' && p[1] >= '1' && p[1] <= '9' && p[2] == 'S') {
			middle_nS_format = p;
			right_subformat = &p[3];
			break;
		}
	}

	// 3. If "%nS" (for n in 1..9) is not a substring of the format string, just use strftime.
	if (middle_nS_format == NULL) {
		char* output_string = mlr_malloc_or_die(NZBUFLEN+1);
		int written_length = strftime(output_string, NZBUFLEN, format_string, &tm);
		if (written_length > NZBUFLEN || written_length == 0) {
			fprintf(stderr, "%s: could not strftime(%lf, \"%s\"). See \"%s --help-function strftime\".\n",
				MLR_GLOBALS.bargv0, seconds_since_the_epoch, format_string, MLR_GLOBALS.bargv0);
			exit(1);
		}

		return output_string;
	}

	// Now we know "%nS" (for n in 1..9) is a substring of the format string.  Idea is to
	// copy the subformats to the left and right of the %nS part and format them both using
	// strftime, and format the middle part ourselves using sprintf.  Then concatenate all
	// those pieces.

	// 5. Find the left-of-%nS subformat, and format the input using that.
	int left_subformat_length = middle_nS_format - format_string;
	char* left_subformat = mlr_malloc_or_die(left_subformat_length + 1);
	memcpy(left_subformat, format_string, left_subformat_length);
	left_subformat[left_subformat_length] = 0;

	char left_formatted[NZBUFLEN+1];
	if (*left_subformat == 0) {
		// There's nothing to the left of %nS. strftime will error on empty format string, so we can
		// just map empty format to empty result ourselves.
		*left_formatted = 0;
	} else {
		int written_length = strftime(left_formatted, NZBUFLEN, left_subformat, &tm);
		if (written_length > NZBUFLEN || written_length == 0) {
			fprintf(stderr, "%s: could not strftime(%lf, \"%s\"). See \"%s --help-function strftime\".\n",
				MLR_GLOBALS.bargv0, seconds_since_the_epoch, format_string, MLR_GLOBALS.bargv0);
			exit(1);
		}
	}
	free(left_subformat);

	// 6. There are two parts in the middle: the integer part which strftime can populate
	//    from the struct tm, using %S format, and the fractional-seconds part which we sprintf.
	//    First do the int part.
	char middle_int_formatted[NZBUFLEN+1];
	char* middle_int_format = "%S";
	int written_length = strftime(middle_int_formatted, NZBUFLEN, middle_int_format, &tm);
	if (written_length > NZBUFLEN || written_length == 0) {
		fprintf(stderr, "%s: could not strftime(%lf, \"%s\"). See \"%s --help-function strftime\".\n",
			MLR_GLOBALS.bargv0, seconds_since_the_epoch, format_string, MLR_GLOBALS.bargv0);
		exit(1);
	}

	// 7. Do the fractional-seconds part. One key point is that sprintf always writes a leading zero,
	//    e.g. .123456 becomes "0.123456". We'll take off the leading zero later.
	char middle_sprintf_format[] = "%.xlf";
	char middle_fractional_formatted[16];
	// "%6S" maps to "%.6lf" and so on. We found the middle_nS_format by searching for "%nS" for
	// n in 1..9 so sprintf-format subscript 2 is the same as strftime format subscript 1.
	middle_sprintf_format[2] = middle_nS_format[1];
	sprintf(middle_fractional_formatted, middle_sprintf_format, fracsec);

	// 8. Format the right-of-%nS part, also using strftime.
	char right_formatted[NZBUFLEN];
	if (*right_subformat == 0) {
		// There's nothing to the right of %nS. strftime will error on empty format string, so we can
		// just map empty format to empty result ourselves.
		*right_formatted = 0;
	} else {
		int written_length = strftime(right_formatted, NZBUFLEN, right_subformat, &tm);
		if (written_length > NZBUFLEN || written_length == 0) {
			fprintf(stderr, "%s: could not strftime(%lf, \"%s\"). See \"%s --help-function strftime\".\n",
				MLR_GLOBALS.bargv0, seconds_since_the_epoch, format_string, MLR_GLOBALS.bargv0);
			exit(1);
		}
	}

	// 9. Concatenate the output. For string_builder, the size argument is just an initial size;
	//    it can realloc beyond that initial estimate if necessary.

	string_builder_t* psb = sb_alloc(NZBUFLEN+1);
	sb_append_string(psb, left_formatted);
	sb_append_string(psb, middle_int_formatted);
	MLR_INTERNAL_CODING_ERROR_IF(middle_fractional_formatted[0] != '0');
	sb_append_string(psb, &middle_fractional_formatted[1]); // xxx comment
	sb_append_string(psb, right_formatted);
	char* output_string = sb_finish(psb);
	sb_free(psb);

	return output_string;
}

// ----------------------------------------------------------------
double mlr_seconds_from_time_string(char* string, char* format) {
	struct tm tm;
	memset(&tm, 0, sizeof(tm));
	char* retval = mlr_arch_strptime(string, format, &tm);
	if (retval == NULL) {
		fprintf(stderr, "%s: could not strptime(\"%s\", \"%s\"). See \"%s --help-function strptime\".\n",
			MLR_GLOBALS.bargv0, string, format, MLR_GLOBALS.bargv0);
		exit(1);
	}
	MLR_INTERNAL_CODING_ERROR_IF(*retval != 0); // Parseable input followed by non-parseable
	time_t iseconds = mlr_arch_timegm(&tm);
	return (double)iseconds;
}

//#define MYBUFLEN 256

//// ----------------------------------------------------------------
//static int p(char* time_string, char* format_string) {
//
//	struct tm tm;
//	memset(&tm, 0, sizeof(tm));
//
//	// xxx cmt try the non-floating-point-seconds case first and return quickly if so.
//	char* strptime_retval = mlr_arch_strptime(time_string, format_string, &tm);
//	if (strptime_retval != NULL) {
//		if (*strptime_retval != 0) { // xxx extraneous
//			fprintf(stderr, "%s: could not strptime(\"%s\", \"%s\"). See \"%s --help-function strptime\".\n",
//				MLR_GLOBALS.bargv0, time_string, format_string, MLR_GLOBALS.bargv0);
//			exit(1);
//		}
//		time_t iseconds = mlr_arch_timegm(&tm);
//		printf("%s, %s -> %u, \"%s\"\n", time_string, format_string, (unsigned)iseconds, strptime_retval);
//		return 0;
//	}
//
//	char* pS = strstr(format_string, "%S");
//	if (pS == NULL) {
//		// Couldn't have been because of floating-point-seconds stuff. No reason to try any harder.
//		fprintf(stderr, "%s: could not strptime(\"%s\", \"%s\"). See \"%s --help-function strptime\".\n",
//			MLR_GLOBALS.bargv0, time_string, format_string, MLR_GLOBALS.bargv0);
//		exit(1);
//	}
//
//	// xxx
//	// Input is    "2017-04-09T00:51:09.123456 TZBLAHBLAH"
//	// with format "%Y-%m-%dT%H:%M:%S TZBLAHBLAH"
//
//	// 1. Copy the format up to the %S but with nothing else after. This is temporary to help us locate
//	//    the fractional-seconds part of the input string.
//	//    Example temporary format: "%Y-%m-%dT%H:%M:%S"
//
//	int truncated_format_string_length = pS - format_string + 2;
//	char* truncated_format_string = mlr_malloc_or_die(truncated_format_string_length + 1);
//	memcpy(truncated_format_string, format_string, truncated_format_string_length);
//	truncated_format_string[truncated_format_string_length] = 0;
//	//printf("ORIGINAL  FORMAT \"%s\"\n", format_string);
//	//printf("TRUNCATED FORMAT \"%s\"\n", truncated_format_string);
//
//	// 2. strptime using that truncated format and ignore the tm. Only look at the string return value.
//	//    Example return value: ".123456 TZBLAHBLAH"
//
//	strptime_retval = mlr_arch_strptime(time_string, truncated_format_string, &tm);
//	if (strptime_retval == NULL) {
//		fprintf(stderr, "%s: could not strptime(\"%s\", \"%s\"). See \"%s --help-function strptime\".\n",
//			MLR_GLOBALS.bargv0, time_string, format_string, MLR_GLOBALS.bargv0);
//		exit(1);
//	}
//	//printf("STRPTIME TEMP RETVAL \"%s\"\n", strptime_retval);
//
//	// 3. strtod the return value to find the fractional-seconds part, and whatever's after.
//	//    Example fractional-seconds part: ".123456"
//	//    Example stuff after: " TZBLAHBLAH"
//
//	char* stuff_after = NULL;
//	double fractional_seconds = strtod(strptime_retval, &stuff_after);
//	//printf("FRACTIONAL SECONDS %.6lf\n", fractional_seconds);
//	//printf("STUFF AFTER  \"%s\"\n", stuff_after);
//
//	// 4. Make a copy of the input string with the fractional seconds elided.
//	//    Example: "2017-04-09T00:51:09 TZBLAHBLAH"
//	char* elided_fraction_input = mlr_malloc_or_die(strlen(time_string) + 1);
//	int input_length_to_fractional_seconds = strptime_retval - time_string;
//	memcpy(elided_fraction_input, time_string, input_length_to_fractional_seconds);
//	strcpy(&elided_fraction_input[input_length_to_fractional_seconds], stuff_after);
//	//printf("ELIDE \"%s\"\n", elided_fraction_input);
//
//	// 5. strptime the elided-fraction input string using the original format string. Get the tm.
//	memset(&tm, 0, sizeof(tm));
//	strptime_retval = mlr_arch_strptime(elided_fraction_input, format_string, &tm);
//	if (strptime_retval == NULL) {
//		fprintf(stderr, "%s: could not strptime(\"%s\", \"%s\"). See \"%s --help-function strptime\".\n",
//			MLR_GLOBALS.bargv0, time_string, format_string, MLR_GLOBALS.bargv0);
//		exit(1);
//	}
//	//printf("STRPTIME ELIDE RETVAL \"%s\"\n", strptime_retval);
//
//	// 6. Convert the tm to a time_t (seconds since the epoch) and then add the fractional seceonds.
//	time_t iseconds = mlr_arch_timegm(&tm);
//	//printf("ISECONDS %u\n", (unsigned)iseconds);
//	double fseconds = iseconds + fractional_seconds;
//	//printf("FSECONDS %.6lf\n", fseconds);
//	printf("%s %s -> %.6lf\n", time_string, format_string, fseconds);
//
//	return 0;
//}
