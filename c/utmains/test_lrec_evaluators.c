#include "mapping/lrec_evaluators.h"

// test_lrec_evaluators has the MinUnit inside lrec_evaluators, as it tests
// many private methods. (The other option is to make them all public.)
int main(int argc, char **argv) {
	return test_lrec_evaluators_main(argc, argv);
}
