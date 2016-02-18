#ifndef DVECTOR_H
#define DVECTOR_H

typedef struct _dvector_t {
	double* data;
	unsigned long long size;
	unsigned long long capacity;
} dvector_t;

dvector_t* dvector_alloc(unsigned long long initial_capacity);
void dvector_free(dvector_t* pdvector);
void dvector_append(dvector_t* pdvector, double value);

#endif // DVECTOR_H
