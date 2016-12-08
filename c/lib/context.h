#ifndef CONTEXT_H
#define CONTEXT_H

// File-level context for Miller's NR, FNR, FILENAME, and FILENUM variables, as
// well as for error messages
typedef struct _context_t {
	long long nr;
	long long fnr;
	int       filenum;
	char*     filename;
	int       force_eof; // e.g. mlr head

	char*     ips;
	char*     ifs;
	char*     irs;
	char*     ops;
	char*     ofs;
	char*     ors;
} context_t;

void context_init(context_t* pctx, char* first_file_name);
void context_print(context_t* pctx, char* indent);

#endif // CONTEXT_H
