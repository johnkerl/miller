#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "lemon_assert.h"
#include "lemon_action.h"

// xxx move:
char *msort(char *list, char **next, int (*cmp)());


/*
** Routines processing parser actions in the LEMON parser generator.
*/

/* Allocate a new parser action */
struct action *Action_new() {
	static struct action *freelist = 0;
	struct action *new;

	if (freelist==0) {
		int i;
		int amt = 100;
		freelist = (struct action *)malloc (sizeof(struct action)*amt) ;
		if (freelist==0) {
			fprintf(stderr,"Unable to allocate memory for a new parser action.");
			exit(1);
		}
		for(i=0; i<amt-1; i++) freelist[i].next = &freelist[i+1];
		freelist[amt-1].next = 0;
	}
	new = freelist;
	freelist = freelist->next;
	return new;
}

/* Compare two actions */
static int actioncmp(struct action *ap1, struct action *ap2)
{
	int rc;
	rc = ap1->sp->index - ap2->sp->index;
	if (rc==0)  rc = (int)ap1->type - (int)ap2->type;
	if (rc==0) {
		assert (ap1->type==REDUCE || ap1->type==RD_RESOLVED || ap1->type==CONFLICT);
		assert (ap2->type==REDUCE || ap2->type==RD_RESOLVED || ap2->type==CONFLICT);
		rc = ap1->x.rp->index - ap2->x.rp->index;
	}
	return rc;
}

/* Sort parser actions */
struct action *Action_sort(struct action *ap)
{
	ap = (struct action *)msort((char *)ap,(char **)&ap->next,actioncmp);
	return ap;
}

void Action_add(
	struct action **app,
	enum e_action type,
	struct symbol *sp,
	char *arg)
{
	struct action *new;
	new = Action_new();
	new->next = *app;
	*app = new;
	new->type = type;
	new->sp = sp;
	if (type==SHIFT) {
		new->x.stp = (struct state *)arg;
	} else {
		new->x.rp = (struct rule *)arg;
	}
}

/* Free all memory associated with the given acttab */
void acttab_free(acttab *p) {
	free (p->aAction);
	free (p->aLookahead);
	free (p);
}

/* Allocate a new acttab structure */
acttab *acttab_alloc(void) {
	acttab *p = malloc (sizeof(*p)) ;
	if (p==0) {
		fprintf(stderr,"Unable to allocate memory for a new acttab.");
		exit(1);
	}
	memset(p, 0, sizeof(*p));
	return p;
}

/* Add a new action to the current transaction set
*/
void acttab_action(acttab *p, int lookahead, int action) {
	if (p->nLookahead>=p->nLookaheadAlloc) {
		p->nLookaheadAlloc += 25;
		p->aLookahead = realloc (p->aLookahead, sizeof(p->aLookahead[0])*p->nLookaheadAlloc) ;
		if (p->aLookahead==0) {
			fprintf(stderr,"malloc failed\n");
			exit(1);
		}
	}
	if (p->nLookahead==0) {
		p->mxLookahead = lookahead;
		p->mnLookahead = lookahead;
		p->mnAction = action;
	} else {
		if (p->mxLookahead<lookahead)  p->mxLookahead = lookahead;
		if (p->mnLookahead>lookahead) {
			p->mnLookahead = lookahead;
			p->mnAction = action;
		}
	}
	p->aLookahead[p->nLookahead].lookahead = lookahead;
	p->aLookahead[p->nLookahead].action = action;
	p->nLookahead++;
}

/*
** Add the transaction set built up with prior calls to acttab_action()
** into the current action table.  Then reset the transaction set back
** to an empty set in preparation for a new round of acttab_action() calls.
**
** Return the offset into the action table of the new transaction.
*/
int acttab_insert(acttab *p) {
	int i, j, k, n;
	assert (p->nLookahead>0) ;

	/* Make sure we have enough space to hold the expanded action table
	** in the worst case.  The worst case occurs if the transaction set
	** must be appended to the current action table
	*/
	n = p->mxLookahead + 1;
	if (p->nAction + n >= p->nActionAlloc) {
		int oldAlloc = p->nActionAlloc;
		p->nActionAlloc = p->nAction + n + p->nActionAlloc + 20;
		p->aAction = realloc (p->aAction, sizeof(p->aAction[0])*p->nActionAlloc);
		if (p->aAction==0) {
			fprintf(stderr,"malloc failed\n");
			exit(1);
		}
		for(i=oldAlloc; i<p->nActionAlloc; i++){
			p->aAction[i].lookahead = -1;
			p->aAction[i].action = -1;
		}
	}

	/* Scan the existing action table looking for an offset where we can
	** insert the current transaction set.  Fall out of the loop when that
	** offset is found.  In the worst case, we fall out of the loop when
	** i reaches p->nAction, which means we append the new transaction set.
	**
	** i is the index in p->aAction[] where p->mnLookahead is inserted.
	*/
	for(i=0; i<p->nAction+p->mnLookahead; i++){
		if (p->aAction[i].lookahead<0) {
			for(j=0; j<p->nLookahead; j++){
				k = p->aLookahead[j].lookahead - p->mnLookahead + i;
				if (k<0)  break;
				if (p->aAction[k].lookahead>=0)  break;
			}
			if (j<p->nLookahead)  continue;
			for(j=0; j<p->nAction; j++){
				if (p->aAction[j].lookahead==j+p->mnLookahead-i)  break;
			}
			if (j==p->nAction) {
				break;  /* Fits in empty slots */
			}
		} else if (p->aAction[i].lookahead==p->mnLookahead) {
			if (p->aAction[i].action!=p->mnAction)  continue;
			for(j=0; j<p->nLookahead; j++){
				k = p->aLookahead[j].lookahead - p->mnLookahead + i;
				if (k<0 || k>=p->nAction)  break;
				if (p->aLookahead[j].lookahead!=p->aAction[k].lookahead)  break;
				if (p->aLookahead[j].action!=p->aAction[k].action)  break;
			}
			if (j<p->nLookahead)  continue;
			n = 0;
			for(j=0; j<p->nAction; j++){
				if (p->aAction[j].lookahead<0)  continue;
				if (p->aAction[j].lookahead==j+p->mnLookahead-i)  n++;
			}
			if (n==p->nLookahead) {
				break;  /* Same as a prior transaction set */
			}
		}
	}
	/* Insert transaction set at index i. */
	for(j=0; j<p->nLookahead; j++){
		k = p->aLookahead[j].lookahead - p->mnLookahead + i;
		p->aAction[k] = p->aLookahead[j];
		if (k>=p->nAction)  p->nAction = k+1;
	}
	p->nLookahead = 0;

	/* Return the offset that is added to the lookahead in order to get the
	** index into yy_action of the action */
	return i - p->mnLookahead;
}
