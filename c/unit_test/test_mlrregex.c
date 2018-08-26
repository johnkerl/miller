#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/mlrregex.h"

int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

// ----------------------------------------------------------------
static char * test_save_regex_captures() {

	// Capture-groups \1 through \9 supported, along with entire-string match in \0
	const size_t nmatchmax = 10;
	regmatch_t matches[nmatchmax];
	string_array_t* pregex_captures = NULL;
	regex_t regex;

	char* input  = "abcde";
	char* sregex = "abcde";
	regcomp_or_die(&regex, sregex, 0);
	regmatch_or_die(&regex, input, nmatchmax, matches);
	save_regex_captures(&pregex_captures, input, matches, nmatchmax);
	mu_assert_lf(pregex_captures != NULL);
	mu_assert_lf(pregex_captures->length == 1);
	mu_assert_lf(pregex_captures->strings[0] != NULL);
	mu_assert_lf(streq(pregex_captures->strings[0], "abcde"));
	regfree(&regex);

	input  = "abcde";
	sregex = "a(.*)e";
	regcomp_or_die(&regex, sregex, 0);
	regmatch_or_die(&regex, input, nmatchmax, matches);
	save_regex_captures(&pregex_captures, input, matches, nmatchmax);
	mu_assert_lf(pregex_captures->length == 2);
	mu_assert_lf(pregex_captures->strings[0] != NULL);
	mu_assert_lf(streq(pregex_captures->strings[1], "bcd"));
	regfree(&regex);

	input  = "abcde";
	sregex = "a(b)(.)(d)e";
	regcomp_or_die(&regex, sregex, 0);
	regmatch_or_die(&regex, input, nmatchmax, matches);
	save_regex_captures(&pregex_captures, input, matches, nmatchmax);
	mu_assert_lf(pregex_captures->length == 4);
	mu_assert_lf(pregex_captures->strings[0] != NULL);
	mu_assert_lf(streq(pregex_captures->strings[1], "b"));
	mu_assert_lf(streq(pregex_captures->strings[2], "c"));
	mu_assert_lf(streq(pregex_captures->strings[3], "d"));
	regfree(&regex);

	input  = "abcdefghij";
	sregex = "(a)(b)(c)(d)(e)(f)(g)(h)(i)";
	regcomp_or_die(&regex, sregex, 0);
	regmatch_or_die(&regex, input, nmatchmax, matches);
	save_regex_captures(&pregex_captures, input, matches, nmatchmax);
	mu_assert_lf(pregex_captures->length == 10);
	mu_assert_lf(pregex_captures->strings[0] != NULL);
	mu_assert_lf(streq(pregex_captures->strings[1], "a"));
	mu_assert_lf(streq(pregex_captures->strings[2], "b"));
	mu_assert_lf(streq(pregex_captures->strings[3], "c"));
	mu_assert_lf(streq(pregex_captures->strings[4], "d"));
	mu_assert_lf(streq(pregex_captures->strings[5], "e"));
	mu_assert_lf(streq(pregex_captures->strings[6], "f"));
	mu_assert_lf(streq(pregex_captures->strings[7], "g"));
	mu_assert_lf(streq(pregex_captures->strings[8], "h"));
	mu_assert_lf(streq(pregex_captures->strings[9], "i"));
	regfree(&regex);

	string_array_free(pregex_captures);

	return 0;
}

// ----------------------------------------------------------------
static char * test_interpolate_regex_captures() {
	int was_allocated = FALSE;

	char* output = interpolate_regex_captures("hello", NULL, &was_allocated);
	mu_assert_lf(streq(output, "hello"));
	mu_assert_lf(was_allocated == FALSE);

	string_array_t* psa = string_array_alloc(0);
	output = interpolate_regex_captures("hello", psa, &was_allocated);
	mu_assert_lf(streq(output, "hello"));
	mu_assert_lf(was_allocated == FALSE);
	string_array_free(psa);

	psa = string_array_from_line(mlr_strdup_or_die("X,a,b,c"), ',');
	output = interpolate_regex_captures("hello", psa, &was_allocated);
	mu_assert_lf(streq(output, "hello"));
	mu_assert_lf(was_allocated == FALSE);
	string_array_free(psa);

	psa = string_array_from_line(mlr_strdup_or_die("X,a,b,c"), ',');
	output = interpolate_regex_captures("h\\3ello", psa, &was_allocated);
	printf("output=[%s]\n", output);
	mu_assert_lf(streq(output, "hcello"));
	mu_assert_lf(was_allocated == TRUE);
	string_array_free(psa);

	psa = string_array_from_line(mlr_strdup_or_die("X,a,b,c"), ',');
	output = interpolate_regex_captures("h\\1ello", psa, &was_allocated);
	printf("output=[%s]\n", output);
	mu_assert_lf(streq(output, "haello"));
	mu_assert_lf(was_allocated == TRUE);
	string_array_free(psa);

	psa = string_array_from_line(mlr_strdup_or_die("X,a,b,c"), ',');
	output = interpolate_regex_captures("h\\4ello", psa, &was_allocated);
	printf("output=[%s]\n", output);
	mu_assert_lf(streq(output, "hello"));
	mu_assert_lf(was_allocated == TRUE);
	string_array_free(psa);

	psa = string_array_from_line(mlr_strdup_or_die("X,a,b,c"), ',');
	output = interpolate_regex_captures("h\\0ello", psa, &was_allocated);
	printf("output=[%s]\n", output);
	mu_assert_lf(streq(output, "hXello"));
	mu_assert_lf(was_allocated == TRUE);
	string_array_free(psa);

	psa = string_array_from_line(mlr_strdup_or_die("X,a,b,c"), ',');
	output = interpolate_regex_captures("h\\3e\\1l\\2l\\4o", psa, &was_allocated);
	printf("output=[%s]\n", output);
	mu_assert_lf(streq(output, "hcealblo"));
	mu_assert_lf(was_allocated == TRUE);
	string_array_free(psa);

	return 0;
}

// ----------------------------------------------------------------
static char * test_regextract() {
	char* input = NULL;
	char* sregex = NULL;
	char* output = NULL;
	regex_t regex;
	int cflags = 0;

	input = "abcdef";
	sregex = ".+";
	regcomp_or_die(&regex, sregex, cflags);
	output = regextract(input, &regex);
	mu_assert_lf(output != NULL);
	mu_assert_lf(streq(output, input));
	printf("regextract input=\"%s\" regex=\"%s\" output=\"%s\"\n", input, sregex, output);
	free(output);

	input = "abcdef";
	sregex = "[a-z]+";
	regcomp_or_die(&regex, sregex, cflags);
	output = regextract(input, &regex);
	mu_assert_lf(output != NULL);
	mu_assert_lf(streq(output, input));
	printf("regextract input=\"%s\" regex=\"%s\" output=\"%s\"\n", input, sregex, output);
	free(output);

	input = "abcdef";
	sregex = "[0-9]+";
	regcomp_or_die(&regex, sregex, cflags);
	output = regextract(input, &regex);
	mu_assert_lf(output == NULL);
	printf("regextract input=\"%s\" regex=\"%s\" output=NULL\n", input, sregex);
	free(output);

	input = "abc345";
	sregex = "[0-9]+";
	regcomp_or_die(&regex, sregex, cflags);
	output = regextract(input, &regex);
	printf("regextract input=\"%s\" regex=\"%s\" output=\"%s\"\n", input, sregex, output);
	mu_assert_lf(output != NULL);
	mu_assert_lf(streq(output, "345"));
	free(output);

	input = "789xyz";
	sregex = "[0-9]+";
	regcomp_or_die(&regex, sregex, cflags);
	output = regextract(input, &regex);
	printf("regextract input=\"%s\" regex=\"%s\" output=\"%s\"\n", input, sregex, output);
	mu_assert_lf(output != NULL);
	mu_assert_lf(streq(output, "789"));
	free(output);

	return 0;
}

// ----------------------------------------------------------------
static char * test_regextract_or_else() {
	char* input = NULL;
	char* sregex = NULL;
	char* default_value = "DEFAULT";
	char* output = NULL;
	regex_t regex;
	int cflags = 0;

	input = "abcdef";
	sregex = ".+";
	regcomp_or_die(&regex, sregex, cflags);
	output = regextract_or_else(input, &regex, default_value);
	mu_assert_lf(output != NULL);
	mu_assert_lf(streq(output, input));
	printf("regextract_or_else input=\"%s\" regex=\"%s\" default=\"%s\" output=\"%s\"\n", input, sregex, default_value, output);
	free(output);

	input = "abcdef";
	sregex = "[a-z]+";
	regcomp_or_die(&regex, sregex, cflags);
	output = regextract_or_else(input, &regex, default_value);
	mu_assert_lf(output != NULL);
	mu_assert_lf(streq(output, input));
	printf("regextract_or_else input=\"%s\" regex=\"%s\" default=\"%s\" output=\"%s\"\n", input, sregex, default_value, output);
	free(output);

	input = "abcdef";
	sregex = "[0-9]+";
	regcomp_or_die(&regex, sregex, cflags);
	output = regextract_or_else(input, &regex, default_value);
	mu_assert_lf(output != NULL);
	mu_assert_lf(streq(output, default_value));
	printf("regextract_or_else input=\"%s\" regex=\"%s\" default=\"%s\" output=NULL\n", input, sregex, default_value);
	free(output);

	input = "abc345";
	sregex = "[0-9]+";
	regcomp_or_die(&regex, sregex, cflags);
	output = regextract_or_else(input, &regex, default_value);
	printf("regextract_or_else input=\"%s\" regex=\"%s\" default=\"%s\" output=\"%s\"\n", input, sregex, default_value, output);
	mu_assert_lf(output != NULL);
	mu_assert_lf(streq(output, "345"));
	free(output);

	input = "789xyz";
	sregex = "[0-9]+";
	regcomp_or_die(&regex, sregex, cflags);
	output = regextract_or_else(input, &regex, default_value);
	printf("regextract_or_else input=\"%s\" regex=\"%s\" default=\"%s\" output=\"%s\"\n", input, sregex, default_value, output);
	mu_assert_lf(output != NULL);
	mu_assert_lf(streq(output, "789"));
	free(output);

	return 0;
}

// ================================================================
static char * all_tests() {
	mu_run_test(test_save_regex_captures);
	mu_run_test(test_interpolate_regex_captures);
	mu_run_test(test_regextract);
	mu_run_test(test_regextract_or_else);
	return 0;
}

int main(int argc, char **argv) {
	mlr_global_init(argv[0], NULL);
	printf("TEST_MLRREGEX ENTER\n");
	char *result = all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_MLRREGEX: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
