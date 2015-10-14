#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlrutil.h"
#include "lib/string_builder.h"

int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

// ----------------------------------------------------------------
static char * test_simple() {
	string_builder_t* psb = sb_alloc(1);

	sb_init(psb, 1);
	mu_assert("error: case 0", streq("", sb_finish(psb)));

	sb_init(psb, 1);
	sb_append_char(psb, 'a');
	mu_assert("error: case 1", streq("a", sb_finish(psb)));

	sb_init(psb, 1);
	sb_append_char(psb, 'a');
	sb_append_char(psb, 'b');
	mu_assert("error: case 2", streq("ab", sb_finish(psb)));

	sb_init(psb, 1);
	sb_append_char(psb, 'a');
	sb_append_char(psb, 'b');
	sb_append_char(psb, 'b');
	sb_append_char(psb, 'c');
	sb_append_char(psb, 'c');
	sb_append_char(psb, 'e');
	mu_assert("error: case 3", streq("abbcce", sb_finish(psb)));

	sb_init(psb, 1);
	sb_append_string(psb, "");
	mu_assert("error: case 4", streq("", sb_finish(psb)));

	sb_init(psb, 1);
	sb_append_string(psb, "hello");
	mu_assert("error: case 5", streq("hello", sb_finish(psb)));

	sb_init(psb, 1);
	sb_append_string(psb, "hello");
	sb_append_char(psb, ',');
	sb_append_char(psb, ' ');
	sb_append_string(psb, "world");
	sb_append_char(psb, '!');
	mu_assert("error: case 6", streq("hello, world!", sb_finish(psb)));

	sb_init(psb, 2);
	sb_append_string(psb, "hello");
	sb_append_char(psb, ',');
	sb_append_char(psb, ' ');
	sb_append_string(psb, "world");
	sb_append_char(psb, '!');
	mu_assert("error: case 7", streq("hello, world!", sb_finish(psb)));

	sb_init(psb, 32768);
	sb_append_string(psb, "hello");
	sb_append_char(psb, ',');
	sb_append_char(psb, ' ');
	sb_append_string(psb, "world");
	sb_append_char(psb, '!');
	mu_assert("error: case 8", streq("hello, world!", sb_finish(psb)));

	return 0;
}

// ================================================================
static char * all_tests() {
	mu_run_test(test_simple);
	return 0;
}

int main(int argc, char **argv) {
	printf("TEST_STRING_BUILDER ENTER\n");
	char *result = all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_STRING_BUILDER: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
