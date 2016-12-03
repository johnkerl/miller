#ifndef LEMON_STRING_H
#define LEMON_STRING_H

int strhash(char *x);

char *Strsafe(char*);
void Strsafe_init();
int Strsafe_insert(char *);
char *Strsafe_find(char *);

#endif // LEMON_STRING_H
