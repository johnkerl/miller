// http://www.jera.com/techinfo/jtns/jtn002.html
// MinUnit license:
// "You may use the code in this tech note for any purpose, with the understanding that it comes with NO WARRANTY."

#ifndef MINUNIT_H
#define MINUNIT_H

#define MU_STRINGIFY_2(x) #x
#define MU_STRINGIFY_1(x) MU_STRINGIFY_2(x)

#define mu_assert_lf(test) mu_assert(__FILE__ " line " MU_STRINGIFY_1(__LINE__), test)

#define mu_assert(message, test) do { \
	assertions_run++; \
	if (!(test)) { \
		assertions_failed++; \
		return message; \
	} \
} while (0)

#define mu_run_test(test) do { \
	char *message = test(); \
	tests_run++; \
	if (message) { \
		tests_failed++; \
		printf("Failure at %s, invoked from file %s line %d\n", message, __FILE__, __LINE__); \
		return message; \
	} \
} while (0)

extern int tests_run;
extern int tests_failed;
extern int assertions_run;
extern int assertions_failed;

#endif // MINUNIT_H
