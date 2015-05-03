#ifndef HDR_KEEPER_H
#define HDR_KEEPER_H

#include "containers/slls.h"

// xxx cmt

typedef struct _hdr_keeper_t {
	char*   line;
	slls_t* pkeys;
} hdr_keeper_t;

hdr_keeper_t* hdr_keeper_alloc(char* line, slls_t* pkeys);
void hdr_keeper_free(hdr_keeper_t* phdr_keeper);

#endif // HDR_KEEPER_H
