#ifndef MLR_GLOBALS_H
#define MLR_GLOBALS_H

typedef struct _mlr_globals_t {
	char* bargv0; // basename of argv0
	char* ofmt;
} mlr_globals_t;
extern mlr_globals_t MLR_GLOBALS;
void mlr_global_init(char* argv0, char* ofmt);

#endif // MLR_GLOBALS_H
