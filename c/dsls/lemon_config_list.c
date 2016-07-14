#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "lemon_config_list.h"

#include "lemon_assert.h"
#include "lemon_error.h"
#include "lemon_msort.h"
#include "lemon_plink.h"
#include "lemon_set.h"

/*
** Routines to processing a configuration list and building a state
** in the LEMON parser generator.
*/

// ----------------------------------------------------------------
/* Compare two configurations */
int Configcmp(a,b)
struct config *a;
struct config *b;
{
	int x;
	x = a->rp->index - b->rp->index;
	if (x==0)  x = a->dot - b->dot;
	return x;
}

// ----------------------------------------------------------------
static struct config *freelist = 0;      /* List of free configurations */
static struct config *current = 0;       /* Top of list of configurations */
static struct config **currentend = 0;   /* Last on list of configs */
static struct config *basis = 0;         /* Top of list of basis configs */
static struct config **basisend = 0;     /* End of list of basis configs */

/* Return a pointer to a new configuration */
static struct config *newconfig() {
	struct config *new;
	if (freelist==0) {
		int i;
		int amt = 3;
		freelist = (struct config *)malloc (sizeof(struct config)*amt) ;
		if (freelist==0) {
			fprintf(stderr,"Unable to allocate memory for a new configuration.");
			exit(1);
		}
		for(i=0; i<amt-1; i++) freelist[i].next = &freelist[i+1];
		freelist[amt-1].next = 0;
	}
	new = freelist;
	freelist = freelist->next;
	return new;
}

/* The configuration "old" is no longer used */
static void deleteconfig(struct config *old)
{
	old->next = freelist;
	freelist = old;
}

/* Initialized the configuration list builder */
void Configlist_init(){
	current = 0;
	currentend = &current;
	basis = 0;
	basisend = &basis;
	Configtable_init();
	return;
}

/* Initialized the configuration list builder */
void Configlist_reset(){
	current = 0;
	currentend = &current;
	basis = 0;
	basisend = &basis;
	Configtable_clear(0);
	return;
}

/* Add another configuration to the configuration list */
struct config *Configlist_add(
	struct rule *rp,    /* The rule */
	int dot)            /* Index into the RHS of the rule where the dot goes */
{
	struct config *cfp, model;

	assert (currentend!=0) ;
	model.rp = rp;
	model.dot = dot;
	cfp = Configtable_find(&model);
	if (cfp==0) {
		cfp = newconfig();
		cfp->rp = rp;
		cfp->dot = dot;
		cfp->fws = SetNew();
		cfp->stp = 0;
		cfp->fplp = cfp->bplp = 0;
		cfp->next = 0;
		cfp->bp = 0;
		*currentend = cfp;
		currentend = &cfp->next;
		Configtable_insert(cfp);
	}
	return cfp;
}

/* Add a basis configuration to the configuration list */
struct config *Configlist_addbasis(struct rule *rp, int dot)
{
	struct config *cfp, model;

	assert (basisend!=0);
	assert (currentend!=0);
	model.rp = rp;
	model.dot = dot;
	cfp = Configtable_find(&model);
	if (cfp == 0) {
		cfp = newconfig();
		cfp->rp = rp;
		cfp->dot = dot;
		cfp->fws = SetNew();
		cfp->stp = 0;
		cfp->fplp = cfp->bplp = 0;
		cfp->next = 0;
		cfp->bp = 0;
		*currentend = cfp;
		currentend = &cfp->next;
		*basisend = cfp;
		basisend = &cfp->bp;
		Configtable_insert(cfp);
	}
	return cfp;
}

/* Compute the closure of the configuration list */
void Configlist_closure(struct lemon *lemp)
{
	struct config *cfp, *newcfp;
	struct rule *rp, *newrp;
	struct symbol *sp, *xsp;
	int i, dot;

	assert (currentend!=0) ;
	for(cfp=current; cfp; cfp=cfp->next){
		rp = cfp->rp;
		dot = cfp->dot;
		if (dot>=rp->nrhs)  continue;
		sp = rp->rhs[dot];
		if (sp->type==NONTERMINAL) {
			if (sp->rule==0 && sp!=lemp->errsym) {
				ErrorMsg(lemp->filename,rp->line,"Nonterminal \"%s\" has no rules.",
					sp->name);
				lemp->errorcnt++;
			}
			for(newrp=sp->rule; newrp; newrp=newrp->nextlhs){
				newcfp = Configlist_add(newrp,0);
				for (i=dot+1; i<rp->nrhs; i++) {
					xsp = rp->rhs[i];
					if (xsp->type==TERMINAL) {
						SetAdd(newcfp->fws,xsp->index);
						break;
					} else {
						SetUnion(newcfp->fws,xsp->firstset);
						if (xsp->lambda==B_FALSE)  break;
					}
				}
				if (i==rp->nrhs)  Plink_add(&cfp->fplp,newcfp);
			}
		}
	}
	return;
}

/* Sort the configuration list */
void Configlist_sort() {
	current = (struct config *)msort((char *)current,(char **)&(current->next),Configcmp);
	currentend = 0;
	return;
}

/* Sort the basis configuration list */
void Configlist_sortbasis() {
	basis = (struct config *)msort((char *)current,(char **)&(current->bp),Configcmp);
	basisend = 0;
	return;
}

/* Return a pointer to the head of the configuration list and
** reset the list */
struct config *Configlist_return() {
	struct config *old;
	old = current;
	current = 0;
	currentend = 0;
	return old;
}

/* Return a pointer to the head of the configuration list and
** reset the list */
struct config *Configlist_basis() {
	struct config *old;
	old = basis;
	basis = 0;
	basisend = 0;
	return old;
}

/* Free all elements of the given configuration list */
void Configlist_eat(struct config *cfp)
{
	struct config *nextcfp;
	for(; cfp; cfp=nextcfp){
		nextcfp = cfp->next;
		assert (cfp->fplp==0) ;
		assert (cfp->bplp==0) ;
		if (cfp->fws)  SetFree(cfp->fws);
		deleteconfig(cfp);
	}
	return;
}

// ================================================================
/* Hash a configuration */
static int confighash(a)
struct config *a;
{
	int h=0;
	h = h*571 + a->rp->index*37 + a->dot;
	return h;
}

/* There is one instance of the following structure for each
** associative array of type "x4".
*/
struct s_x4 {
	int size;               /* The number of available slots.  Must be a power of 2, >= 1. */
	int count;              /* Number of currently slots filled */
	struct s_x4node *tbl;  /* The data stored here */
	struct s_x4node **ht;  /* Hash table for lookups */
};

/* There is one instance of this structure for every data element
** in an associative array of type "x4".
*/
typedef struct s_x4node {
	struct config *data;     /* The data */
	struct s_x4node *next;   /* Next entry with the same hash */
	struct s_x4node **from;  /* Previous link */
} x4node;

/* There is only one instance of the array, which is the following */
static struct s_x4 *x4a;

/* Allocate a new associative array */
void Configtable_init() {
	if (x4a)  return;
	x4a = (struct s_x4*)malloc (sizeof(struct s_x4)) ;
	if (x4a) {
		x4a->size = 64;
		x4a->count = 0;
		x4a->tbl = (x4node*)malloc (
			(sizeof(x4node) + sizeof(x4node*))*64) ;
		if (x4a->tbl==0) {
			free(x4a);
			x4a = 0;
		} else {
			int i;
			x4a->ht = (x4node**)&(x4a->tbl[64]);
			for(i=0; i<64; i++) x4a->ht[i] = 0;
		}
	}
}

/* Insert a new record into the array.  Return TRUE if successful.
** Prior data with the same key is NOT overwritten */
int Configtable_insert(struct config *data) {
	x4node *np;
	int h;
	int ph;

	if (x4a==0)  return 0;
	ph = confighash(data);
	h = ph & (x4a->size-1);
	np = x4a->ht[h];
	while (np) {
		if (Configcmp(np->data,data)==0) {
			/* An existing entry with the same key is found. */
			/* Fail because overwrite is not allows. */
			return 0;
		}
		np = np->next;
	}
	if (x4a->count>=x4a->size) {
		/* Need to make the hash table bigger */
		int i,size;
		struct s_x4 array;
		array.size = size = x4a->size*2;
		array.count = x4a->count;
		array.tbl = (x4node*)malloc(
			(sizeof(x4node) + sizeof(x4node*))*size) ;
		if (array.tbl==0)  return 0;  /* Fail due to malloc failure */
		array.ht = (x4node**)&(array.tbl[size]);
		for(i=0; i<size; i++) array.ht[i] = 0;
		for(i=0; i<x4a->count; i++){
			x4node *oldnp, *newnp;
			oldnp = &(x4a->tbl[i]);
			h = confighash(oldnp->data) & (size-1);
			newnp = &(array.tbl[i]);
			if (array.ht[h])  array.ht[h]->from = &(newnp->next);
			newnp->next = array.ht[h];
			newnp->data = oldnp->data;
			newnp->from = &(array.ht[h]);
			array.ht[h] = newnp;
		}
		free(x4a->tbl);
		*x4a = array;
	}
	/* Insert the new data */
	h = ph & (x4a->size-1);
	np = &(x4a->tbl[x4a->count++]);
	np->data = data;
	if (x4a->ht[h])  x4a->ht[h]->from = &(np->next);
	np->next = x4a->ht[h];
	x4a->ht[h] = np;
	np->from = &(x4a->ht[h]);
	return 1;
}

/* Return a pointer to data assigned to the given key.  Return NULL
** if no such key. */
struct config *Configtable_find(struct config *key)
{
	int h;
	x4node *np;

	if (x4a==0)  return 0;
	h = confighash(key) & (x4a->size-1);
	np = x4a->ht[h];
	while (np) {
		if (Configcmp(np->data,key)==0)  break;
		np = np->next;
	}
	return np ? np->data : 0;
}

/* Remove all data from the table.  Pass each data to the function "f"
** as it is removed.  ("f" may be null to avoid this step.) */
void Configtable_clear(int(*f)(/* struct config * */))
{
	int i;
	if (x4a==0 || x4a->count==0)  return;
	if (f)  for(i=0; i<x4a->count; i++) (*f)(x4a->tbl[i].data);
	for(i=0; i<x4a->size; i++) x4a->ht[i] = 0;
	x4a->count = 0;
	return;
}
