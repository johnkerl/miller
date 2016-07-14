#ifndef LEMON_SYMBOL_H
#define LEMON_SYMBOL_H

#include "lemon_structs.h"

struct symbol *Symbol_new(char *x);
void Symbol_init();
int  Symbol_insert(struct symbol *data, char *key);
struct symbol *Symbol_find(char *key);
struct symbol *Symbol_Nth(int n);
int Symbol_count();
struct symbol **Symbol_arrayof();
int Symbolcmpp(struct symbol **a, struct symbol **b);

#endif // LEMON_SYMBOL_H
