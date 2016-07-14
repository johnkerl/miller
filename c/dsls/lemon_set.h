#ifndef LEMON_SET_H
#define LEMON_SET_H

void  SetSize(int N);             /* All sets will be of size N */
char *SetNew(void);               /* A new set for element 0..N */
void  SetFree(char*);             /* Deallocate a set */

int SetAdd(char*,int);            /* Add element to a set */
int SetUnion(char *A,char *B);    /* A <- A U B, thru element N */

#define SetFind(X,Y) (X[Y])       /* True if Y is in set X */

#endif // LEMON_SET_H
