#ifndef PARSE_TRIE_H
#define PARSE_TRIE_H

// ----------------------------------------------------------------
// xxx cmt for data parse, not DSL parse.
// xxx cmt why not flex/lemon: want streaming, not ingest-all.
struct _parse_trie_node_t;
typedef struct _parse_trie_node_t {
	struct _parse_trie_node_t* pnexts[256];
	char c;      // current character at this node
	int  stridx; // which string was stored ending here; -1 if not end of string.
	int  strlen; // length of string stored ending here; -1 if not end of string.
} parse_trie_node_t;
typedef struct _parse_trie_t {
	parse_trie_node_t* proot;
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

// xxx cmt longest-prefix match

// xxx cmt assumption that enough has been peeked -- there is no check here.
// (recall that this is called on every single byte of each data file, e.g.
// for a 1GB data file we'll be in this function a billion times ... it needs
// to be tight.)

// xxx cmt re sob/mask semantics (ring buffer)

static inline int parse_trie_match(parse_trie_t* ptrie, char* buf, int sob, int buflen, int mask, int* pstridx, int* pmatchlen) {
	parse_trie_node_t* pnode = ptrie->proot;
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

#endif // PARSE_TRIE_H
