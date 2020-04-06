#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"

int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

// ----------------------------------------------------------------
static char * test_canonical_mod() {
	mu_assert("error: canonical_mod -7", mlr_canonical_mod(-7, 5) == 3);
	mu_assert("error: canonical_mod -6", mlr_canonical_mod(-6, 5) == 4);
	mu_assert("error: canonical_mod -5", mlr_canonical_mod(-5, 5) == 0);
	mu_assert("error: canonical_mod -4", mlr_canonical_mod(-4, 5) == 1);
	mu_assert("error: canonical_mod -3", mlr_canonical_mod(-3, 5) == 2);
	mu_assert("error: canonical_mod -2", mlr_canonical_mod(-2, 5) == 3);
	mu_assert("error: canonical_mod -1", mlr_canonical_mod(-1, 5) == 4);
	mu_assert("error: canonical_mod  0", mlr_canonical_mod(0, 5) == 0);
	mu_assert("error: canonical_mod  1", mlr_canonical_mod(1, 5) == 1);
	mu_assert("error: canonical_mod  2", mlr_canonical_mod(2, 5) == 2);
	mu_assert("error: canonical_mod  3", mlr_canonical_mod(3, 5) == 3);
	mu_assert("error: canonical_mod  4", mlr_canonical_mod(4, 5) == 4);
	mu_assert("error: canonical_mod  5", mlr_canonical_mod(5, 5) == 0);
	mu_assert("error: canonical_mod  6", mlr_canonical_mod(6, 5) == 1);
	mu_assert("error: canonical_mod  7", mlr_canonical_mod(7, 5) == 2);
	return 0;
}

// ----------------------------------------------------------------
static char * test_power_of_two_above() {
	mu_assert("error: power_of_two_above 0", power_of_two_above(0) == 1);
	mu_assert("error: power_of_two_above 1", power_of_two_above(1) == 2);
	mu_assert("error: power_of_two_above 2", power_of_two_above(2) == 4);
	mu_assert("error: power_of_two_above 3", power_of_two_above(3) == 4);
	mu_assert("error: power_of_two_above 4", power_of_two_above(4) == 8);
	mu_assert("error: power_of_two_above 5", power_of_two_above(5) == 8);
	mu_assert("error: power_of_two_above 6", power_of_two_above(6) == 8);
	mu_assert("error: power_of_two_above 7", power_of_two_above(7) == 8);
	mu_assert("error: power_of_two_above 8", power_of_two_above(8) == 16);
	mu_assert("error: power_of_two_above 9", power_of_two_above(9) == 16);
	mu_assert("error: power_of_two_above 1023", power_of_two_above(1023) == 1024);
	mu_assert("error: power_of_two_above 1024", power_of_two_above(1024) == 2048);
	mu_assert("error: power_of_two_above 1025", power_of_two_above(1025) == 2048);
	return 0;
}

// ----------------------------------------------------------------
static char * test_streq() {
	char* x;
	char* y;

	x = "";   y = "";   mu_assert_lf(streq(x, y) == !strcmp(x, y));
	x = "";   y = "1";  mu_assert_lf(streq(x, y) == !strcmp(x, y));
	x = "";   y = "12"; mu_assert_lf(streq(x, y) == !strcmp(x, y));

	x = "";   y = "";   mu_assert_lf(streq(x, y) == !strcmp(x, y));
	x = "a";  y = "";   mu_assert_lf(streq(x, y) == !strcmp(x, y));
	x = "ab"; y = "";   mu_assert_lf(streq(x, y) == !strcmp(x, y));

	x = "1";  y = "";   mu_assert_lf(streq(x, y) == !strcmp(x, y));
	x = "1";  y = "1";  mu_assert_lf(streq(x, y) == !strcmp(x, y));
	x = "1";  y = "12"; mu_assert_lf(streq(x, y) == !strcmp(x, y));

	x = "12"; y = "";   mu_assert_lf(streq(x, y) == !strcmp(x, y));
	x = "12"; y = "1";  mu_assert_lf(streq(x, y) == !strcmp(x, y));
	x = "12"; y = "12"; mu_assert_lf(streq(x, y) == !strcmp(x, y));

	x = "";   y = "a";  mu_assert_lf(streq(x, y) == !strcmp(x, y));
	x = "a";  y = "a";  mu_assert_lf(streq(x, y) == !strcmp(x, y));
	x = "ab"; y = "a";  mu_assert_lf(streq(x, y) == !strcmp(x, y));

	x = "";   y = "ab"; mu_assert_lf(streq(x, y) == !strcmp(x, y));
	x = "a";  y = "ab"; mu_assert_lf(streq(x, y) == !strcmp(x, y));
	x = "ab"; y = "ab"; mu_assert_lf(streq(x, y) == !strcmp(x, y));

	return 0;
}

// ----------------------------------------------------------------
static char * test_streqn() {
	char* x;
	char* y;

	x = "";    y = "";    mu_assert_lf(streqn(x, y, 2) == !strncmp(x, y, 2));
	x = "";    y = "1";   mu_assert_lf(streqn(x, y, 2) == !strncmp(x, y, 2));
	x = "";    y = "12";  mu_assert_lf(streqn(x, y, 2) == !strncmp(x, y, 2));
	x = "";    y = "123"; mu_assert_lf(streqn(x, y, 2) == !strncmp(x, y, 2));

	x = "";    y = "";    mu_assert_lf(streqn(x, y, 2) == !strncmp(x, y, 2));
	x = "a";   y = "";    mu_assert_lf(streqn(x, y, 2) == !strncmp(x, y, 2));
	x = "ab";  y = "";    mu_assert_lf(streqn(x, y, 2) == !strncmp(x, y, 2));
	x = "abc"; y = "";    mu_assert_lf(streqn(x, y, 2) == !strncmp(x, y, 2));

	x = "a";   y = "a";   mu_assert_lf(streqn(x, y, 2) == !strncmp(x, y, 2));
	x = "a";   y = "aa";  mu_assert_lf(streqn(x, y, 2) == !strncmp(x, y, 2));
	x = "a";   y = "ab";  mu_assert_lf(streqn(x, y, 2) == !strncmp(x, y, 2));
	x = "a";   y = "abd"; mu_assert_lf(streqn(x, y, 2) == !strncmp(x, y, 2));

	x = "ab";  y = "a";   mu_assert_lf(streqn(x, y, 2) == !strncmp(x, y, 2));
	x = "ab";  y = "ab";  mu_assert_lf(streqn(x, y, 2) == !strncmp(x, y, 2));
	x = "ab";  y = "abd"; mu_assert_lf(streqn(x, y, 2) == !strncmp(x, y, 2));

	x = "abc"; y = "a";   mu_assert_lf(streqn(x, y, 2) == !strncmp(x, y, 2));
	x = "abc"; y = "ab";  mu_assert_lf(streqn(x, y, 2) == !strncmp(x, y, 2));
	x = "abc"; y = "abc"; mu_assert_lf(streqn(x, y, 2) == !strncmp(x, y, 2));
	x = "abc"; y = "abd"; mu_assert_lf(streqn(x, y, 2) == !strncmp(x, y, 2));

	return 0;
}

// ----------------------------------------------------------------
static char * test_strdup_quoted() {

	mu_assert_lf(streq(mlr_strdup_quoted_or_die(""), "\"\""));
	mu_assert_lf(streq(mlr_strdup_quoted_or_die("x"), "\"x\""));
	mu_assert_lf(streq(mlr_strdup_quoted_or_die("xy"), "\"xy\""));
	mu_assert_lf(streq(mlr_strdup_quoted_or_die("xyz"), "\"xyz\""));

	return 0;
}

// ----------------------------------------------------------------
static char * test_starts_or_ends_with() {

	mu_assert_lf(string_starts_with("abcde", ""));
	mu_assert_lf(string_starts_with("abcde", "a"));
	mu_assert_lf(string_starts_with("abcde", "abcd"));
	mu_assert_lf(string_starts_with("abcde", "abcde"));
	mu_assert_lf(!string_starts_with("abcde", "abcdef"));

	mu_assert_lf(string_ends_with("abcde", "", NULL));
	mu_assert_lf(string_ends_with("abcde", "e", NULL));
	mu_assert_lf(string_ends_with("abcde", "de", NULL));
	mu_assert_lf(string_ends_with("abcde", "abcde", NULL));
	mu_assert_lf(!string_ends_with("abcde", "0abcde", NULL));
	int len = -1;
	mu_assert_lf(!string_ends_with("abcde", "0abcde", &len));
	mu_assert_lf(len == 5);

	return 0;
}

// ----------------------------------------------------------------
static char * test_scanners() {
	mu_assert("error: mlr_alloc_string_from_double", streq(mlr_alloc_string_from_double(4.25, "%.4f"), "4.2500"));
	mu_assert("error: mlr_alloc_string_from_ull", streq(mlr_alloc_string_from_ull(12345LL), "12345"));
	mu_assert("error: mlr_alloc_string_from_int", streq(mlr_alloc_string_from_int(12345), "12345"));
	return 0;
}

// ----------------------------------------------------------------
static char * test_paste() {
	mu_assert("error: paste 2", streq(mlr_paste_2_strings("ab", "cd"), "abcd"));
	mu_assert("error: paste 3", streq(mlr_paste_3_strings("ab", "cd", "ef"), "abcdef"));
	mu_assert("error: paste 4", streq(mlr_paste_4_strings("ab", "cd", "ef", "gh"), "abcdefgh"));
	mu_assert("error: paste 5", streq(mlr_paste_5_strings("ab", "cd", "ef", "gh", "ij"), "abcdefghij"));
	return 0;
}

// ----------------------------------------------------------------
static char * test_unbackslash() {
	mu_assert_lf(streq(mlr_alloc_unbackslash(""), ""));
	mu_assert_lf(streq(mlr_alloc_unbackslash("hello"), "hello"));
	mu_assert_lf(streq(mlr_alloc_unbackslash("\\r\\n"), "\r\n"));
	mu_assert_lf(streq(mlr_alloc_unbackslash("\\t\\\\"), "\t\\"));
	mu_assert_lf(streq(mlr_alloc_unbackslash("[\\132]"), "[Z]"));
	mu_assert_lf(streq(mlr_alloc_unbackslash("[\\x59]"), "[Y]"));
	return 0;
}

// ----------------------------------------------------------------
static char * test_rstrip() {

	char* a = NULL;
	mlr_rstrip(a);
	mu_assert_lf(a == NULL);

	a = mlr_strdup_or_die("");
	mlr_rstrip(a);
	mu_assert_lf(streq(a, ""));

	a = mlr_strdup_or_die("foo");
	mlr_rstrip(a);
	mu_assert_lf(streq(a, "foo"));

	a = mlr_strdup_or_die("\r");
	mlr_rstrip(a);
	mu_assert_lf(streq(a, ""));

	a = mlr_strdup_or_die("\n");
	mlr_rstrip(a);
	mu_assert_lf(streq(a, ""));

	a = mlr_strdup_or_die("\r\n");
	mlr_rstrip(a);
	mu_assert_lf(streq(a, ""));

	a = mlr_strdup_or_die("x\r");
	mlr_rstrip(a);
	mu_assert_lf(streq(a, "x"));

	a = mlr_strdup_or_die("x\n");
	mlr_rstrip(a);
	mu_assert_lf(streq(a, "x"));

	a = mlr_strdup_or_die("x\r\n");
	mlr_rstrip(a);
	mu_assert_lf(streq(a, "x"));

	return 0;
}

// ================================================================
static char * all_tests() {
	mu_run_test(test_canonical_mod);
	mu_run_test(test_power_of_two_above);
	mu_run_test(test_streq);
	mu_run_test(test_streqn);
mu_run_test(test_strdup_quoted);
	mu_run_test(test_starts_or_ends_with);
	mu_run_test(test_scanners);
	mu_run_test(test_paste);
	mu_run_test(test_unbackslash);
	mu_run_test(test_rstrip);
	return 0;
}

int main(int argc, char **argv) {
	mlr_global_init(argv[0], NULL);
	printf("TEST_MLRUTIL ENTER\n");
	char *result = all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_MLRUTIL: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
