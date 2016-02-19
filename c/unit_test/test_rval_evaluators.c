#include "mapping/rval_evaluators.h"

// test_rval_evaluators has the MinUnit inside rval_evaluators, as it tests
// many private methods. (The other option is to make them all public.)
int main(int argc, char **argv) {
	return test_rval_evaluators_main(argc, argv);
}
