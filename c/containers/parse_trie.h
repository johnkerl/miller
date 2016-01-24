#ifndef PARSE_TRIE_H
#define PARSE_TRIE_H

// ----------------------------------------------------------------
// This is for parsing of Miller RFC-CSV data, not for parsing of DSL
// expressions for the put and filter verbs. This is used instead of flex or
// lemon since such parsers, by design, must read to end of input to complete
// parsing. Here, by contrast, we want to split the input stream by delimiters,
// identifying one record at a time. This is so that data may be processed in a
// streaming manner rather than an ingest-all manner.

struct _parse_trie_node_t;
typedef struct _parse_trie_node_t {
	struct _parse_trie_node_t* pnexts[256];
	char c;      // current character at this node
	int  stridx; // which string was stored ending here; -1 if not end of string.
	int  strlen; // length of string stored ending here; -1 if not end of string.
} parse_trie_node_t;
typedef struct _parse_trie_t {
	parse_trie_node_t* past;
	int maxlen;
} parse_trie_t;

// ----------------------------------------------------------------
parse_trie_t* parse_trie_alloc();
void parse_trie_free(parse_trie_t* ptrie);
void parse_trie_print(parse_trie_t* ptrie);
void parse_trie_add_string(parse_trie_t* ptrie, char* string, int stridx);

// ----------------------------------------------------------------
// Example input:
// * string 0 is "a"
// * string 1 is "aa"
// * buf is "aaabc"
// Output:
// * return value is TRUE
// * stridx is 1 since longest match is "aa"
// * matchlen is 2 since "aa" has length 2

// This does a longest-prefix match: input data "\"\nabcdefg" is matched
// against "\"\n" rather than against "\"".

// We assume that enough data has been peeked into the ring buffer for the
// parse-trie's maxlen.  There is no check here. This function is called on
// every single character of RFC-CSV input data so the error-checking would be
// inefficient here, as well as misplaced.

// The start of buffer (sob), buflen, and mask attributes are nominally
// presented from a ring_buffer object.

static inline int parse_trie_ring_match(parse_trie_t* ptrie, char* buf, int sob, int buflen, int mask,
	int* pstridx, int* pmatchlen)
{
	parse_trie_node_t* pnode = ptrie->past;
	parse_trie_node_t* pnext;
	parse_trie_node_t* pterm = NULL;
	for (int i = 0; i < buflen; i++) {
		char c = buf[(sob+i)&mask];
		pnext = pnode->pnexts[(unsigned char) c];
		if (pnext == NULL)
			break;
		if (pnext->strlen > 0) {
			pterm = pnext;
		}
		pnode = pnext;
	}
	if (pterm == NULL) {
		return FALSE;
	} else {
		*pstridx   = pterm->stridx;
		*pmatchlen = pterm->strlen;
		return TRUE;
	}
}

static inline int parse_trie_match(parse_trie_t* ptrie, char* p, char* e, int* pstridx, int* pmatchlen) {
	parse_trie_node_t* pnode = ptrie->past;
	parse_trie_node_t* pnext;
	parse_trie_node_t* pterm = NULL;
	for ( ; p < e; p++) {
		char c = *p;
		pnext = pnode->pnexts[(unsigned char) c];
		if (pnext == NULL)
			break;
		if (pnext->strlen > 0) {
			pterm = pnext;
		}
		pnode = pnext;
	}
	if (pterm == NULL) {
		return FALSE;
	} else {
		*pstridx   = pterm->stridx;
		*pmatchlen = pterm->strlen;
		return TRUE;
	}
}

#endif // PARSE_TRIE_H
