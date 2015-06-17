#ifndef MLR_GLOBALS_H
#define MLR_GLOBALS_H
#include "cli/mlrcli.h"

typedef struct _mlr_globals_t {
	char*       argv0;
	char*       ofmt;
	// These are shared by mlrcli.c and mlrmain.c. The only reason for their
	// exposure anywhere else is to communicate format and separator flags to
	// mapper_join, which (unlike other mappers) needs to do its own file I/O.
	cli_opts_t* popts;
} mlr_globals_t;
extern mlr_globals_t MLR_GLOBALS;
void mlr_global_init(char* argv0, char* ofmt, cli_opts_t* popts);

#endif // MLR_GLOBALS_H
