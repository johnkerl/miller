#ifndef LEMON_UTIL_H
#define LEMON_UTIL_H

int strhash(char *x);

char *Strsafe();
void Strsafe_init();
int Strsafe_insert(char *);
char *Strsafe_find(char *);

#endif // LEMON_UTIL_H
