#include <stdio.h>
#include <string.h>
#include "lemon_error.h"

/***************** From the file "error.c" *********************************/
/*
** Code for printing error message.
*/

/* Find a good place to break "msg" so that its length is at least "min"
** but no more than "max".  Make the point as close to max as possible.
*/
static int findbreak(char *msg, int min, int max) {
	int i,spot;
	char c;
	for(i=spot=min; i<=max; i++){
		c = msg[i];
		if (c=='\t')  msg[i] = ' ';
		if (c=='\n') { msg[i] = ' '; spot = i; break; }
		if (c==0) { spot = i; break; }
		if (c=='-' && i<max-1)  spot = i+1;
		if (c==' ')  spot = i;
	}
	return spot;
}

/*
** The error message is split across multiple lines if necessary.  The
** splits occur at a space, if there is a space available near the end
** of the line.
*/
#define ERRMSGSIZE  10000 /* Hope this is big enough.  No way to error check */
#define LINEWIDTH      79 /* Max width of any output line */
#define PREFIXLIMIT    30 /* Max width of the prefix on each line */
void ErrorMsg(const char *filename, int lineno, const char *format, ...) {
	char errmsg[ERRMSGSIZE];
	char prefix[PREFIXLIMIT+10];
	int errmsgsize;
	int prefixsize;
	int availablewidth;
	va_list ap;
	int end, restart, base;

	va_start(ap, format);
	/* Prepare a prefix to be prepended to every output line */
	if (lineno>0) {
		sprintf(prefix,"%.*s:%d: ",PREFIXLIMIT-10,filename,lineno);
	} else {
		sprintf(prefix,"%.*s: ",PREFIXLIMIT-10,filename);
	}
	prefixsize = strlen(prefix);
	availablewidth = LINEWIDTH - prefixsize;

	/* Generate the error message */
	vsprintf(errmsg,format,ap);
	va_end(ap);
	errmsgsize = strlen(errmsg);
	/* Remove trailing '\n's from the error message. */
	while (errmsgsize>0 && errmsg[errmsgsize-1]=='\n') {
		 errmsg[--errmsgsize] = 0;
	}

	/* Print the error message */
	base = 0;
	while (errmsg[base]!=0) {
		end = restart = findbreak(&errmsg[base],0,availablewidth);
		restart += base;
		while (errmsg[restart]==' ')  restart++;
		fprintf(stdout,"%s%.*s\n",prefix,end,&errmsg[base]);
		base = restart;
	}
}
