#ifndef LEMON_STATE_TABLE_H
#define LEMON_STATE_TABLE_H

/* Routines to manage the state table */

#include "lemon_structs.h"

struct state *State_new();
void State_init(void);
int State_insert(struct state *, struct config *);
struct state *State_find(struct config *);
struct state **State_arrayof();

#endif // LEMON_STATE_TABLE_H
