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
static char * test_numbers_only() {
	json_char* input = "123";
	json_char* end   = NULL;
	json_value_t* pvalue = json_parse_for_unit_test(input, &end);
	mu_assert_lf(pvalue != NULL);
	mu_assert_lf(pvalue->type == JSON_INTEGER);
	printf("sval: [%s]\n", pvalue->u.integer.sval);
	// xxx make a dump-node method: non-recursive version
	// xxx make an argv-option :)

	return 0;
}

//typedef struct _json_value_t {
//	struct _json_value_t * parent;
//	json_type_t type;
//	union {
//		struct {
//			int nval;
//			unsigned int length;
//			char* sval;
//		} boolean;
//		struct {
//			json_int_t nval;
//			unsigned int length;
//			char* sval;
//		} integer;
//		struct {
//			double nval;
//			unsigned int length;
//			char* sval;
//		} dbl;
//		struct {
//			unsigned int length;
//			json_char * ptr; /* null-terminated */
//		} string;
//		struct {
//			unsigned int length;
//			union {
//				json_object_entry_t * values;
//				char* mem;
//			} p;
//		} object;
//		struct {
//			unsigned int length;
//			struct _json_value_t ** values;
//		} array;
//	} u;
//	union {
//		struct _json_value_t * next_alloc;
//		union {
//			void * pvobject_mem;
//			char * pobject_mem;
//		} p;
//	} _reserved;
//	unsigned int line, col;
//} json_value_t;

// ================================================================
static char * all_tests() {
	mu_run_test(test_numbers_only);
	return 0;
}

int main(int argc, char **argv) {
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
