#ifndef HEADER_KEEPER_H
#define HEADER_KEEPER_H

#include "containers/slls.h"

// xxx cmt w/ xref

typedef struct _header_keeper_t {
	char*   line;
	slls_t* pkeys;
} header_keeper_t;

header_keeper_t* header_keeper_alloc(char* line, slls_t* pkeys);
void header_keeper_free(header_keeper_t* pheader_keeper);

#endif // HEADER_KEEPER_H
