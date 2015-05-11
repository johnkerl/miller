#ifndef CONTEXT_H
#define CONTEXT_H

typedef struct _context_t {
	long long nr;
	long long fnr;
	int       filenum;
	char*     filename;
} context_t;

#endif // CONTEXT_H
