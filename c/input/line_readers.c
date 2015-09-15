#include <stdio.h>
#include "lib/mlrutil.h"

// xxx under construction

// Use powers of two exclusively, to help avoid heap fragmentation
#define INITIAL_SIZE 128

// xxx look up what restrict do ... should i be using these more often?
// xxx limited semantics: initial linep & linecapp are disregarded; doesn't return the delimiter in the string.
// xxx getcdelim is just for comparison to getdelim. getsdelim is the deliverable.
size_t mlr_getcdelim(char ** restrict ppline, size_t * restrict plinecap, int delimiter, FILE * restrict fp) {
	size_t linecap = INITIAL_SIZE;
	char* pline = mlr_malloc_or_die(INITIAL_SIZE);
	char* p = pline;
	int len = 0;

	while (TRUE) {
		if (len >= linecap) {
			linecap = linecap << 1;
			// xxx mlr_realloc_or_die
			pline = realloc(pline, linecap);
			p = pline;
		}
		int c = getc_unlocked(fp);
		if (c == EOF) {
			break;
		} else if (c == delimiter) {
			*(p++) = 0;
		} else {
			*(p++) = c;
		}
	}

	*ppline = pline;
	*plinecap = linecap;
	len = p - pline;
	return len;
}

size_t mlr_getsdelim(char ** restrict ppline, size_t * restrict plinecap, char* delimiter, FILE * restrict fp) {
	return 0; // xxx stub
}
