#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "lemon_option.h"

/************************ From the file "option.c" **************************/
static char **argv;
static struct s_options *op;
static FILE *errstream;

#define ISOPT(X) ((X)[0]=='-'||(X)[0]=='+'||strchr((X),'=')!=0)

// ----------------------------------------------------------------
// Print the command line with a caret pointing to the k-th character
// of the nth field.
static void errline(int n, int k, FILE *err) {
	int spcnt, i;
	spcnt = 0;
	if (argv[0])  fprintf(err,"%s",argv[0]);
	spcnt = strlen(argv[0]) + 1;
	for(i=1; i<n && argv[i]; i++){
		fprintf(err," %s",argv[i]);
		spcnt += strlen(argv[i]+1);
	}
	spcnt += k;
	for(; argv[i]; i++) fprintf(err," %s",argv[i]);
	if (spcnt<20) {
		fprintf(err,"\n%*s^-- here\n",spcnt,"");
	} else {
		fprintf(err,"\n%*shere --^\n",spcnt-7,"");
	}
}

// ----------------------------------------------------------------
// Return the index of the N-th non-switch argument.  Return -1
// if N is out of range.
static int argindex(int n) {
	int i;
	int dashdash = 0;
	if (argv!=0 && *argv!=0) {
		for(i=1; argv[i]; i++){
			if (dashdash || !ISOPT(argv[i])) {
				if (n==0)  return i;
				n--;
			}
			if (strcmp(argv[i],"--")==0)  dashdash = 1;
		}
	}
	return -1;
}

static char emsg[] = "Command line syntax error: ";

// ----------------------------------------------------------------
// Process a flag command-line argument.
static int handleflags(int i, FILE *err) {
	int v;
	int errcnt = 0;
	int j;
	for(j=0; op[j].label; j++){
		if (strncmp(&argv[i][1],op[j].label,strlen(op[j].label))==0)  break;
	}
	v = argv[i][0]=='-' ? 1 : 0;
	if (op[j].label==0) {
		if (err) {
			fprintf(err,"%sundefined option.\n",emsg);
			errline(i,1,err);
		}
		errcnt++;
	} else if (op[j].type==OPT_FLAG) {
		*((int*)op[j].arg) = v;
	} else if (op[j].type==OPT_FFLAG) {
		(*(void(*)())(op[j].arg))(v);
	} else if (op[j].type==OPT_FSTR) {
		(*(void(*)())(op[j].arg))(&argv[i][2]);
	} else {
		if (err) {
			fprintf(err,"%smissing argument on switch.\n",emsg);
			errline(i,1,err);
		}
		errcnt++;
	}
	return errcnt;
}

// ----------------------------------------------------------------
// Process a command-line switch which has an argument.
static int handleswitch(int i, FILE *err) {
	int lv = 0;
	double dv = 0.0;
	char *sv = 0, *end;
	char *cp;
	int j;
	int errcnt = 0;
	cp = strchr(argv[i],'=');
	*cp = 0;
	for(j=0; op[j].label; j++){
		if (strcmp(argv[i],op[j].label)==0)  break;
	}
	*cp = '=';
	if (op[j].label==0) {
		if (err) {
			fprintf(err,"%sundefined option.\n",emsg);
			errline(i,0,err);
		}
		errcnt++;
	} else {
		cp++;
		switch (op[j].type) {
			case OPT_FLAG:
			case OPT_FFLAG:
				if (err) {
					fprintf(err,"%soption requires an argument.\n",emsg);
					errline(i,0,err);
				}
				errcnt++;
				break;
			case OPT_DBL:
			case OPT_FDBL:
				dv = strtod(cp,&end);
				if (*end) {
					if (err) {
						fprintf(err,"%sillegal character in floating-point argument.\n",emsg);
						errline(i,((unsigned long)end)-(unsigned long)argv[i],err);
					}
					errcnt++;
				}
				break;
			case OPT_INT:
			case OPT_FINT:
				lv = strtol(cp,&end,0);
				if (*end) {
					if (err) {
						fprintf(err,"%sillegal character in integer argument.\n",emsg);
						errline(i,((unsigned long)end)-(unsigned long)argv[i],err);
					}
					errcnt++;
				}
				break;
			case OPT_STR:
			case OPT_FSTR:
				sv = cp;
				break;
		}
		switch (op[j].type) {
			case OPT_FLAG:
			case OPT_FFLAG:
				break;
			case OPT_DBL:
				*(double*)(op[j].arg) = dv;
				break;
			case OPT_FDBL:
				(*(void(*)())(op[j].arg))(dv);
				break;
			case OPT_INT:
				*(int*)(op[j].arg) = lv;
				break;
			case OPT_FINT:
				(*(void(*)())(op[j].arg))((int)lv);
				break;
			case OPT_STR:
				*(char**)(op[j].arg) = sv;
				break;
			case OPT_FSTR:
				(*(void(*)())(op[j].arg))(sv);
				break;
		}
	}
	return errcnt;
}

// ----------------------------------------------------------------
int OptInit(char **a, struct s_options *o, FILE *err) {
	int errcnt = 0;
	argv = a;
	op = o;
	errstream = err;
	if (argv && *argv && op) {
		int i;
		for(i=1; argv[i]; i++){
			if (argv[i][0]=='+' || argv[i][0]=='-') {
				errcnt += handleflags(i,err);
			} else if (strchr(argv[i],'=')) {
				errcnt += handleswitch(i,err);
			}
		}
	}
	if (errcnt>0) {
		fprintf(err,"Valid command line options for \"%s\" are:\n",*a);
		OptPrint();
		exit(1);
	}
	return 0;
}

// ----------------------------------------------------------------
int OptNArgs() {
	int cnt = 0;
	int dashdash = 0;
	int i;
	if (argv!=0 && argv[0]!=0) {
		for(i=1; argv[i]; i++){
			if (dashdash || !ISOPT(argv[i]))  cnt++;
			if (strcmp(argv[i],"--")==0)  dashdash = 1;
		}
	}
	return cnt;
}

// ----------------------------------------------------------------
char *OptArg(int n) {
	int i;
	i = argindex(n);
	return i>=0 ? argv[i] : 0;
}

// ----------------------------------------------------------------
void OptErr(int n) {
	int i;
	i = argindex(n);
	if (i>=0)  errline(i,0,errstream);
}

// ----------------------------------------------------------------
void OptPrint() {
	int i;
	int max, len;
	max = 0;
	for(i=0; op[i].label; i++){
		len = strlen(op[i].label) + 1;
		switch (op[i].type) {
			case OPT_FLAG:
			case OPT_FFLAG:
				break;
			case OPT_INT:
			case OPT_FINT:
				len += 9;       /* length of "<integer>" */
				break;
			case OPT_DBL:
			case OPT_FDBL:
				len += 6;       /* length of "<real>" */
				break;
			case OPT_STR:
			case OPT_FSTR:
				len += 8;       /* length of "<string>" */
				break;
		}
		if (len>max)  max = len;
	}
	for(i=0; op[i].label; i++){
		switch (op[i].type) {
			case OPT_FLAG:
			case OPT_FFLAG:
				fprintf(errstream,"  -%-*s  %s\n",max,op[i].label,op[i].message);
				break;
			case OPT_INT:
			case OPT_FINT:
				fprintf(errstream,"  %s=<integer>%*s  %s\n",op[i].label,
					(int)(max-strlen(op[i].label)-9),"",op[i].message);
				break;
			case OPT_DBL:
			case OPT_FDBL:
				fprintf(errstream,"  %s=<real>%*s  %s\n",op[i].label,
					(int)(max-strlen(op[i].label)-6),"",op[i].message);
				break;
			case OPT_STR:
			case OPT_FSTR:
				fprintf(errstream,"  %s=<string>%*s  %s\n",op[i].label,
					(int)(max-strlen(op[i].label)-8),"",op[i].message);
				break;
		}
	}
}
