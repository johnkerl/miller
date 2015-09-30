#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlrutil.h"
#include "cli/argparse.h"

int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

// ----------------------------------------------------------------
// Assumes null-termination of input array
static int compute_argc(char** argv) {
	int argc = 0;
	while (*argv) {
		argv++;
		argc++;
	}
	return argc;
}

static char * test1() {
	int     bflag1 = -4;
	int     bflag2 = -5;
	int     intv1  = -6;
	int     intv2  = -7;
	int     intv3  = -8;
	double  dblv   = -9.5;
	char*   string = NULL;
	slls_t* plist  = NULL;
	ap_state_t* pstate = ap_alloc();

	ap_define_true_flag(pstate,        "-t",   &bflag1);
	ap_define_false_flag(pstate,       "-f",   &bflag2);
	ap_define_int_value_flag(pstate,   "-100", 100, &intv1);
	ap_define_int_value_flag(pstate,   "-200", 200, &intv2);
	ap_define_int_flag(pstate,         "-i",   &intv3);
	ap_define_double_flag(pstate,      "-d",   &dblv);
	ap_define_string_flag(pstate,      "-s",   &string);
	ap_define_string_list_flag(pstate, "-S",   &plist);

	char* argv[] = { "test-verb", NULL };
	int argc = compute_argc(argv);
	char* verb = argv[0];
	int argi = 1;
	mu_assert_lf(ap_parse(pstate, verb, &argi, argc, argv) == TRUE);
	mu_assert_lf(bflag1 == -4);
	mu_assert_lf(bflag2 == -5);
	mu_assert_lf(intv1 == -6);
	mu_assert_lf(intv2 == -7);
	mu_assert_lf(intv3 == -8);
	mu_assert_lf(dblv == -9.5);
	mu_assert_lf(string == NULL);
	mu_assert_lf(plist == NULL);
	mu_assert_lf(argi == 1);

	ap_free(pstate);
	return 0;
}

static char * test2() {
	int     bflag1 = -4;
	int     bflag2 = -5;
	int     intv1  = -6;
	int     intv2  = -7;
	int     intv3  = -8;
	double  dblv   = -9.5;
	char*   string = NULL;
	slls_t* plist  = NULL;
	ap_state_t* pstate = ap_alloc();

	ap_define_true_flag(pstate,        "-t",   &bflag1);
	ap_define_false_flag(pstate,       "-f",   &bflag2);
	ap_define_int_value_flag(pstate,   "-100", 100, &intv1);
	ap_define_int_value_flag(pstate,   "-200", 200, &intv2);
	ap_define_int_flag(pstate,         "-i",   &intv3);
	ap_define_double_flag(pstate,      "-d",   &dblv);
	ap_define_string_flag(pstate,      "-s",   &string);
	ap_define_string_list_flag(pstate, "-S",   &plist);

	char* argv[] = {
		"test-verb",
		"-t",
		"-f",
		"-100",
		"-200",
		"-i", "555",
		"-d", "4.25",
		"-s", "hello",
		"-S", mlr_strdup_or_die("a,b,c,d,e"),
		"do", "re", "mi",
		NULL
	};
	int argc = compute_argc(argv);
	char* verb = argv[0];
	int argi = 1;
	mu_assert_lf(ap_parse(pstate, verb, &argi, argc, argv) == TRUE);
	mu_assert_lf(bflag1 == TRUE);
	mu_assert_lf(bflag2 == FALSE);
	mu_assert_lf(intv1 == 100);
	mu_assert_lf(intv2 == 200);
	mu_assert_lf(intv3 == 555);
	mu_assert_lf(dblv == 4.25);
	mu_assert_lf(string != NULL);
	mu_assert_lf(streq(string, "hello"));
	mu_assert_lf(plist != NULL);
	mu_assert_lf(slls_equals(plist, slls_from_line(mlr_strdup_or_die("a,b,c,d,e"), ',', FALSE)));
	mu_assert_lf(argi == 13);

	ap_free(pstate);
	return 0;
}

// ================================================================
static char * all_tests() {
	mu_run_test(test1);
	mu_run_test(test2);
	return 0;
}

int main(int argc, char **argv) {
	printf("TEST_ARGPARSE ENTER\n");
	char *result = all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_ARGPARSE: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
