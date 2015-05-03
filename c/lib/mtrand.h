#ifndef MTRAND_H
#define MTRAND_H

void     mtrand_init_default();
void     mtrand_init(unsigned s);
void     mtrand_init_from_array(unsigned init_key[], int key_length);
unsigned get_mtrand_int32(void);
int      get_mtrand_int31(void);
double   get_mtrand_float(void);
double   get_mtrand_double(void);

#endif // MTRAND_H
