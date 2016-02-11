#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlrutil.h"
#include "input/json_parser.h"

int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

// ----------------------------------------------------------------
static void try_out(char* input) {
	json_char* end = NULL;
	json_value_t* pvalue = json_parse_for_unit_test((json_char*)input, &end);
	json_print_recursive(pvalue);
	// xxx make a dump-node method: recursive version
}

// ----------------------------------------------------------------
static char * test_numbers_only() {
	json_char* input = "123";
	json_char* end   = NULL;
	json_value_t* pvalue = json_parse_for_unit_test(input, &end);
	mu_assert_lf(pvalue != NULL);
	mu_assert_lf(pvalue->type == JSON_INTEGER);
	json_print_recursive(pvalue);
	return 0;
}

// ----------------------------------------------------------------
static char * test_array() {
	json_char* input = "[ 1238139939387114097,  1213418874328430463 ]";
	json_char* end   = NULL;
	json_value_t* pvalue = json_parse_for_unit_test(input, &end);
	mu_assert_lf(pvalue != NULL);
	mu_assert_lf(pvalue->type == JSON_ARRAY);
	// xxx more
	json_print_recursive(pvalue);
	return 0;
}

// ================================================================
static char * all_tests() {
	mu_run_test(test_numbers_only);
	mu_run_test(test_array);
	return 0;
}

int main(int argc, char **argv) {
	if (argc > 1) { // Manual mode
		for (int argi = 1; argi < argc; argi++) {
			try_out(argv[argi]);
		}
	} else { // Unit-test mode
		printf("TEST_JSON_PARSER ENTER\n");
		char *result = all_tests();
		printf("\n");
		if (result != 0) {
			printf("Not all unit tests passed\n");
		}
		else {
			printf("TEST_JSON_PARSER: ALL UNIT TESTS PASSED\n");
		}
		printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
		printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

		return result != 0;
	}
}
