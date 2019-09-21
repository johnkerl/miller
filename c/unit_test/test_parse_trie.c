#include <stdio.h>
#include <string.h>
#include <ctype.h>
#include "lib/minunit.h"
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/parse_trie.h"

int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

static char* sep = "================================================================";


static void print_with_unprintables(char* s) {
	for (char* p = s; *p; p++) {
		char c = *p;
		printf("%c[%02x]", isprint((unsigned char)c) ? c : '?', ((unsigned)c) & 0xff);
	}
}

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
		printf("Adding string[%d] = [", stridx);
		print_with_unprintables(strings[stridx]);
		printf("]\n");
		parse_trie_add_string(ptrie, strings[stridx], stridx);
		parse_trie_print(ptrie);
	}

	stridx = -2;
	matchlen = -2;
	rc = parse_trie_ring_match(ptrie, buf, 0, strlen(buf), 0xff, &stridx, &matchlen);

	parse_trie_free(ptrie);

	printf("buf      = %s\n", buf);
	printf("rc       = %d\n", rc);
	printf("stridx   = %d (%s)\n", stridx, strings[stridx]);
	printf("matchlen = %d\n", matchlen);

	*prc       = rc;
	*pstridx   = stridx;
	*pmatchlen = matchlen;
}

// ----------------------------------------------------------------
static char* test_simplest() {
	char* strings[] = { "a" };
	char* buf = "a";
	int expect_rc = TRUE, expect_stridx = 0, expect_matchlen = 1;

	int num_strings = sizeof(strings) / sizeof(strings[0]);
	int stridx, matchlen, rc;
	test_case("simplest", strings, num_strings, buf, &rc, &stridx, &matchlen);
	mu_assert_lf(rc == expect_rc);
	mu_assert_lf(stridx == expect_stridx);
	mu_assert_lf(matchlen == expect_matchlen);
	return 0;
}

// ----------------------------------------------------------------
static char* test_disjoint() {
	char* strings[] = { "abc" , "fg" };
	char* buf = "abcde";
	int expect_rc = TRUE, expect_stridx = 0, expect_matchlen = 3;

	int num_strings = sizeof(strings) / sizeof(strings[0]);
	int stridx, matchlen, rc;
	test_case("disjoint", strings, num_strings, buf, &rc, &stridx, &matchlen);
	mu_assert_lf(rc == expect_rc);
	mu_assert_lf(stridx == expect_stridx);
	mu_assert_lf(matchlen == expect_matchlen);
	return 0;
}

// ----------------------------------------------------------------
static char* test_short_long() {
	char* strings[] = { "a" , "aa" };
	char* buf = "aaabc";
	int expect_rc = TRUE, expect_stridx = 1, expect_matchlen = 2;

	int num_strings = sizeof(strings) / sizeof(strings[0]);
	int stridx, matchlen, rc;
	test_case("short_long", strings, num_strings, buf, &rc, &stridx, &matchlen);
	mu_assert_lf(rc == expect_rc);
	mu_assert_lf(stridx == expect_stridx);
	mu_assert_lf(matchlen == expect_matchlen);
	return 0;
}

// ----------------------------------------------------------------
static char* test_long_short() {
	char* strings[] = { "aa" , "a" };
	char* buf = "aaabc";
	int expect_rc = TRUE, expect_stridx = 0, expect_matchlen = 2;

	int num_strings = sizeof(strings) / sizeof(strings[0]);
	int stridx, matchlen, rc;
	test_case("long_short", strings, num_strings, buf, &rc, &stridx, &matchlen);
	mu_assert_lf(rc == expect_rc);
	mu_assert_lf(stridx == expect_stridx);
	mu_assert_lf(matchlen == expect_matchlen);
	return 0;
}

// ----------------------------------------------------------------
static char* test_dkvp() {
	char* test_name = "dkvp";
	char* strings[] = { "=" , ",", "\r\n", "\xff" };
	const int PS_TOKEN  = 0;
	const int FS_TOKEN  = 1;
	const int RS_TOKEN  = 2;
	const int EOF_TOKEN = 3;
	int num_strings = sizeof(strings) / sizeof(strings[0]);
	char* buf =
		"abc=123,def=456\r\n"
		"ghi=789\xff";
	char* p = buf;

	printf("%s %s\n", sep, test_name);
	int stridx, matchlen, rc;

	parse_trie_t* ptrie = parse_trie_alloc();
	parse_trie_print(ptrie);
	for (stridx = 0; stridx < num_strings; stridx++) {
		printf("Adding string[%d] = [", stridx);
		print_with_unprintables(strings[stridx]);
		printf("]\n");
		parse_trie_add_string(ptrie, strings[stridx], stridx);
	}
	parse_trie_print(ptrie);

	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen); mu_assert_lf(rc == FALSE); p++;
	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen); mu_assert_lf(rc == FALSE); p++;
	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen); mu_assert_lf(rc == FALSE); p++;

	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen);
	mu_assert_lf(rc == TRUE);
	mu_assert_lf(stridx == PS_TOKEN);
	mu_assert_lf(matchlen == strlen(strings[PS_TOKEN]));
	p += matchlen;

	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen); mu_assert_lf(rc == FALSE); p++;
	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen); mu_assert_lf(rc == FALSE); p++;
	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen); mu_assert_lf(rc == FALSE); p++;

	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen);
	mu_assert_lf(rc == TRUE);
	mu_assert_lf(stridx == FS_TOKEN);
	mu_assert_lf(matchlen == strlen(strings[FS_TOKEN]));
	p += matchlen;

	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen); mu_assert_lf(rc == FALSE); p++;
	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen); mu_assert_lf(rc == FALSE); p++;
	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen); mu_assert_lf(rc == FALSE); p++;

	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen);
	mu_assert_lf(rc == TRUE);
	mu_assert_lf(stridx == PS_TOKEN);
	mu_assert_lf(matchlen == strlen(strings[PS_TOKEN]));
	p += matchlen;

	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen); mu_assert_lf(rc == FALSE); p++;
	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen); mu_assert_lf(rc == FALSE); p++;
	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen); mu_assert_lf(rc == FALSE); p++;

	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen);
	mu_assert_lf(rc == TRUE);
	mu_assert_lf(stridx == RS_TOKEN);
	mu_assert_lf(matchlen == strlen(strings[RS_TOKEN]));
	p += matchlen;

	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen); mu_assert_lf(rc == FALSE); p++;
	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen); mu_assert_lf(rc == FALSE); p++;
	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen); mu_assert_lf(rc == FALSE); p++;

	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen);
	mu_assert_lf(rc == TRUE);
	mu_assert_lf(stridx == PS_TOKEN);
	mu_assert_lf(matchlen == strlen(strings[PS_TOKEN]));
	p += matchlen;

	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen); mu_assert_lf(rc == FALSE); p++;
	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen); mu_assert_lf(rc == FALSE); p++;
	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen); mu_assert_lf(rc == FALSE); p++;

	rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen);
	mu_assert_lf(rc == TRUE);
	mu_assert_lf(stridx == EOF_TOKEN);
	mu_assert_lf(matchlen == strlen(strings[EOF_TOKEN]));
	p += matchlen;

	return 0;
}

// ----------------------------------------------------------------
static char* show_it() {
	char* test_name = "show_it";
	char* strings[] = { "=" , ",", "\r\n", "\xff" };
	const int EOF_TOKEN = 3;
	int num_strings = sizeof(strings) / sizeof(strings[0]);
	char* buf =
		"abc=123,def=456\r\n"
		"ghi=789\xff";
	char* p = buf;

	printf("%s %s\n", sep, test_name);
	int stridx, matchlen, rc;

	parse_trie_t* ptrie = parse_trie_alloc();
	parse_trie_print(ptrie);
	for (stridx = 0; stridx < num_strings; stridx++) {
		printf("Adding string[%d] = [", stridx);
		print_with_unprintables(strings[stridx]);
		printf("]\n");
		parse_trie_add_string(ptrie, strings[stridx], stridx);
	}
	parse_trie_print(ptrie);

	while (TRUE) {
		rc = parse_trie_ring_match(ptrie, p, 0, strlen(p), 0xff, &stridx, &matchlen);
		if (rc) {
			printf("match token %d (", stridx);
			print_with_unprintables(strings[stridx]);
			printf(")\n");

			p += matchlen;
			if (stridx == EOF_TOKEN) {
				break;
			}
		} else {
			char c = *p;
			printf("c %c[%02x]\n", isprint((unsigned char)c) ? c : '?', ((unsigned)c)&0xff);
			p++;
		}
	}

	mu_assert_lf(*p == 0);

	return 0;
}

// ================================================================
static char* all_tests() {
	mu_run_test(test_simplest);
	mu_run_test(test_disjoint);
	mu_run_test(test_short_long);
	mu_run_test(test_long_short);
	mu_run_test(test_dkvp);
	mu_run_test(show_it);
	return 0;
}

int main(int argc, char** argv) {
	mlr_global_init(argv[0], NULL);
	printf("TEST_PARSE_TRIE ENTER\n");
	char* result = all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_PARSE_TRIE: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
