#ifndef LEMON_OPTION_H
#define LEMON_OPTION_H

struct s_options {
	enum {
		OPT_FLAG=1,
		OPT_INT,
		OPT_DBL,
		OPT_STR,
		OPT_FFLAG,
		OPT_FINT,
		OPT_FDBL,
		OPT_FSTR
	} type;
	char *label;
	char *arg;
	char *message;
};

int    OptInit(/* char**,struct s_options*,FILE* */);
int    OptNArgs(/* void */);
char  *OptArg(/* int */);
void   OptErr(/* int */);
void   OptPrint(/* void */);

#endif // LEMON_OPTION_H
