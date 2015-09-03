#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlrutil.h"
#include "containers/parse_trie.h"

#ifdef __TEST_PARSE_TRIE_MAIN__
int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

static char* sep = "================================================================";

// ----------------------------------------------------------------
static char * test_simplest() {
	parse_trie_t* ptrie = parse_trie_alloc();
	parse_trie_print(ptrie);
	parse_trie_add_string(ptrie, "a", 0);
	parse_trie_print(ptrie);

	char* buf = "a";
	int stridx = -2;
	int matchlen = -2;
	int rc = parse_trie_match(ptrie, buf, strlen(buf), &stridx, &matchlen);
	printf("buf      = %s\n", buf);
	printf("rc       = %d\n", rc);
	printf("stridx   = %d\n", stridx);
	printf("matchlen = %d\n", matchlen);
	mu_assert_lf(rc == TRUE);
	mu_assert_lf(stridx == 0);
	mu_assert_lf(matchlen == 1);

	parse_trie_free(ptrie);
	return 0;
}

// ----------------------------------------------------------------
static char * test_disjoint() {
	parse_trie_t* ptrie = parse_trie_alloc();
	parse_trie_print(ptrie);
	parse_trie_add_string(ptrie, "abc", 0);
	parse_trie_print(ptrie);
	parse_trie_add_string(ptrie, "fg", 1);
	parse_trie_print(ptrie);

	char* buf = "abcde";
	int stridx = -2;
	int matchlen = -2;
	int rc = parse_trie_match(ptrie, buf, strlen(buf), &stridx, &matchlen);
	printf("buf      = %s\n", buf);
	printf("rc       = %d\n", rc);
	printf("stridx   = %d\n", stridx);
	printf("matchlen = %d\n", matchlen);
	mu_assert_lf(rc == TRUE);
	mu_assert_lf(stridx == 0);
	mu_assert_lf(matchlen == 3);

	parse_trie_free(ptrie);
	return 0;
}

// ----------------------------------------------------------------
static char * test_short_long() {
	parse_trie_t* ptrie = parse_trie_alloc();
	parse_trie_print(ptrie);
	parse_trie_add_string(ptrie, "a", 0);
	parse_trie_print(ptrie);
	parse_trie_add_string(ptrie, "aa", 1);
	parse_trie_print(ptrie);

	char* buf = "aaabc";
	int stridx = -2;
	int matchlen = -2;
	int rc = parse_trie_match(ptrie, buf, strlen(buf), &stridx, &matchlen);
	printf("buf      = %s\n", buf);
	printf("rc       = %d\n", rc);
	printf("stridx   = %d\n", stridx);
	printf("matchlen = %d\n", matchlen);
	mu_assert_lf(rc == TRUE);
	mu_assert_lf(stridx == 0);
	mu_assert_lf(matchlen == 3);

	parse_trie_free(ptrie);
	return 0;
}

// ----------------------------------------------------------------
static char * test_long_short() {
	parse_trie_t* ptrie = parse_trie_alloc();
	parse_trie_print(ptrie);
	parse_trie_add_string(ptrie, "aa", 0);
	parse_trie_print(ptrie);
	parse_trie_add_string(ptrie, "a", 1);
	parse_trie_print(ptrie);

	char* buf = "aaabc";
	int stridx = -2;
	int matchlen = -2;
	int rc = parse_trie_match(ptrie, buf, strlen(buf), &stridx, &matchlen);
	printf("buf      = %s\n", buf);
	printf("rc       = %d\n", rc);
	printf("stridx   = %d\n", stridx);
	printf("matchlen = %d\n", matchlen);
	mu_assert_lf(rc == TRUE);
	mu_assert_lf(stridx == 0);
	mu_assert_lf(matchlen == 3);

	parse_trie_free(ptrie);
	return 0;
}

// ================================================================
static char * all_tests() {
	printf("%s\n", sep); mu_run_test(test_simplest);
	printf("%s\n", sep); mu_run_test(test_disjoint);
	printf("%s\n", sep); mu_run_test(test_short_long);
	printf("%s\n", sep); mu_run_test(test_long_short);
	printf("%s\n", sep);
	return 0;
}

int main(int argc, char **argv) {
	char *result = all_tests();
	printf("\n");
	if (result != 0) {
		//printf("%s\n", result);
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_PARSE_TRIE: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
#endif // __TEST_PARSE_TRIE_MAIN__
