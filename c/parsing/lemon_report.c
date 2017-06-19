#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <ctype.h>

#include "lemon_report.h"

#include "lemon_error.h"
#include "lemon_assert.h"
#include "lemon_string.h"
#include "lemon_action.h"
#include "lemon_dims.h"
#include "lemon_symbol.h"

/*
** Procedures for generating reports and tables in the LEMON parser generator.
*/

/* Generate a filename with the given suffix.  Space to hold the
** name comes from malloc() and must be freed by the calling
** function.
*/
static char *file_makename(struct lemon *lemp, char *suffix)
{
	char *name;
	char *cp;

	name = malloc (strlen(lemp->filename) + strlen(suffix) + 5) ;
	if (name==0) {
		fprintf(stderr,"Can't allocate space for a filename.\n");
		exit(1);
	}
	strcpy(name,lemp->filename);
	cp = strrchr(name,'.');
	if (cp)  *cp = 0;
	strcat(name,suffix);
	return name;
}
static char* capitalize(char* string) {
	char* p;
	for (p = string; *p; p++) {
		*p = toupper(*p);
	}
	return string;
}

/* Open a file with a name based on the name of the input file,
** but with a different (specified) suffix, and return a pointer
** to the stream */
static FILE *file_open(struct lemon *lemp, char *suffix, char *mode)
{
	FILE *fp;

	free(lemp->outname);
	lemp->outname = file_makename(lemp, suffix);
	fp = fopen(lemp->outname,mode);
	if (fp==0 && *mode=='w') {
		fprintf(stderr,"Can't open file \"%s\".\n",lemp->outname);
		lemp->errorcnt++;
		return 0;
	}
	return fp;
}

/* Duplicate the input file without comments and without actions on rules */
void Reprint(struct lemon *lemp)
{
	struct rule *rp;
	struct symbol *sp;
	int i, j, maxlen, len, ncolumns, skip;
	printf("// Reprint of input file \"%s\".\n// Symbols:\n",lemp->filename);
	maxlen = 10;
	for(i=0; i<lemp->nsymbol; i++){
		sp = lemp->symbols[i];
		len = strlen(sp->name);
		if (len>maxlen)  maxlen = len;
	}
	ncolumns = 76/(maxlen+5);
	if (ncolumns<1)  ncolumns = 1;
	skip = (lemp->nsymbol + ncolumns - 1)/ncolumns;
	for(i=0; i<skip; i++){
		printf("//");
		for(j=i; j<lemp->nsymbol; j+=skip){
			sp = lemp->symbols[j];
			assert (sp->index==j) ;
			printf(" %3d %-*.*s",j,maxlen,maxlen,sp->name);
		}
		printf("\n");
	}
	for(rp=lemp->rule; rp; rp=rp->next){
		printf("%s",rp->lhs->name);
/*    if (rp->lhsalias)  printf("(%s)",rp->lhsalias); */
		printf(" ::=");
		for(i=0; i<rp->nrhs; i++){
			printf(" %s",rp->rhs[i]->name);
/*      if (rp->rhsalias[i])  printf("(%s)",rp->rhsalias[i]); */
		}
		printf(".");
		if (rp->precsym)  printf(" [%s]",rp->precsym->name);
/*    if (rp->code)  printf("\n    %s",rp->code); */
		printf("\n");
	}
}

void ConfigPrint(FILE *fp, struct config *cfp)
{
	struct rule *rp;
	int i;
	rp = cfp->rp;
	fprintf(fp,"%s ::=",rp->lhs->name);
	for(i=0; i<=rp->nrhs; i++){
		if (i==cfp->dot)  fprintf(fp," *");
		if (i==rp->nrhs)  break;
		fprintf(fp," %s",rp->rhs[i]->name);
	}
}

/* #define TEST */
#ifdef TEST
/* Print a set */
static void SetPrint(FILE *out, char *set, struct lemon *lemp)
{
	int i;
	char *spacer;
	spacer = "";
	fprintf(out,"%12s[","");
	for(i=0; i<lemp->nterminal; i++){
		if (SetFind(set,i)) {
			fprintf(out,"%s%s",spacer,lemp->symbols[i]->name);
			spacer = " ";
		}
	}
	fprintf(out,"]\n");
}

/* Print a plink chain */
static void PlinkPrint(FILE *out, struct plink *plp, char *tag)
{
	while (plp) {
		fprintf(out,"%12s%s (state %2d) ","",tag,plp->cfp->stp->index);
		ConfigPrint(out,plp->cfp);
		fprintf(out,"\n");
		plp = plp->next;
	}
}
#endif

/* Print an action to the given file descriptor.  Return FALSE if
** nothing was actually printed.
*/
int PrintAction(struct action *ap, FILE *fp, int indent) {
	int result = 1;
	switch (ap->type) {
		case SHIFT:
			fprintf(fp,"%*s shift  %d",indent,ap->sp->name,ap->x.stp->index);
			break;
		case REDUCE:
			fprintf(fp,"%*s reduce %d",indent,ap->sp->name,ap->x.rp->index);
			break;
		case ACCEPT:
			fprintf(fp,"%*s accept",indent,ap->sp->name);
			break;
		case ERROR:
			fprintf(fp,"%*s error",indent,ap->sp->name);
			break;
		case CONFLICT:
			fprintf(fp,"%*s reduce %-3d ** Parsing conflict **",
				indent,ap->sp->name,ap->x.rp->index);
			break;
		case SH_RESOLVED:
		case RD_RESOLVED:
		case NOT_USED:
			result = 0;
			break;
	}
	return result;
}

/* Generate the "y.output" log file */
void ReportOutput(struct lemon *lemp)
{
	int i;
	struct state *stp;
	struct config *cfp;
	struct action *ap;
	FILE *fp;

	fp = file_open(lemp,".out","w");
	if (fp==0)  return;
	fprintf(fp," \b");
	for(i=0; i<lemp->nstate; i++){
		stp = lemp->sorted[i];
		fprintf(fp,"State %d:\n",stp->index);
		if (lemp->basisflag)  cfp=stp->bp;
		else                  cfp=stp->cfp;
		while (cfp) {
			char buf[20];
			if (cfp->dot==cfp->rp->nrhs) {
				sprintf(buf,"(%d)",cfp->rp->index);
				fprintf(fp,"    %5s ",buf);
			} else {
				fprintf(fp,"          ");
			}
			ConfigPrint(fp,cfp);
			fprintf(fp,"\n");
#ifdef TEST
			SetPrint(fp,cfp->fws,lemp);
			PlinkPrint(fp,cfp->fplp,"To  ");
			PlinkPrint(fp,cfp->bplp,"From");
#endif
			if (lemp->basisflag)  cfp=cfp->bp;
			else                  cfp=cfp->next;
		}
		fprintf(fp,"\n");
		for(ap=stp->ap; ap; ap=ap->next){
			if (PrintAction(ap,fp,30))  fprintf(fp,"\n");
		}
		fprintf(fp,"\n");
	}
	fclose(fp);
	return;
}

/* Search for the file "name" which is in the same directory as
** the exacutable */
static char *pathsearch(char *argv0, char *name, int modemask)
{
	char *pathlist;
	char *path,*cp;
	char c;

#ifdef __WIN32__
	cp = strrchr(argv0,'\\');
#else
	cp = strrchr(argv0,'/');
#endif
	if (cp) {
		c = *cp;
		*cp = 0;
		path = (char *)malloc (strlen(argv0) + strlen(name) + 2) ;
		if (path)  sprintf(path,"%s/%s",argv0,name);
		*cp = c;
	} else {
		extern char *getenv();
		pathlist = getenv("PATH");
		if (pathlist==0)  pathlist = ".:/bin:/usr/bin";
		path = (char *)malloc (strlen(pathlist)+strlen(name)+2) ;
		if (path!=0) {
			while (*pathlist) {
				cp = strchr(pathlist,':');
				if (cp==0)  cp = &pathlist[strlen(pathlist)];
				c = *cp;
				*cp = 0;
				sprintf(path,"%s/%s",pathlist,name);
				*cp = c;
				if (c==0)  pathlist = "";
				else pathlist = &cp[1];
				if (access(path,modemask)==0)  break;
			}
		}
	}
	return path;
}

/* Given an action, compute the integer value for that action
** which is to be put in the action table of the generated machine.
** Return negative if no action should be generated.
*/
static int compute_action(struct lemon *lemp, struct action *ap)
{
	int act;
	switch (ap->type) {
		case SHIFT:  act = ap->x.stp->index;               break;
		case REDUCE: act = ap->x.rp->index + lemp->nstate; break;
		case ERROR:  act = lemp->nstate + lemp->nrule;     break;
		case ACCEPT: act = lemp->nstate + lemp->nrule + 1; break;
		default:     act = -1; break;
	}
	return act;
}

#define LINESIZE 1000
/* The next cluster of routines are for reading the template file
** and writing the results to the generated parser */
/* The first function transfers data from "in" to "out" until
** a line is seen which begins with "%%".  The line number is
** tracked.
**
** if name!=0, then any word that begin with "Parse" is changed to
** begin with *name instead.
*/
static void tplt_xfer(
	char *name,
	FILE *in,
	FILE *out,
	int *lineno)
{
	int i, iStart;
	char line[LINESIZE];
	while (fgets(line,LINESIZE,in) && (line[0]!='%' || line[1]!='%')) {
		(*lineno)++;
		iStart = 0;
		if (name) {
			for(i=0; line[i]; i++){
				if (line[i]=='P' && strncmp(&line[i],"Parse",5)==0
					&& (i==0 || !isalpha(line[i-1]))
			 ) {
					if (i>iStart)  fprintf(out,"%.*s",i-iStart,&line[iStart]);
					fprintf(out,"%s",name);
					i += 4;
					iStart = i+1;
				}
			}
		}
		fprintf(out,"%s",&line[iStart]);
	}
}

/* The next function finds the template file and opens it, returning
** a pointer to the opened file. */
static FILE *tplt_open(struct lemon *lemp)
{
	static char templatename[] = "lempar.c";
	char buf[1000];
	FILE *in;
	char *tpltname;
	char *cp;

	cp = strrchr(lemp->filename,'.');
	if (cp) {
		sprintf(buf,"%.*s.lt",(int)(cp-lemp->filename),lemp->filename);
	} else {
		sprintf(buf,"%s.lt",lemp->filename);
	}
	if (access(buf,004)==0) {
		tpltname = buf;
	} else if (access(templatename,004)==0) {
		tpltname = templatename;
	} else {
		tpltname = pathsearch(lemp->argv0,templatename,0);
	}
	if (tpltname==0) {
		fprintf(stderr,"Can't find the parser driver template file \"%s\".\n",
		templatename);
		lemp->errorcnt++;
		return 0;
	}
	in = fopen(tpltname,"r");
	if (in==0) {
		fprintf(stderr,"Can't open the template file \"%s\".\n",templatename);
		lemp->errorcnt++;
		return 0;
	}
	return in;
}

/* Print a string to the file and keep the linenumber up to date */
static void tplt_print(
	FILE *out,
	struct lemon *lemp,
	char *str,
	int strln,
	int *lineno,
	int suppress_line_directives)
{
	if (str == 0)
		return;
	if (!suppress_line_directives) {
		fprintf(out,"#line %d \"%s\"\n",strln,lemp->filename); (*lineno)++;
	}
	while (*str) {
		if (*str=='\n')  (*lineno)++;
		putc(*str,out);
		str++;
	}
	if (!suppress_line_directives) {
		fprintf(out,"\n#line %d \"%s\"\n",*lineno+2,lemp->outname); (*lineno)+=2;
	}
	return;
}

/*
** The following routine emits code for the destructor for the
** symbol sp
*/
void emit_destructor_code(
	FILE *out,
	struct symbol *sp,
	struct lemon *lemp,
	int *lineno,
	int suppress_line_directives)
{
	char *cp = 0;

	int linecnt = 0;
	if (sp->type==TERMINAL) {
		cp = lemp->tokendest;
		if (cp==0)  return;
		if (!suppress_line_directives) {
			fprintf(out,"#line %d \"%s\"\n{",lemp->tokendestln,lemp->filename);
		}
	} else if (sp->destructor) {
		cp = sp->destructor;
		if (!suppress_line_directives) {
			fprintf(out,"#line %d \"%s\"\n{",sp->destructorln,lemp->filename);
		}
	} else if (lemp->vardest) {
		cp = lemp->vardest;
		if (cp==0)  return;
		if (!suppress_line_directives) {
			fprintf(out,"#line %d \"%s\"\n",lemp->vardestln,lemp->filename);
		}
		fprintf(out, "{");
	} else {
		assert (0) ;  /* Cannot happen */
	}
	for(; *cp; cp++){
		if (*cp=='$' && cp[1]=='$') {
			fprintf(out,"(yypminor->yy%d)",sp->dtnum);
			cp++;
			continue;
		}
		if (*cp=='\n')  linecnt++;
		fputc(*cp,out);
	}
	(*lineno) += 3 + linecnt;
	fprintf(out,"}\n");
	if (!suppress_line_directives) {
		fprintf(out,"#line %d \"%s\"\n",*lineno,lemp->outname);
	}
	return;
}

/*
** Return TRUE (non-zero) if the given symbol has a destructor.
*/
int has_destructor(struct symbol *sp, struct lemon *lemp)
{
	int ret;
	if (sp->type==TERMINAL) {
		ret = lemp->tokendest!=0;
	} else {
		ret = lemp->vardest!=0 || sp->destructor!=0;
	}
	return ret;
}

/*
** Append text to a dynamically allocated string.  If zText is 0 then
** reset the string to be empty again.  Always return the complete text
** of the string (which is overwritten with each call).
**
** n bytes of zText are stored.  If n==0 then all of zText up to the first
** \000 terminator is stored.  zText can contain up to two instances of
** %d.  The values of p1 and p2 are written into the first and second
** %d.
**
** If n==-1, then the previous character is overwritten.
*/
static char *append_str(char *zText, int n, int p1, int p2) {
	static char *z = 0;
	static int alloced = 0;
	static int used = 0;
	int i, c;
	char zInt[40];

	if (zText==0) {
		used = 0;
		return z;
	}
	if (n<=0) {
		if (n<0) {
			used += n;
			assert (used>=0) ;
		}
		n = strlen(zText);
	}
	if (n+sizeof(zInt)*2+used >= alloced) {
		alloced = n + sizeof(zInt)*2 + used + 200;
		z = realloc(z,  alloced);
	}
	if (z==0)  return "";
	while (n-- > 0) {
		c = *(zText++);
		if (c=='%' && zText[0]=='d') {
			sprintf(zInt, "%d", p1);
			p1 = p2;
			strcpy(&z[used], zInt);
			used += strlen(&z[used]);
			zText++;
			n--;
		} else {
			z[used++] = c;
		}
	}
	z[used] = 0;
	return z;
}

/*
** zCode is a string that is the action associated with a rule.  Expand
** the symbols in this string so that the refer to elements of the parser
** stack.  Return a new string stored in space obtained from malloc.
*/
static void translate_code(struct lemon *lemp, struct rule *rp) {
	char *cp, *xp;
	int i;
	char lhsused = 0;    /* True if the LHS element has been used */
	char used[MAXRHS];   /* True for each RHS element which is used */

	for(i=0; i<rp->nrhs; i++) used[i] = 0;
	lhsused = 0;

	append_str(0,0,0,0);
	for(cp=rp->code; *cp; cp++){
		if (isalpha(*cp) && (cp==rp->code || (!isalnum(cp[-1]) && cp[-1]!='_'))) {
			char saved;
			for(xp= &cp[1]; isalnum(*xp) || *xp=='_'; xp++);
			saved = *xp;
			*xp = 0;
			if (rp->lhsalias && strcmp(cp,rp->lhsalias)==0) {
				append_str("yygotominor.yy%d",0,rp->lhs->dtnum,0);
				cp = xp;
				lhsused = 1;
			} else {
				for(i=0; i<rp->nrhs; i++){
					if (rp->rhsalias[i] && strcmp(cp,rp->rhsalias[i])==0) {
						if (cp!=rp->code && cp[-1]=='@') {
							/* If the argument is of the form @X then substituted
							** the token number of X, not the value of X */
							append_str("yymsp[%d].major",-1,i-rp->nrhs+1,0);
						} else {
							append_str("yymsp[%d].minor.yy%d",0,
												 i-rp->nrhs+1,rp->rhs[i]->dtnum);
						}
						cp = xp;
						used[i] = 1;
						break;
					}
				}
			}
			*xp = saved;
		}
		append_str(cp, 1, 0, 0);
	} /* End loop */

	/* Check to make sure the LHS has been used */
	if (rp->lhsalias && !lhsused) {
		ErrorMsg(lemp->filename,rp->ruleline,
			"Label \"%s\" for \"%s(%s)\" is never used.",
				rp->lhsalias,rp->lhs->name,rp->lhsalias);
		lemp->errorcnt++;
	}

	/* Generate destructor code for RHS symbols which are not used in the
	** reduce code */
	for(i=0; i<rp->nrhs; i++){
		if (rp->rhsalias[i] && !used[i]) {
			ErrorMsg(lemp->filename,rp->ruleline,
				"Label %s for \"%s(%s)\" is never used.",
				rp->rhsalias[i],rp->rhs[i]->name,rp->rhsalias[i]);
			lemp->errorcnt++;
		} else if (rp->rhsalias[i]==0) {
			if (has_destructor(rp->rhs[i],lemp)) {
				append_str("  yy_destructor(%d,&yymsp[%d].minor);\n", 0,
					 rp->rhs[i]->index,i-rp->nrhs+1);
			} else {
				/* No destructor defined for this term */
			}
		}
	}
	cp = append_str(0,0,0,0);
	rp->code = Strsafe(cp);
}

/*
** Generate code which executes when the rule "rp" is reduced.  Write
** the code to "out".  Make sure lineno stays up-to-date.
*/
static void emit_code(
	FILE *out,
	struct rule *rp,
	struct lemon *lemp,
	int *lineno,
	int suppress_line_directives)
{
	char *cp;
	int linecnt = 0;

	/* Generate code to do the reduce action */
	if (rp->code) {
		if (!suppress_line_directives) {
			fprintf(out,"#line %d \"%s\"\n",rp->line,lemp->filename);
		}
		fprintf(out,"{");
		fprintf(out,"%s",rp->code);
		for(cp=rp->code; *cp; cp++){
			if (*cp=='\n')  linecnt++;
		} /* End loop */
		(*lineno) += 3 + linecnt;
		fprintf(out,"}\n");
		if (!suppress_line_directives) {
			fprintf(out,"#line %d \"%s\"\n",*lineno,lemp->outname);
		}
	} /* End if (rp->code)  */

	return;
}

/*
** Print the definition of the union used for the parser's data stack.
** This union contains fields for every possible data type for tokens
** and nonterminals.  In the process of computing and printing this
** union, also set the ".dtnum" field of every terminal and nonterminal
** symbol.
*/
void print_stack_union(
	FILE *out,                  /* The output stream */
	struct lemon *lemp,         /* The main info structure for this parser */
	int *plineno,               /* Pointer to the line number */
	int mhflag)                 /* True if generating makeheaders output */
{
	int lineno = *plineno;    /* The line number of the output */
	char **types;             /* A hash table of datatypes */
	int arraysize;            /* Size of the "types" array */
	int maxdtlength;          /* Maximum length of any ".datatype" field. */
	char *stddt;              /* Standardized name for a datatype */
	int i,j;                  /* Loop counters */
	unsigned hash;            /* For hashing the name of a type */
	char *name;               /* Name of the parser */

	/* Allocate and initialize types[] and allocate stddt[] */
	arraysize = lemp->nsymbol * 2;
	types = (char**)malloc (arraysize * sizeof(char*)) ;
	for(i=0; i<arraysize; i++) types[i] = 0;
	maxdtlength = 0;
	if (lemp->vartype) {
		maxdtlength = strlen(lemp->vartype);
	}
	for(i=0; i<lemp->nsymbol; i++){
		int len;
		struct symbol *sp = lemp->symbols[i];
		if (sp->datatype==0)  continue;
		len = strlen(sp->datatype);
		if (len>maxdtlength)  maxdtlength = len;
	}
	stddt = (char*)malloc (maxdtlength*2 + 1) ;
	if (types==0 || stddt==0) {
		fprintf(stderr,"Out of memory.\n");
		exit(1);
	}

	/* Build a hash table of datatypes. The ".dtnum" field of each symbol
	** is filled in with the hash index plus 1.  A ".dtnum" value of 0 is
	** used for terminal symbols.  If there is no %default_type defined then
	** 0 is also used as the .dtnum value for nonterminals which do not specify
	** a datatype using the %type directive.
	*/
	for(i=0; i<lemp->nsymbol; i++){
		struct symbol *sp = lemp->symbols[i];
		char *cp;
		if (sp==lemp->errsym) {
			sp->dtnum = arraysize+1;
			continue;
		}
		if (sp->type!=NONTERMINAL || (sp->datatype==0 && lemp->vartype==0)) {
			sp->dtnum = 0;
			continue;
		}
		cp = sp->datatype;
		if (cp==0)  cp = lemp->vartype;
		j = 0;
		while (isspace(*cp))  cp++;
		while (*cp)  stddt[j++] = *cp++;
		while (j>0 && isspace(stddt[j-1]))  j--;
		stddt[j] = 0;
		hash = 0;
		for(j=0; stddt[j]; j++){
			hash = hash*53 + (unsigned)stddt[j];
		}
		hash = (hash & 0x7fffffff)%arraysize;
		while (types[hash]) {
			if (strcmp(types[hash],stddt)==0) {
				sp->dtnum = hash + 1;
				break;
			}
			hash++;
			if (hash>=arraysize)  hash = 0;
		}
		if (types[hash]==0) {
			sp->dtnum = hash + 1;
			types[hash] = (char*)malloc (strlen(stddt)+1) ;
			if (types[hash]==0) {
				fprintf(stderr,"Out of memory.\n");
				exit(1);
			}
			strcpy(types[hash],stddt);
		}
	}

	/* Print out the definition of YYTOKENTYPE and YYMINORTYPE */
	name = lemp->name ? lemp->name : "Parse";
	lineno = *plineno;
	if (mhflag) { fprintf(out,"#if INTERFACE\n"); lineno++; }
	fprintf(out,"#define %sTOKENTYPE %s\n",name,
		lemp->tokentype?lemp->tokentype:"void*");  lineno++;
	if (mhflag) { fprintf(out,"#endif\n"); lineno++; }
	fprintf(out,"typedef union {\n"); lineno++;
	fprintf(out,"  %sTOKENTYPE yy0;\n",name); lineno++;
	for(i=0; i<arraysize; i++){
		if (types[i]==0)  continue;
		fprintf(out,"  %s yy%d;\n",types[i],i+1); lineno++;
		free(types[i]);
	}
	fprintf(out,"  int yy%d;\n",lemp->errsym->dtnum); lineno++;
	free(stddt);
	free(types);
	fprintf(out,"} YYMINORTYPE;\n"); lineno++;
	*plineno = lineno;
}

/*
** Return the name of a C datatype able to represent values between
** lwr and upr, inclusive.
*/
static const char *minimum_size_type(int lwr, int upr){
	if (lwr>=0) {
		if (upr<=255) {
			return "unsigned char";
		} else if (upr<65535) {
			return "unsigned short int";
		} else {
			return "unsigned int";
		}
	} else if (lwr>=-127 && upr<=127) {
		return "signed char";
	} else if (lwr>=-32767 && upr<32767) {
		return "short";
	} else {
		return "int";
	}
}

/*
** Each state contains a set of token transaction and a set of
** nonterminal transactions.  Each of these sets makes an instance
** of the following structure.  An array of these structures is used
** to order the creation of entries in the yy_action[] table.
*/
struct axset {
	struct state *stp;   /* A pointer to a state */
	int isTkn;           /* True to use tokens.  False for non-terminals */
	int nAction;         /* Number of actions */
};

/*
** Compare to axset structures for sorting purposes
*/
static int axset_compare(const void *a, const void *b){
	struct axset *p1 = (struct axset*)a;
	struct axset *p2 = (struct axset*)b;
	return p2->nAction - p1->nAction;
}

/* Generate C source code for the parser */
/* Output in makeheaders format if mhflag==true */
void ReportTable(struct lemon *lemp, int mhflag, int suppress_line_directives)
{
	FILE *out, *in;
	char line[LINESIZE];
	int  lineno;
	struct state *stp;
	struct action *ap;
	struct rule *rp;
	struct acttab *pActtab;
	int i, j, n;
	char *name;
	int mnTknOfst, mxTknOfst;
	int mnNtOfst, mxNtOfst;
	struct axset *ax;

	in = tplt_open(lemp);
	if (in==0)  return;
	out = file_open(lemp,".c","w");
	if (out==0) {
		fclose(in);
		return;
	}
	lineno = 1;
	tplt_xfer(lemp->name,in,out,&lineno);

	/* Generate the include code, if any */
	tplt_print(out,lemp,lemp->include,lemp->includeln,&lineno, suppress_line_directives);
	if (mhflag) {
		char *name = file_makename(lemp, ".h");
		fprintf(out,"#include \"%s\"\n", name); lineno++;
		free(name);
	}
	tplt_xfer(lemp->name,in,out,&lineno);

	/* Generate #defines for all tokens */
	if (mhflag) {
		char *prefix;
		fprintf(out,"#if INTERFACE\n"); lineno++;
		if (lemp->tokenprefix)  prefix = lemp->tokenprefix;
		else                    prefix = "";
		for(i=1; i<lemp->nterminal; i++){
			fprintf(out,"#define %s%-30s %2d\n",prefix,lemp->symbols[i]->name,i);
			lineno++;
		}
		fprintf(out,"#endif\n"); lineno++;
	}
	tplt_xfer(lemp->name,in,out,&lineno);

	/* Generate the defines */
	fprintf(out,"#define YYCODETYPE %s\n",
		minimum_size_type(0, lemp->nsymbol+5)); lineno++;
	fprintf(out,"#define YYNOCODE %d\n",lemp->nsymbol+1);  lineno++;
	fprintf(out,"#define YYACTIONTYPE %s\n",
		minimum_size_type(0, lemp->nstate+lemp->nrule+5));  lineno++;
	print_stack_union(out,lemp,&lineno,mhflag);
	if (lemp->stacksize) {
		if (atoi(lemp->stacksize)<=0) {
			ErrorMsg(lemp->filename,0,
"Illegal stack size: [%s].  The stack size should be an integer constant.",
				lemp->stacksize);
			lemp->errorcnt++;
			lemp->stacksize = "800";
		}
		fprintf(out,"#define YYSTACKDEPTH %s\n",lemp->stacksize);  lineno++;
	} else {
		fprintf(out,"#define YYSTACKDEPTH 800\n");  lineno++;
	}
	if (mhflag) {
		fprintf(out,"#if INTERFACE\n"); lineno++;
	}
	name = lemp->name ? lemp->name : "Parse";
	if (lemp->arg && lemp->arg[0]) {
		int i;
		i = strlen(lemp->arg);
		while (i>=1 && isspace(lemp->arg[i-1]))  i--;
		while (i>=1 && (isalnum(lemp->arg[i-1]) || lemp->arg[i-1]=='_'))  i--;
		fprintf(out,"#define %sARG_SDECL %s;\n",name,lemp->arg);  lineno++;
		fprintf(out,"#define %sARG_PDECL ,%s\n",name,lemp->arg);  lineno++;
		fprintf(out,"#define %sARG_FETCH %s = pparser->%s\n",
			name,lemp->arg,&lemp->arg[i]);  lineno++;
		fprintf(out,"#define %sARG_STORE pparser->%s = %s\n",
			name,&lemp->arg[i],&lemp->arg[i]);  lineno++;
	} else {
		fprintf(out,"#define %sARG_SDECL\n",name);  lineno++;
		fprintf(out,"#define %sARG_PDECL\n",name);  lineno++;
		fprintf(out,"#define %sARG_FETCH\n",name); lineno++;
		fprintf(out,"#define %sARG_STORE\n",name); lineno++;
	}
	if (mhflag) {
		fprintf(out,"#endif\n"); lineno++;
	}
	fprintf(out,"#define YYNSTATE %d\n",lemp->nstate);  lineno++;
	fprintf(out,"#define YYNRULE %d\n",lemp->nrule);  lineno++;
	fprintf(out,"#define YYERRORSYMBOL %d\n",lemp->errsym->index);  lineno++;
	fprintf(out,"#define YYERRSYMDT yy%d\n",lemp->errsym->dtnum);  lineno++;
	if (lemp->has_fallback) {
		fprintf(out,"#define YYFALLBACK 1\n");  lineno++;
	}
	tplt_xfer(lemp->name,in,out,&lineno);

	/* Generate the action table and its associates:
	**
	**  yy_action[]        A single table containing all actions.
	**  yy_lookahead[]     A table containing the lookahead for each entry in
	**                     yy_action.  Used to detect hash collisions.
	**  yy_shift_ofst[]    For each state, the offset into yy_action for
	**                     shifting terminals.
	**  yy_reduce_ofst[]   For each state, the offset into yy_action for
	**                     shifting non-terminals after a reduce.
	**  yy_default[]       Default action for each state.
	*/

	/* Compute the actions on all states and count them up */
	ax = malloc (sizeof(ax[0])*lemp->nstate*2) ;
	if (ax==0) {
		fprintf(stderr,"malloc failed\n");
		exit(1);
	}
	for(i=0; i<lemp->nstate; i++){
		stp = lemp->sorted[i];
		stp->nTknAct = stp->nNtAct = 0;
		stp->iDflt = lemp->nstate + lemp->nrule;
		stp->iTknOfst = NO_OFFSET;
		stp->iNtOfst = NO_OFFSET;
		for(ap=stp->ap; ap; ap=ap->next){
			if (compute_action(lemp,ap)>=0) {
				if (ap->sp->index<lemp->nterminal) {
					stp->nTknAct++;
				} else if (ap->sp->index<lemp->nsymbol) {
					stp->nNtAct++;
				} else {
					stp->iDflt = compute_action(lemp, ap);
				}
			}
		}
		ax[i*2].stp = stp;
		ax[i*2].isTkn = 1;
		ax[i*2].nAction = stp->nTknAct;
		ax[i*2+1].stp = stp;
		ax[i*2+1].isTkn = 0;
		ax[i*2+1].nAction = stp->nNtAct;
	}
	mxTknOfst = mnTknOfst = 0;
	mxNtOfst = mnNtOfst = 0;

	/* Compute the action table.  In order to try to keep the size of the
	** action table to a minimum, the heuristic of placing the largest action
	** sets first is used.
	*/
	qsort(ax, lemp->nstate*2, sizeof(ax[0]), axset_compare);
	pActtab = acttab_alloc();
	for(i=0; i<lemp->nstate*2 && ax[i].nAction>0; i++){
		stp = ax[i].stp;
		if (ax[i].isTkn) {
			for(ap=stp->ap; ap; ap=ap->next){
				int action;
				if (ap->sp->index>=lemp->nterminal)  continue;
				action = compute_action(lemp, ap);
				if (action<0)  continue;
				acttab_action(pActtab, ap->sp->index, action);
			}
			stp->iTknOfst = acttab_insert(pActtab);
			if (stp->iTknOfst<mnTknOfst)  mnTknOfst = stp->iTknOfst;
			if (stp->iTknOfst>mxTknOfst)  mxTknOfst = stp->iTknOfst;
		} else {
			for(ap=stp->ap; ap; ap=ap->next){
				int action;
				if (ap->sp->index<lemp->nterminal)  continue;
				if (ap->sp->index==lemp->nsymbol)  continue;
				action = compute_action(lemp, ap);
				if (action<0)  continue;
				acttab_action(pActtab, ap->sp->index, action);
			}
			stp->iNtOfst = acttab_insert(pActtab);
			if (stp->iNtOfst<mnNtOfst)  mnNtOfst = stp->iNtOfst;
			if (stp->iNtOfst>mxNtOfst)  mxNtOfst = stp->iNtOfst;
		}
	}
	free(ax);

	/* Output the yy_action table */
	fprintf(out,"static YYACTIONTYPE yy_action[] = {\n"); lineno++;
	n = acttab_size(pActtab);
	for(i=j=0; i<n; i++){
		int action = acttab_yyaction(pActtab, i);
		if (action<0)  action = lemp->nsymbol + lemp->nrule + 2;
		if (j==0)  fprintf(out," /* %5d */ ", i);
		fprintf(out, " %4d,", action);
		if (j==9 || i==n-1) {
			fprintf(out, "\n"); lineno++;
			j = 0;
		} else {
			j++;
		}
	}
	fprintf(out, "};\n"); lineno++;

	/* Output the yy_lookahead table */
	fprintf(out,"static YYCODETYPE yy_lookahead[] = {\n"); lineno++;
	for(i=j=0; i<n; i++){
		int la = acttab_yylookahead(pActtab, i);
		if (la<0)  la = lemp->nsymbol;
		if (j==0)  fprintf(out," /* %5d */ ", i);
		fprintf(out, " %4d,", la);
		if (j==9 || i==n-1) {
			fprintf(out, "\n"); lineno++;
			j = 0;
		} else {
			j++;
		}
	}
	fprintf(out, "};\n"); lineno++;

	/* Output the yy_shift_ofst[] table */
	fprintf(out, "#define YY_SHIFT_USE_DFLT (%d)\n", mnTknOfst-1); lineno++;
	fprintf(out, "static %s yy_shift_ofst[] = {\n",
					minimum_size_type(mnTknOfst-1, mxTknOfst)); lineno++;
	n = lemp->nstate;
	for(i=j=0; i<n; i++){
		int ofst;
		stp = lemp->sorted[i];
		ofst = stp->iTknOfst;
		if (ofst==NO_OFFSET)  ofst = mnTknOfst - 1;
		if (j==0)  fprintf(out," /* %5d */ ", i);
		fprintf(out, " %4d,", ofst);
		if (j==9 || i==n-1) {
			fprintf(out, "\n"); lineno++;
			j = 0;
		} else {
			j++;
		}
	}
	fprintf(out, "};\n"); lineno++;

	/* Output the yy_reduce_ofst[] table */
	fprintf(out, "#define YY_REDUCE_USE_DFLT (%d)\n", mnNtOfst-1); lineno++;
	fprintf(out, "static %s yy_reduce_ofst[] = {\n",
					minimum_size_type(mnNtOfst-1, mxNtOfst)); lineno++;
	n = lemp->nstate;
	for(i=j=0; i<n; i++){
		int ofst;
		stp = lemp->sorted[i];
		ofst = stp->iNtOfst;
		if (ofst==NO_OFFSET)  ofst = mnNtOfst - 1;
		if (j==0)  fprintf(out," /* %5d */ ", i);
		fprintf(out, " %4d,", ofst);
		if (j==9 || i==n-1) {
			fprintf(out, "\n"); lineno++;
			j = 0;
		} else {
			j++;
		}
	}
	fprintf(out, "};\n"); lineno++;

	/* Output the default action table */
	fprintf(out, "static YYACTIONTYPE yy_default[] = {\n"); lineno++;
	n = lemp->nstate;
	for(i=j=0; i<n; i++){
		stp = lemp->sorted[i];
		if (j==0)  fprintf(out," /* %5d */ ", i);
		fprintf(out, " %4d,", stp->iDflt);
		if (j==9 || i==n-1) {
			fprintf(out, "\n"); lineno++;
			j = 0;
		} else {
			j++;
		}
	}
	fprintf(out, "};\n"); lineno++;
	tplt_xfer(lemp->name,in,out,&lineno);

	/* Generate the table of fallback tokens.
	*/
	if (lemp->has_fallback) {
		for(i=0; i<lemp->nterminal; i++){
			struct symbol *p = lemp->symbols[i];
			if (p->fallback==0) {
				fprintf(out, "    0,  /* %10s => nothing */\n", p->name);
			} else {
				fprintf(out, "  %3d,  /* %10s => %s */\n", p->fallback->index,
					p->name, p->fallback->name);
			}
			lineno++;
		}
	}
	tplt_xfer(lemp->name, in, out, &lineno);

	/* Generate a table containing the symbolic name of every symbol
	*/
	for(i=0; i<lemp->nsymbol; i++){
		sprintf(line,"\"%s\",",lemp->symbols[i]->name);
		fprintf(out,"  %-15s",line);
		if ((i&3)==3) { fprintf(out,"\n"); lineno++; }
	}
	if ((i&3)!=0) { fprintf(out,"\n"); lineno++; }
	tplt_xfer(lemp->name,in,out,&lineno);

	/* Generate a table containing a text string that describes every
	** rule in the rule set of the grammer.  This information is used
	** when tracing REDUCE actions.
	*/
	for(i=0, rp=lemp->rule; rp; rp=rp->next, i++){
		assert (rp->index==i) ;
		fprintf(out," /* %3d */ \"%s ::=", i, rp->lhs->name);
		for(j=0; j<rp->nrhs; j++) fprintf(out," %s",rp->rhs[j]->name);
		fprintf(out,"\",\n"); lineno++;
	}
	tplt_xfer(lemp->name,in,out,&lineno);

	/* Generate code which executes every time a symbol is popped from
	** the stack while processing errors or while destroying the parser.
	** (In other words, generate the %destructor actions)
	*/
	if (lemp->tokendest) {
		for(i=0; i<lemp->nsymbol; i++){
			struct symbol *sp = lemp->symbols[i];
			if (sp==0 || sp->type!=TERMINAL)  continue;
			fprintf(out,"    case %d:\n",sp->index); lineno++;
		}
		for(i=0; i<lemp->nsymbol && lemp->symbols[i]->type!=TERMINAL; i++);
		if (i<lemp->nsymbol) {
			emit_destructor_code(out,lemp->symbols[i],lemp,&lineno, suppress_line_directives);
			fprintf(out,"      break;\n"); lineno++;
		}
	}
	for(i=0; i<lemp->nsymbol; i++){
		struct symbol *sp = lemp->symbols[i];
		if (sp==0 || sp->type==TERMINAL || sp->destructor==0)  continue;
		fprintf(out,"    case %d:\n",sp->index); lineno++;

		/* Combine duplicate destructors into a single case */
		for(j=i+1; j<lemp->nsymbol; j++){
			struct symbol *sp2 = lemp->symbols[j];
			if (sp2 && sp2->type!=TERMINAL && sp2->destructor
					&& sp2->dtnum==sp->dtnum
					&& strcmp(sp->destructor,sp2->destructor)==0) {
				 fprintf(out,"    case %d:\n",sp2->index); lineno++;
				 sp2->destructor = 0;
			}
		}

		emit_destructor_code(out,lemp->symbols[i],lemp,&lineno, suppress_line_directives);
		fprintf(out,"      break;\n"); lineno++;
	}
	if (lemp->vardest) {
		struct symbol *dflt_sp = 0;
		for(i=0; i<lemp->nsymbol; i++){
			struct symbol *sp = lemp->symbols[i];
			if (sp==0 || sp->type==TERMINAL ||
					sp->index<=0 || sp->destructor!=0)  continue;
			fprintf(out,"    case %d:\n",sp->index); lineno++;
			dflt_sp = sp;
		}
		if (dflt_sp!=0) {
			emit_destructor_code(out,dflt_sp,lemp,&lineno, suppress_line_directives);
			fprintf(out,"      break;\n"); lineno++;
		}
	}
	tplt_xfer(lemp->name,in,out,&lineno);

	/* Generate code which executes whenever the parser stack overflows */
	tplt_print(out,lemp,lemp->overflow,lemp->overflowln,&lineno, suppress_line_directives);
	tplt_xfer(lemp->name,in,out,&lineno);

	/* Generate the table of rule information
	**
	** Note: This code depends on the fact that rules are number
	** sequentually beginning with 0.
	*/
	for(rp=lemp->rule; rp; rp=rp->next){
		fprintf(out,"  { %d, %d },\n",rp->lhs->index,rp->nrhs); lineno++;
	}
	tplt_xfer(lemp->name,in,out,&lineno);

	/* Generate code which execution during each REDUCE action */
	for(rp=lemp->rule; rp; rp=rp->next){
		if (rp->code)  translate_code(lemp, rp);
	}
	for(rp=lemp->rule; rp; rp=rp->next){
		struct rule *rp2;
		if (rp->code==0)  continue;
		fprintf(out,"      case %d:\n",rp->index); lineno++;
		for(rp2=rp->next; rp2; rp2=rp2->next){
			if (rp2->code==rp->code) {
				fprintf(out,"      case %d:\n",rp2->index); lineno++;
				rp2->code = 0;
			}
		}
		emit_code(out,rp,lemp,&lineno, suppress_line_directives);
		fprintf(out,"        break;\n"); lineno++;
	}
	tplt_xfer(lemp->name,in,out,&lineno);

	/* Generate code which executes if a parse fails */
	tplt_print(out,lemp,lemp->failure,lemp->failureln,&lineno, suppress_line_directives);
	tplt_xfer(lemp->name,in,out,&lineno);

	/* Generate code which executes when a syntax error occurs */
	tplt_print(out,lemp,lemp->error,lemp->errorln,&lineno, suppress_line_directives);
	tplt_xfer(lemp->name,in,out,&lineno);

	/* Generate code which executes when the parser accepts its input */
	tplt_print(out,lemp,lemp->accept,lemp->acceptln,&lineno, suppress_line_directives);
	tplt_xfer(lemp->name,in,out,&lineno);

	/* Append any addition code the user desires */
	tplt_print(out,lemp,lemp->extracode,lemp->extracodeln,&lineno, suppress_line_directives);

	fclose(in);
	fclose(out);
	return;
}

/* Generate a header file for the parser */
void ReportHeader(struct lemon *lemp)
{
	FILE *out, *in;
	char *prefix;
	char line[LINESIZE];
	char pattern[LINESIZE];
	int i;

	if (lemp->tokenprefix)  prefix = lemp->tokenprefix;
	else                    prefix = "";
	in = file_open(lemp,".h","r");
	if (in) {
		for(i=1; i<lemp->nterminal && fgets(line,LINESIZE,in); i++){
			sprintf(pattern,"#define %s%-30s %2d\n",prefix,lemp->symbols[i]->name,i);
			if (strcmp(line,pattern))  break;
		}
		fclose(in);
		if (i==lemp->nterminal) {
			/* No change in the file.  Don't rewrite it. */
			return;
		}
	}
	out = file_open(lemp,".h","w");
	if (out) {
		char* ident = capitalize(file_makename(lemp, "_h"));
		fprintf(out, "#ifndef %s\n", ident);
		fprintf(out, "#define %s\n", ident);
		fprintf(out, "\n");
		for(i=1; i<lemp->nterminal; i++){
			fprintf(out,"#define %s%-30s %2d\n",prefix,lemp->symbols[i]->name,i);
		}
		fprintf(out, "\n");
		fprintf(out, "#endif // %s\n", ident);
		fclose(out);
	}
	return;
}

/* Reduce the size of the action tables, if possible, by making use
** of defaults.
**
** In this version, we take the most frequent REDUCE action and make
** it the default.  Only default a reduce if there are more than one.
*/

void CompressTables(struct lemon *lemp)
{
	struct state *stp;
	struct action *ap, *ap2;
	struct rule *rp, *rp2, *rbest;
	int nbest, n;
	int i;

	for(i=0; i<lemp->nstate; i++){
		stp = lemp->sorted[i];
		nbest = 0;
		rbest = 0;

		for(ap=stp->ap; ap; ap=ap->next){
			if (ap->type!=REDUCE)  continue;
			rp = ap->x.rp;
			if (rp==rbest)  continue;
			n = 1;
			for(ap2=ap->next; ap2; ap2=ap2->next){
				if (ap2->type!=REDUCE)  continue;
				rp2 = ap2->x.rp;
				if (rp2==rbest)  continue;
				if (rp2==rp)  n++;
			}
			if (n>nbest) {
				nbest = n;
				rbest = rp;
			}
		}

		/* Do not make a default if the number of rules to default
		** is not at least 2 */
		if (nbest<2)  continue;


		/* Combine matching REDUCE actions into a single default */
		for(ap=stp->ap; ap; ap=ap->next){
			if (ap->type==REDUCE && ap->x.rp==rbest)  break;
		}
		assert (ap) ;
		ap->sp = Symbol_new("{default}");
		for(ap=ap->next; ap; ap=ap->next){
			if (ap->type==REDUCE && ap->x.rp==rbest)  ap->type = NOT_USED;
		}
		stp->ap = Action_sort(stp->ap);
	}
}

