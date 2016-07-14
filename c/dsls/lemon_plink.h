#ifndef LEMON_PLINK_H
#define LEMON_PLINK_H

#include "lemon_structs.h"

struct plink *Plink_new(void);
void Plink_add(struct plink **, struct config *);
void Plink_copy(struct plink **, struct plink *);
void Plink_delete(struct plink *);

#endif // LEMON_PLINK_H
