#ifndef DVECTOR_H
#define DVECTOR_H

typedef struct _dvector_t {
	double* data;
	int     size;
	int     capacity;
} dvector_t;

dvector_t* dvector_alloc(int initial_capacity);
void dvector_free(dvector_t* pdvector);
void dvector_append(dvector_t* pdvector, double value);

#endif // DVECTOR_H
