#ifndef MAPPER_H
#define MAPPER_H

#include "containers/lrec.h"
#include "containers/sllv.h"
#include "mapping/context.h"

// xxx cmt conventions:
// * mem-mgt: pass same lrec or allow new & free old
// * null-terminated streams for then-chaining
// * drain state at EOS

// ----------------------------------------------------------------
// Data plane:

// Returns linked list of records (lrec_t*).
typedef sllv_t* mapper_process_func_t(lrec_t* pinrec, context_t* pctx, void* pvstate);
typedef void    mapper_free_func_t(void* pvstate);

typedef struct _mapper_t {
	void* pvstate;
	mapper_process_func_t* pmapper_process_func;
	mapper_free_func_t*    pmapper_free_func;
} mapper_t;

// ----------------------------------------------------------------
// Control plane:

typedef void mapper_usage_func_t(char* argv0, char* verb);
typedef      mapper_t* mapper_parse_cli_func_t(int* pargi, int argc, char** argv);

typedef struct _mapper_setup_t {
	char*                    verb;
	mapper_usage_func_t*     pusage_func;
	mapper_parse_cli_func_t* pparse_func;
} mapper_setup_t;

#endif // MAPPER_H
