#ifndef RETURN_STATE_H
#define RETURN_STATE_H

#include "containers/mlhmmv.h"

typedef struct _return_state_t {
	boxed_xval_t retval;
	int returned;
} return_state_t;

#endif // RETURN_STATE_H
