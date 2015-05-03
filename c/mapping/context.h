#ifndef CONTEXT_H
#define CONTEXT_H

typedef struct _static_context_t {
	char*     argv0;
	char*     ofmt;
} static_context_t;

typedef struct _dynamic_context_t {
	long long nr;
	long long fnr;
	int       filenum;
	char*     filename;
} dynamic_context_t;

typedef struct _context_t {
	static_context_t  statx;
	dynamic_context_t dynx;
} context_t;

#endif // CONTEXT_H
