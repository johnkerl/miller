#include "lemon_msort.h"

/*
** A generic merge-sort program.
**
** USAGE:
** Let "ptr" be a pointer to some structure which is at the head of
** a null-terminated list.  Then to sort the list call:
**
**     ptr = msort(ptr,&(ptr->next),cmpfnc);
**
** In the above, "cmpfnc" is a pointer to a function which compares
** two instances of the structure and returns an integer, as in
** strcmp.  The second argument is a pointer to the pointer to the
** second element of the linked list.  This address is used to compute
** the offset to the "next" field within the structure.  The offset to
** the "next" field must be constant for all structures in the list.
**
** The function returns a new pointer which is the head of the list
** after sorting.
**
** ALGORITHM:
** Merge-sort.
*/

/*
** Return a pointer to the next structure in the linked list.
*/
#define NEXT(A) (*(char**)(((unsigned long)A)+offset))

/*
** Inputs:
**   a:       A sorted, null-terminated linked list.  (May be null).
**   b:       A sorted, null-terminated linked list.  (May be null).
**   cmp:     A pointer to the comparison function.
**   offset:  Offset in the structure to the "next" field.
**
** Return Value:
**   A pointer to the head of a sorted list containing the elements
**   of both a and b.
**
** Side effects:
**   The "next" pointers for elements in the lists a and b are
**   changed.
*/
static char *merge(
	char *a,
	char *b,
	int (*cmp)(),
	int offset)
{
	char *ptr, *head;

	if (a==0) {
		head = b;
	} else if (b==0) {
		head = a;
	} else {
		if ((*cmp)(a,b)<0) {
			ptr = a;
			a = NEXT(a);
		} else {
			ptr = b;
			b = NEXT(b);
		}
		head = ptr;
		while (a && b) {
			if ((*cmp)(a,b)<0) {
				NEXT(ptr) = a;
				ptr = a;
				a = NEXT(a);
			} else {
				NEXT(ptr) = b;
				ptr = b;
				b = NEXT(b);
			}
		}
		if (a)  NEXT(ptr) = a;
		else    NEXT(ptr) = b;
	}
	return head;
}

/*
** Inputs:
**   list:      Pointer to a singly-linked list of structures.
**   next:      Pointer to pointer to the second element of the list.
**   cmp:       A comparison function.
**
** Return Value:
**   A pointer to the head of a sorted list containing the elements
**   orginally in list.
**
** Side effects:
**   The "next" pointers for elements in list are changed.
*/
#define LISTSIZE 30
char *msort(char *list, char **next, int (*cmp)())
{
	unsigned long offset;
	char *ep;
	char *set[LISTSIZE];
	int i;
	offset = (unsigned long)next - (unsigned long)list;
	for(i=0; i<LISTSIZE; i++) set[i] = 0;
	while (list) {
		ep = list;
		list = NEXT(list);
		NEXT(ep) = 0;
		for(i=0; i<LISTSIZE-1 && set[i]!=0; i++){
			ep = merge(ep,set[i],cmp,offset);
			set[i] = 0;
		}
		set[i] = ep;
	}
	ep = 0;
	for(i=0; i<LISTSIZE; i++) if (set[i])  ep = merge(ep,set[i],cmp,offset);
	return ep;
}
