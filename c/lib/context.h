#ifndef CONTEXT_H
#define CONTEXT_H

// xxx cmt
typedef struct _context_t {
	long long nr;
	long long fnr;
	int       filenum;
	char*     filename;
} context_t;

void context_init(context_t* pctx, char* first_file_name);

void context_print(context_t* pctx, char* indent);

#endif // CONTEXT_H
