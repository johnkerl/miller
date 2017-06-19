#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "lemon_memory.h"
#include "lemon_string.h"

int strhash(char *x) {
	unsigned h = 0;
	while (*x) h = h*13 + *(x++);
	return (int)h;
}

// ================================================================
/* There is one instance of the following structure for each
** associative array of type "x1".
*/
struct s_x1 {
	int size;               /* The number of available slots. Must be a power of 2, >= 1. */
	int count;              /* Number of currently slots filled */
	struct s_x1node *tbl;  /* The data stored here */
	struct s_x1node **ht;  /* Hash table for lookups */
};

/* There is one instance of this structure for every data element
** in an associative array of type "x1".
*/
typedef struct s_x1node {
	char *data;                  /* The data */
	struct s_x1node *next;   /* Next entry with the same hash */
	struct s_x1node **from;  /* Previous link */
} x1node;

/* There is only one instance of the array, which is the following */
static struct s_x1 *x1a;

// ================================================================
/* Works like strdup, sort of.  Save a string in malloced memory, but
** keep strings in a table so that the same string is not in more
** than one place.
*/
char *Strsafe(char *y) {
	char *z;

	z = Strsafe_find(y);
	if (z==0 && (z=malloc (strlen(y)+1) )!=0) {
		strcpy(z,y);
		Strsafe_insert(z);
	}
	MemoryCheck(z);
	return z;
}

/* Allocate a new associative array */
void Strsafe_init(){
	if (x1a)  return;
	x1a = (struct s_x1*)malloc (sizeof(struct s_x1)) ;
	if (x1a) {
		x1a->size = 1024;
		x1a->count = 0;
		x1a->tbl = (x1node*)malloc (
			(sizeof(x1node) + sizeof(x1node*))*1024) ;
		if (x1a->tbl==0) {
			free(x1a);
			x1a = 0;
		} else {
			int i;
			x1a->ht = (x1node**)&(x1a->tbl[1024]);
			for(i=0; i<1024; i++) x1a->ht[i] = 0;
		}
	}
}

/* Insert a new record into the array.  Return TRUE if successful.
** Prior data with the same key is NOT overwritten */
int Strsafe_insert(char *data)
{
	x1node *np;
	int h;
	int ph;

	if (x1a==0)  return 0;
	ph = strhash(data);
	h = ph & (x1a->size-1);
	np = x1a->ht[h];
	while (np) {
		if (strcmp(np->data,data)==0) {
			/* An existing entry with the same key is found. */
			/* Fail because overwrite is not allows. */
			return 0;
		}
		np = np->next;
	}
	if (x1a->count>=x1a->size) {
		/* Need to make the hash table bigger */
		int i,size;
		struct s_x1 array;
		array.size = size = x1a->size*2;
		array.count = x1a->count;
		array.tbl = (x1node*)malloc(
			(sizeof(x1node) + sizeof(x1node*))*size) ;
		if (array.tbl==0)  return 0;  /* Fail due to malloc failure */
		array.ht = (x1node**)&(array.tbl[size]);
		for(i=0; i<size; i++) array.ht[i] = 0;
		for(i=0; i<x1a->count; i++){
			x1node *oldnp, *newnp;
			oldnp = &(x1a->tbl[i]);
			h = strhash(oldnp->data) & (size-1);
			newnp = &(array.tbl[i]);
			if (array.ht[h])  array.ht[h]->from = &(newnp->next);
			newnp->next = array.ht[h];
			newnp->data = oldnp->data;
			newnp->from = &(array.ht[h]);
			array.ht[h] = newnp;
		}
		free(x1a->tbl);
		*x1a = array;
	}
	/* Insert the new data */
	h = ph & (x1a->size-1);
	np = &(x1a->tbl[x1a->count++]);
	np->data = data;
	if (x1a->ht[h])  x1a->ht[h]->from = &(np->next);
	np->next = x1a->ht[h];
	x1a->ht[h] = np;
	np->from = &(x1a->ht[h]);
	return 1;
}

/* Return a pointer to data assigned to the given key.  Return NULL
** if no such key. */
char *Strsafe_find(char *key)
{
	int h;
	x1node *np;

	if (x1a==0)  return 0;
	h = strhash(key) & (x1a->size-1);
	np = x1a->ht[h];
	while (np) {
		if (strcmp(np->data,key)==0)  break;
		np = np->next;
	}
	return np ? np->data : 0;
}
