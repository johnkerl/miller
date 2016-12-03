#ifndef LEMON_STRUCTS_H
#define LEMON_STRUCTS_H

typedef enum {B_FALSE=0, B_TRUE} Boolean;

// Symbols (terminals and nonterminals) of the grammar are stored in the following:
struct symbol {
	char *name;              /* Name of the symbol */
	int index;               /* Index number for this symbol */
	enum {
		TERMINAL,
		NONTERMINAL
	} type;                  /* Symbols are all either TERMINALS or NTs */
	struct rule *rule;       /* Linked list of rules of this (if an NT) */
	struct symbol *fallback; /* fallback token in case this token doesn't parse */
	int prec;                /* Precedence if defined (-1 otherwise) */
	enum e_assoc {
		LEFT,
		RIGHT,
		NONE,
		UNK
	} assoc;                 /* Associativity if predecence is defined */
	char *firstset;          /* First-set for all rules of this symbol */
	Boolean lambda;          /* True if NT and can generate an empty string */
	char *destructor;        /* Code which executes whenever this symbol is
				** popped from the stack during error processing */
	int destructorln;        /* Line number of destructor code */
	char *datatype;          /* The data type of information held by this
				** object. Only used if type==NONTERMINAL */
	int dtnum;               /* The data type number.  In the parser, the value
				** stack is a union.  The .yy%d element of this
				** union is the correct data type for this object */
};

// Each production rule in the grammar is stored in the following structure.
struct rule {
	struct symbol *lhs;      /* Left-hand side of the rule */
	char *lhsalias;          /* Alias for the LHS (NULL if none) */
	int ruleline;            /* Line number for the rule */
	int nrhs;                /* Number of RHS symbols */
	struct symbol **rhs;     /* The RHS symbols */
	char **rhsalias;         /* An alias for each RHS symbol (NULL if none) */
	int line;                /* Line number at which code begins */
	char *code;              /* The code executed when this rule is reduced */
	struct symbol *precsym;  /* Precedence symbol for this rule */
	int index;               /* An index number for this rule */
	Boolean canReduce;       /* True if this rule is ever reduced */
	struct rule *nextlhs;    /* Next rule with the same LHS */
	struct rule *next;       /* Next rule in the global list */
};

// A configuration is a production rule of the grammar together with
// a mark (dot) showing how much of that rule has been processed so far.
// Configurations also contain a follow-set which is a list of terminal
// symbols which are allowed to immediately follow the end of the rule.
// Every configuration is recorded as an instance of the following:
struct config {
	struct rule *rp;         // The rule upon which the configuration is based
	int dot;                 // The parse point
	char *fws;               // Follow-set for this configuration only
	struct plink *fplp;      // Follow-set forward propagation links
	struct plink *bplp;      // Follow-set backwards propagation links
	struct state *stp;       // Pointer to state which contains this
	enum {
		COMPLETE,            // The status is used during followset and
		INCOMPLETE           // shift computations
	} status;
	struct config *next;     // Next configuration in the state
	struct config *bp;       // The next basis configuration
};

// Every shift or reduce operation is stored as one of the following
struct action {
	struct symbol *sp;       // The look-ahead symbol
	enum e_action {
		SHIFT,
		ACCEPT,
		REDUCE,
		ERROR,
		CONFLICT,                // Was a reduce, but part of a conflict
		SH_RESOLVED,             // Was a shift.  Precedence resolved conflict
		RD_RESOLVED,             // Was reduce.  Precedence resolved conflict
		NOT_USED                 // Deleted by compression
	} type;
	union {
		struct state *stp;     // The new state, if a shift
		struct rule *rp;       // The rule, if a reduce
	} x;
	struct action *next;     // Next action for this state
	struct action *collide;  // Next action with the same hash
};

// Each state of the generated parser's finite state machine
// is encoded as an instance of the following structure.
struct state {
	struct config *bp;       // The basis configurations for this state
	struct config *cfp;      // All configurations in this set
	int index;               // Sequencial number for this state
	struct action *ap;       // Array of actions for this state
	int nTknAct, nNtAct;     // Number of actions on terminals and nonterminals
	int iTknOfst, iNtOfst;   // yy_action[] offset for terminals and nonterms
	int iDflt;               // Default action
};
#define NO_OFFSET (-2147483647)

// A followset propagation link indicates that the contents of one
// configuration followset should be propagated to another whenever
// the first changes.
struct plink {
	struct config *cfp;      // The configuration to which linked
	struct plink *next;      // The next propagate link
};

// The state vector for the entire parser generator is recorded as
// follows.  (LEMON uses no global variables and makes little use of
// static variables.  Fields in the following structure can be thought
// of as being global variables in the program.)
struct lemon {
	struct state **sorted;   // Table of states sorted by state number
	struct rule *rule;       // List of all rules
	int nstate;              // Number of states
	int nrule;               // Number of rules
	int nsymbol;             // Number of terminal and nonterminal symbols
	int nterminal;           // Number of terminal symbols
	struct symbol **symbols; // Sorted array of pointers to symbols
	int errorcnt;            // Number of errors
	struct symbol *errsym;   // The error symbol
	char *name;              // Name of the generated parser
	char *arg;               // Declaration of the 3th argument to parser
	char *tokentype;         // Type of terminal symbols in the parser stack
	char *vartype;           // The default type of non-terminal symbols
	char *start;             // Name of the start symbol for the grammar
	char *stacksize;         // Size of the parser stack
	char *include;           // Code to put at the start of the C file
	int   includeln;         // Line number for start of include code
	char *error;             // Code to execute when an error is seen
	int   errorln;           // Line number for start of error code
	char *overflow;          // Code to execute on a stack overflow
	int   overflowln;        // Line number for start of overflow code
	char *failure;           // Code to execute on parser failure
	int   failureln;         // Line number for start of failure code
	char *accept;            // Code to execute when the parser excepts
	int   acceptln;          // Line number for the start of accept code
	char *extracode;         // Code appended to the generated file
	int   extracodeln;       // Line number for the start of the extra code
	char *tokendest;         // Code to execute to destroy token data
	int   tokendestln;       // Line number for token destroyer code
	char *vardest;           // Code for the default non-terminal destructor
	int  vardestln;          // Line number for default non-term destructor code
	char *filename;          // Name of the input file
	char *outname;           // Name of the current output file
	char *tokenprefix;       // A prefix added to token names in the .h file
	int   nconflict;         // Number of parsing conflicts
	int   tablesize;         // Size of the parse tables
	int   basisflag;         // Print only basis configurations
	int   has_fallback;      // True if any %fallback is seen in the grammer
	char *argv0;             // Name of the program
};

#endif // LEMON_STRUCTS_H
