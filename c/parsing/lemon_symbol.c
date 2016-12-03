#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include "lemon_symbol.h"
#include "lemon_memory.h"
#include "lemon_string.h"

// ================================================================
/* There is one instance of the following structure for each
** associative array of type "x2".
*/
struct s_x2 {
	int size;               /* The number of available slots. Must be a power of 2, >= 1. */
	int count;              /* Number of currently slots filled */
	struct s_x2node *tbl;  /* The data stored here */
	struct s_x2node **ht;  /* Hash table for lookups */
};

/* There is one instance of this structure for every data element
** in an associative array of type "x2".
*/
typedef struct s_x2node {
	struct symbol *data;     /* The data */
	char *key;               /* The key */
	struct s_x2node *next;   /* Next entry with the same hash */
	struct s_x2node **from;  /* Previous link */
} x2node;

/* There is only one instance of the array, which is the following */
static struct s_x2 *x2a;

// ================================================================
/* Return a pointer to the (terminal or nonterminal) symbol "x".
** Create a new symbol if this is the first time "x" has been seen.
*/
struct symbol *Symbol_new(char *x)
{
	struct symbol *sp;

	sp = Symbol_find(x);
	if (sp==0) {
		sp = (struct symbol *)malloc (sizeof(struct symbol)) ;
		MemoryCheck(sp);
		sp->name = Strsafe(x);
		sp->type = isupper(*x) ? TERMINAL : NONTERMINAL;
		sp->rule = 0;
		sp->fallback = 0;
		sp->prec = -1;
		sp->assoc = UNK;
		sp->firstset = 0;
		sp->lambda = B_FALSE;
		sp->destructor = 0;
		sp->datatype = 0;
		Symbol_insert(sp,sp->name);
	}
	return sp;
}

/* Allocate a new associative array */
void Symbol_init(){
	if (x2a)  return;
	x2a = (struct s_x2*)malloc (sizeof(struct s_x2)) ;
	if (x2a) {
		x2a->size = 128;
		x2a->count = 0;
		x2a->tbl = (x2node*)malloc (
			(sizeof(x2node) + sizeof(x2node*))*128) ;
		if (x2a->tbl==0) {
			free(x2a);
			x2a = 0;
		} else {
			int i;
			x2a->ht = (x2node**)&(x2a->tbl[128]);
			for(i=0; i<128; i++) x2a->ht[i] = 0;
		}
	}
}

/* Insert a new record into the array.  Return TRUE if successful.
** Prior data with the same key is NOT overwritten */
int Symbol_insert(struct symbol *data, char *key)
{
	x2node *np;
	int h;
	int ph;

	if (x2a==0)  return 0;
	ph = strhash(key);
	h = ph & (x2a->size-1);
	np = x2a->ht[h];
	while (np) {
		if (strcmp(np->key,key)==0) {
			/* An existing entry with the same key is found. */
			/* Fail because overwrite is not allows. */
			return 0;
		}
		np = np->next;
	}
	if (x2a->count>=x2a->size) {
		/* Need to make the hash table bigger */
		int i,size;
		struct s_x2 array;
		array.size = size = x2a->size*2;
		array.count = x2a->count;
		array.tbl = (x2node*)malloc(
			(sizeof(x2node) + sizeof(x2node*))*size) ;
		if (array.tbl==0)  return 0;  /* Fail due to malloc failure */
		array.ht = (x2node**)&(array.tbl[size]);
		for(i=0; i<size; i++) array.ht[i] = 0;
		for(i=0; i<x2a->count; i++){
			x2node *oldnp, *newnp;
			oldnp = &(x2a->tbl[i]);
			h = strhash(oldnp->key) & (size-1);
			newnp = &(array.tbl[i]);
			if (array.ht[h])  array.ht[h]->from = &(newnp->next);
			newnp->next = array.ht[h];
			newnp->key = oldnp->key;
			newnp->data = oldnp->data;
			newnp->from = &(array.ht[h]);
			array.ht[h] = newnp;
		}
		free(x2a->tbl);
		*x2a = array;
	}
	/* Insert the new data */
	h = ph & (x2a->size-1);
	np = &(x2a->tbl[x2a->count++]);
	np->key = key;
	np->data = data;
	if (x2a->ht[h])  x2a->ht[h]->from = &(np->next);
	np->next = x2a->ht[h];
	x2a->ht[h] = np;
	np->from = &(x2a->ht[h]);
	return 1;
}

/* Return a pointer to data assigned to the given key.  Return NULL
** if no such key. */
struct symbol *Symbol_find(char *key)
{
	int h;
	x2node *np;

	if (x2a==0)  return 0;
	h = strhash(key) & (x2a->size-1);
	np = x2a->ht[h];
	while (np) {
		if (strcmp(np->key,key)==0)  break;
		np = np->next;
	}
	return np ? np->data : 0;
}

/* Return the n-th data.  Return NULL if n is out of range. */
struct symbol *Symbol_Nth(int n)
{
	struct symbol *data;
	if (x2a && n>0 && n<=x2a->count) {
		data = x2a->tbl[n-1].data;
	} else {
		data = 0;
	}
	return data;
}

/* Return the size of the array */
int Symbol_count()
{
	return x2a ? x2a->count : 0;
}

/* Return an array of pointers to all data in the table.
** The array is obtained from malloc.  Return NULL if memory allocation
** problems, or if the array is empty. */
struct symbol **Symbol_arrayof()
{
	struct symbol **array;
	int i,size;
	if (x2a==0)  return 0;
	size = x2a->count;
	array = (struct symbol **)malloc (sizeof(struct symbol *)*size) ;
	if (array) {
		for(i=0; i<size; i++) array[i] = x2a->tbl[i].data;
	}
	return array;
}

/* Compare two symbols for working purposes
**
** Symbols that begin with upper case letters (terminals or tokens)
** must sort before symbols that begin with lower case letters
** (non-terminals).  Other than that, the order does not matter.
**
** We find experimentally that leaving the symbols in their original
** order (the order they appeared in the grammar file) gives the
** smallest parser tables in SQLite.
*/
int Symbolcmpp(struct symbol **a, struct symbol **b) {
	int i1 = (**a).index + 10000000*((**a).name[0]>'Z');
	int i2 = (**b).index + 10000000*((**b).name[0]>'Z');
	return i1-i2;
}
