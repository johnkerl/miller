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
int parse_trie_match(parse_trie_t* ptrie, char* buf, int buflen, int* pstridx, int* pmatchlen);

#endif // PARSE_TRIE_H
