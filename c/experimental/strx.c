#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/time.h>
#include <sys/stat.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/mlr_arch.h"
#include "lib/string_builder.h"

#define MYBUFLEN 256

// ----------------------------------------------------------------
static int f(char* float_string, char* format_string) {
	double fseconds;
	if (sscanf(float_string, "%lf", &fseconds) != 1) {
		fprintf(stderr, "!dbl b04k!\n");
		return 1;
	}

	time_t iseconds = (time_t) fseconds;
	double fracsec = fseconds - iseconds;

	struct tm tm = *gmtime(&iseconds);

	// xxx see if %nS is a substring
	// if not then do it like always (w/ malloc)
	// if so:
	// * copy subformats left & right of the %nS
	// * strftime each
	// * myself format the  middle
	// * sb_append

	char* middle_format = NULL;
	char* right_subformat = NULL;
	for (char* p = format_string; *p; p++) {
		if (p[0] == '%' && p[1] >= '1' && p[1] <= '9' && p[2] == 'S') {
			middle_format = p;
			right_subformat = &p[3];
			break;
		}
	}

	if (middle_format == NULL) {
		char* output_string = mlr_malloc_or_die(MYBUFLEN);

		int written_length = strftime(output_string, MYBUFLEN, format_string, &tm);
		if (written_length > MYBUFLEN || written_length == 0) {
			fprintf(stderr, "Could not strftime(\"%s\", \"%s\").\n",
			float_string, format_string);
			return 1;
		} else {
			printf("%s %s -> %s\n", float_string, format_string, output_string);
		}

		free(output_string);

	} else {
		int left_subformat_length = middle_format - format_string;
		char* left_subformat = mlr_malloc_or_die(left_subformat_length + 1);
		memcpy(left_subformat, format_string, left_subformat_length);
		left_subformat[left_subformat_length] = 0;
		char left_formatted[MYBUFLEN];
		char right_formatted[MYBUFLEN];

		if (*left_subformat == 0) {
			*left_formatted = 0;
		} else {
			int written_length = strftime(left_formatted, MYBUFLEN, left_subformat, &tm);
			printf("LEFT  SUBFORMAT \"%s\"\n", left_subformat);
			if (written_length > MYBUFLEN || written_length == 0) {
				fprintf(stderr, "Could not strftime(\"%s\", \"%s\").\n",
				float_string, left_subformat);
				return 1;
			}
		}
		free(left_subformat);

		char middle_int_formatted[16];
		char* middle_int_format = "%S";
		printf("MIDDLI SUBFORMAT \"%s\"\n", middle_int_format);
		int written_length = strftime(middle_int_formatted, sizeof(middle_int_formatted), middle_int_format, &tm);
		if (written_length > MYBUFLEN || written_length == 0) {
			fprintf(stderr, "Could not strftime(\"%s\", \"%s\").\n",
			float_string, middle_int_format);
			return 1;
		}

		char middle_sprintf_format[] = "%.xlf";
		char middle_fractional_formatted[16];
		middle_sprintf_format[2] = middle_format[1];
		printf("MIDDLE_SPRINTF_FORMAT \"%s\"\n", middle_sprintf_format);
		sprintf(middle_fractional_formatted, middle_sprintf_format, fracsec);

		if (*right_subformat == 0) {
			*right_formatted = 0;
		} else {
			printf("RIGHT SUBFORMAT \"%s\"\n", right_subformat);
			int written_length = strftime(right_formatted, MYBUFLEN, right_subformat, &tm);
			if (written_length > MYBUFLEN || written_length == 0) {
				fprintf(stderr, "Could not strftime(\"%s\", \"%s\").\n",
				float_string, right_subformat);
				return 1;
			}
		}

		string_builder_t* psb = sb_alloc(32);
		sb_append_string(psb, left_formatted);
		sb_append_string(psb, middle_int_formatted);
		sb_append_string(psb, &middle_fractional_formatted[1]); // xxx comment
		sb_append_string(psb, right_formatted);
		char* output_string = sb_finish(psb);
		sb_free(psb);

		printf("LEFT   OUT \"%s\"\n", left_formatted);
		printf("MIDDLI OUT \"%s\"\n", middle_int_formatted);
		printf("MIDDLF OUT \"%s\"\n", middle_fractional_formatted);
		printf("RIGHT  OUT \"%s\"\n", right_formatted);

		printf("%s %s -> %s\n", float_string, format_string, output_string);
		free(output_string);
	}

	return 0;
}

// ----------------------------------------------------------------
static int p(char* time_string, char* format_string) {

	struct tm tm;
	memset(&tm, 0, sizeof(tm));

	// xxx cmt try the non-floating-point-seconds case first and return quickly if so.
	char* strptime_retval = mlr_arch_strptime(time_string, format_string, &tm);
	if (strptime_retval != NULL) {
		if (*strptime_retval != 0) { // xxx extraneous
			fprintf(stderr, "%s: could not strptime(\"%s\", \"%s\"). See \"%s --help-function strptime\".\n",
				MLR_GLOBALS.bargv0, time_string, format_string, MLR_GLOBALS.bargv0);
			exit(1);
		}
		time_t iseconds = mlr_arch_timegm(&tm);
		printf("%s, %s -> %u, \"%s\"\n", time_string, format_string, (unsigned)iseconds, strptime_retval);
		return 0;
	}

	char* pS = strstr(format_string, "%S");
	if (pS == NULL) {
		// Couldn't have been because of floating-point-seconds stuff. No reason to try any harder.
		fprintf(stderr, "%s: could not strptime(\"%s\", \"%s\"). See \"%s --help-function strptime\".\n",
			MLR_GLOBALS.bargv0, time_string, format_string, MLR_GLOBALS.bargv0);
		exit(1);
	}

	// xxx
	// Input is    "2017-04-09T00:51:09.123456 TZBLAHBLAH"
	// with format "%Y-%m-%dT%H:%M:%S TZBLAHBLAH"

	// 1. Copy the format up to the %S but with nothing else after. This is temporary to help us locate
	//    the fractional-seconds part of the input string.
	//    Example temporary format: "%Y-%m-%dT%H:%M:%S"

	int truncated_format_string_length = pS - format_string + 2;
	char* truncated_format_string = mlr_malloc_or_die(truncated_format_string_length + 1);
	memcpy(truncated_format_string, format_string, truncated_format_string_length);
	truncated_format_string[truncated_format_string_length] = 0;
	printf("ORIGINAL  FORMAT \"%s\"\n", format_string);
	printf("TRUNCATED FORMAT \"%s\"\n", truncated_format_string);

	// 2. strptime using that truncated format and ignore the tm. Only look at the string return value.
	//    Example return value: ".123456 TZBLAHBLAH"

	strptime_retval = mlr_arch_strptime(time_string, truncated_format_string, &tm);
	if (strptime_retval == NULL) {
		fprintf(stderr, "%s: could not strptime(\"%s\", \"%s\"). See \"%s --help-function strptime\".\n",
			MLR_GLOBALS.bargv0, time_string, format_string, MLR_GLOBALS.bargv0);
		exit(1);
	}
	printf("STRPTIME TEMP RETVAL \"%s\"\n", strptime_retval);

	// 3. strtod the return value to find the fractional-seconds part, and whatever's after.
	//    Example fractional-seconds part: ".123456"
	//    Example stuff after: " TZBLAHBLAH"

	char* stuff_after = NULL;
	double fractional_seconds = strtod(strptime_retval, &stuff_after);
	printf("FRACTIONAL SECONDS %.6lf\n", fractional_seconds);
	printf("STUFF AFTER  \"%s\"\n", stuff_after);

	// 4. Make a copy of the input string with the fractional seconds elided.
	//    Example: "2017-04-09T00:51:09 TZBLAHBLAH"
	char* elided_fraction_input = mlr_malloc_or_die(strlen(time_string) + 1);
	int input_length_to_fractional_seconds = strptime_retval - time_string;
	memcpy(elided_fraction_input, time_string, input_length_to_fractional_seconds);
	strcpy(&elided_fraction_input[input_length_to_fractional_seconds], stuff_after);
	printf("ELIDE \"%s\"\n", elided_fraction_input);

	// 5. strptime the elided-fraction input string using the original format string. Get the tm.
	memset(&tm, 0, sizeof(tm));
	strptime_retval = mlr_arch_strptime(elided_fraction_input, format_string, &tm);
	if (strptime_retval == NULL) {
		fprintf(stderr, "%s: could not strptime(\"%s\", \"%s\"). See \"%s --help-function strptime\".\n",
			MLR_GLOBALS.bargv0, time_string, format_string, MLR_GLOBALS.bargv0);
		exit(1);
	}
	printf("STRPTIME ELIDE RETVAL \"%s\"\n", strptime_retval);

	// 6. Convert the tm to a time_t (seconds since the epoch) and then add the fractional seceonds.
	time_t iseconds = mlr_arch_timegm(&tm);
	printf("ISECONDS %u\n", (unsigned)iseconds);
	double fseconds = iseconds + fractional_seconds;
	printf("FSECONDS %.6lf\n", fseconds);
	printf("%s %s -> %.6lf\n", time_string, format_string, fseconds);

	return 0;
}

//----------------------------------------------------------------------
int main(int argc, char **argv) {
	if (argc != 4) {
		fprintf(stderr, "c!=4 b04k!\n");
		exit(1);
	}
	if (!strcmp(argv[1], "f")) {
		return f(argv[2], argv[3]);
	} else if (!strcmp(argv[1], "p")) {
		return p(argv[2], argv[3]);
	} else {
		fprintf(stderr, "f/p b04k!\n");
		exit(1);
	}

	return 0;
}
