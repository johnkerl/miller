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
static void test_case(
	char*  test_name,
	char** strings,
	int    num_strings,
	char*  buf,
	int*   prc,
	int*   pstridx,
	int*   pmatchlen)
{
	int stridx, matchlen, rc;

	parse_trie_t* ptrie = parse_trie_alloc();
	printf("%s %s\n", sep, test_name);
	parse_trie_print(ptrie);
	for (stridx = 0; stridx < num_strings; stridx++) {
		printf("Adding string[%d] = %s\n", stridx, strings[stridx]);
		parse_trie_add_string(ptrie, strings[stridx], stridx);
		parse_trie_print(ptrie);
	}

	stridx = -2;
	matchlen = -2;
	rc = parse_trie_match(ptrie, buf, strlen(buf), &stridx, &matchlen);

	parse_trie_free(ptrie);

	printf("buf      = %s\n", buf);
	printf("rc       = %d\n", rc);
	printf("stridx   = %d\n", stridx);
	printf("matchlen = %d\n", matchlen);

	*prc       = rc;
	*pstridx   = stridx;
	*pmatchlen = matchlen;
}

// ----------------------------------------------------------------
static char* test_new() {
	int stridx, matchlen, rc;
	char* strings[] = {
		"a"
	};
	int num_strings = sizeof(strings) / sizeof(strings[0]);
	char* buf = "a";

	test_case("simplest", strings, num_strings, buf, &rc, &stridx, &matchlen);

	mu_assert_lf(rc == TRUE);
	mu_assert_lf(stridx == 0);
	mu_assert_lf(matchlen == 1);

	return 0;
}

// ----------------------------------------------------------------
static char* test_simplest() {
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
static char* test_disjoint() {
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
static char* test_short_long() {
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
static char* test_long_short() {
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
static char* all_tests() {
	mu_run_test(test_new);
	mu_run_test(test_simplest);
	mu_run_test(test_disjoint);
	mu_run_test(test_short_long);
	mu_run_test(test_long_short);
	return 0;
}

int main(int argc, char** argv) {
	char* result = all_tests();
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
