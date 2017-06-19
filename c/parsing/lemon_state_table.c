#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "lemon_state_table.h"

#include "lemon_memory.h"

/*
** Code for processing tables in the LEMON parser generator.
*/

/* Compare two states */
static int statecmp(a,b)
struct config *a;
struct config *b;
{
	int rc;
	for(rc=0; rc==0 && a && b;  a=a->bp, b=b->bp){
		rc = a->rp->index - b->rp->index;
		if (rc==0)  rc = a->dot - b->dot;
	}
	if (rc==0) {
		if (a)  rc = 1;
		if (b)  rc = -1;
	}
	return rc;
}

/* Hash a state */
static int statehash(a)
struct config *a;
{
	unsigned h=0;
	while (a) {
		h = h*571 + (unsigned)a->rp->index*37 + (unsigned)a->dot;
		a = a->bp;
	}
	return (int)h;
}

/* Allocate a new state structure */
struct state *State_new()
{
	struct state *new;
	new = (struct state *)malloc (sizeof(struct state)) ;
	MemoryCheck(new);
	return new;
}

/* There is one instance of the following structure for each
** associative array of type "x3".
*/
struct s_x3 {
	int size;               /* The number of available slots. Must be a power of 2, >= 1. */
	int count;              /* Number of currently slots filled */
	struct s_x3node *tbl;  /* The data stored here */
	struct s_x3node **ht;  /* Hash table for lookups */
};

/* There is one instance of this structure for every data element
** in an associative array of type "x3".
*/
typedef struct s_x3node {
	struct state *data;                  /* The data */
	struct config *key;                   /* The key */
	struct s_x3node *next;   /* Next entry with the same hash */
	struct s_x3node **from;  /* Previous link */
} x3node;

/* There is only one instance of the array, which is the following */
static struct s_x3 *x3a;

/* Allocate a new associative array */
void State_init(){
	if (x3a)  return;
	x3a = (struct s_x3*)malloc (sizeof(struct s_x3)) ;
	if (x3a) {
		x3a->size = 128;
		x3a->count = 0;
		x3a->tbl = (x3node*)malloc (
			(sizeof(x3node) + sizeof(x3node*))*128) ;
		if (x3a->tbl==0) {
			free(x3a);
			x3a = 0;
		} else {
			int i;
			x3a->ht = (x3node**)&(x3a->tbl[128]);
			for(i=0; i<128; i++) x3a->ht[i] = 0;
		}
	}
}

/* Insert a new record into the array.  Return TRUE if successful.
** Prior data with the same key is NOT overwritten */
int State_insert(struct state *data, struct config *key)
{
	x3node *np;
	int h;
	int ph;

	if (x3a==0)  return 0;
	ph = statehash(key);
	h = ph & (x3a->size-1);
	np = x3a->ht[h];
	while (np) {
		if (statecmp(np->key,key)==0) {
			/* An existing entry with the same key is found. */
			/* Fail because overwrite is not allows. */
			return 0;
		}
		np = np->next;
	}
	if (x3a->count>=x3a->size) {
		/* Need to make the hash table bigger */
		int i,size;
		struct s_x3 array;
		array.size = size = x3a->size*2;
		array.count = x3a->count;
		array.tbl = (x3node*)malloc(
			(sizeof(x3node) + sizeof(x3node*))*size) ;
		if (array.tbl==0)  return 0;  /* Fail due to malloc failure */
		array.ht = (x3node**)&(array.tbl[size]);
		for(i=0; i<size; i++) array.ht[i] = 0;
		for(i=0; i<x3a->count; i++){
			x3node *oldnp, *newnp;
			oldnp = &(x3a->tbl[i]);
			h = statehash(oldnp->key) & (size-1);
			newnp = &(array.tbl[i]);
			if (array.ht[h])  array.ht[h]->from = &(newnp->next);
			newnp->next = array.ht[h];
			newnp->key = oldnp->key;
			newnp->data = oldnp->data;
			newnp->from = &(array.ht[h]);
			array.ht[h] = newnp;
		}
		free(x3a->tbl);
		*x3a = array;
	}
	/* Insert the new data */
	h = ph & (x3a->size-1);
	np = &(x3a->tbl[x3a->count++]);
	np->key = key;
	np->data = data;
	if (x3a->ht[h])  x3a->ht[h]->from = &(np->next);
	np->next = x3a->ht[h];
	x3a->ht[h] = np;
	np->from = &(x3a->ht[h]);
	return 1;
}

/* Return a pointer to data assigned to the given key.  Return NULL
** if no such key. */
struct state *State_find(struct config *key)
{
	int h;
	x3node *np;

	if (x3a==0)  return 0;
	h = statehash(key) & (x3a->size-1);
	np = x3a->ht[h];
	while (np) {
		if (statecmp(np->key,key)==0)  break;
		np = np->next;
	}
	return np ? np->data : 0;
}

/* Return an array of pointers to all data in the table.
** The array is obtained from malloc.  Return NULL if memory allocation
** problems, or if the array is empty. */
struct state **State_arrayof()
{
	struct state **array;
	int i,size;
	if (x3a==0)  return 0;
	size = x3a->count;
	array = (struct state **)malloc (sizeof(struct state *)*size) ;
	if (array) {
		for(i=0; i<size; i++) array[i] = x3a->tbl[i].data;
	}
	return array;
}
