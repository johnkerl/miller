#ifndef MAPPER_H
#define MAPPER_H

#include <stdio.h>
#include "lib/context.h"
#include "containers/lrec.h"
#include "containers/sllv.h"

// See ../README.md for memory-management conventions.

// ----------------------------------------------------------------
// Data plane:

struct _mapper_t; // forward reference for method declarations

// Returns linked list of records (lrec_t*).
typedef sllv_t* mapper_process_func_t(lrec_t* pinrec, context_t* pctx, void* pvstate);

typedef void mapper_free_func_t(struct _mapper_t* pmapper);

typedef struct _mapper_t {
	void* pvstate;
	mapper_process_func_t* pprocess_func;
	mapper_free_func_t*    pfree_func; // virtual destructor
} mapper_t;

// ----------------------------------------------------------------
// Control plane:

typedef void mapper_usage_func_t(FILE* o, char* argv0, char* verb);
typedef      mapper_t* mapper_parse_cli_func_t(int* pargi, int argc, char** argv);

typedef struct _mapper_setup_t {
	char*                    verb;
	mapper_usage_func_t*     pusage_func;
	mapper_parse_cli_func_t* pparse_func;
} mapper_setup_t;

#endif // MAPPER_H
