#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "lemon_plink.h"

/*************************** From the file "plink.c" *********************/
/*
** Routines processing configuration follow-set propagation links
** in the LEMON parser generator.
*/
static struct plink *plink_freelist = 0;

/* Allocate a new plink */
struct plink *Plink_new(){
	struct plink *new;

	if (plink_freelist==0) {
		int i;
		int amt = 100;
		plink_freelist = (struct plink *)malloc (sizeof(struct plink)*amt) ;
		if (plink_freelist==0) {
			fprintf(stderr,
			"Unable to allocate memory for a new follow-set propagation link.\n");
			exit(1);
		}
		for(i=0; i<amt-1; i++) plink_freelist[i].next = &plink_freelist[i+1];
		plink_freelist[amt-1].next = 0;
	}
	new = plink_freelist;
	plink_freelist = plink_freelist->next;
	return new;
}

/* Add a plink to a plink list */
void Plink_add(struct plink **plpp, struct config *cfp)
{
	struct plink *new;
	new = Plink_new();
	new->next = *plpp;
	*plpp = new;
	new->cfp = cfp;
}

/* Transfer every plink on the list "from" to the list "to" */
void Plink_copy(struct plink **to, struct plink *from)
{
	struct plink *nextpl;
	while (from) {
		nextpl = from->next;
		from->next = *to;
		*to = from;
		from = nextpl;
	}
}

/* Delete every plink on the list */
void Plink_delete(struct plink *plp)
{
	struct plink *nextpl;

	while (plp) {
		nextpl = plp->next;
		plp->next = plink_freelist;
		plink_freelist = plp;
		plp = nextpl;
	}
}
