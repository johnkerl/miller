#include <stdlib.h>
#include <ctype.h>
#include "lib/mlrutil.h"
#include "containers/parse_trie.h"

static parse_trie_node_t* parse_trie_node_alloc(char c);
static void parse_trie_print_aux(parse_trie_node_t* pnode, int depth);
static void parse_trie_add_string_aux(parse_trie_node_t* pnode, char* string, int stridx, int len);

// ----------------------------------------------------------------
parse_trie_t* parse_trie_alloc() {
	parse_trie_t* ptrie = mlr_malloc_or_die(sizeof(parse_trie_t));
	ptrie->proot = parse_trie_node_alloc(0);
	ptrie->maxlen = 0;
	return ptrie;
}

static parse_trie_node_t* parse_trie_node_alloc(char c) {
	parse_trie_node_t* pnode = mlr_malloc_or_die(sizeof(parse_trie_node_t));
	for (int i = 0; i < 256; i++)
		pnode->pnexts[i] = NULL;
	pnode->c = c;
	pnode->stridx = -1;
	pnode->strlen = -1;
	return pnode;
}

// ----------------------------------------------------------------
void parse_trie_free(parse_trie_t* ptrie) {
	return; // xxx stub
}

// ----------------------------------------------------------------
void parse_trie_print(parse_trie_t* ptrie) {
	parse_trie_node_t* pnode = ptrie->proot;
	printf("PARSE TRIE DUMP START\n");
	if (pnode != NULL) {
		parse_trie_print_aux(pnode, 0);
	}
	printf("PARSE TRIE DUMP END\n");
}

static void parse_trie_print_aux(parse_trie_node_t* pnode, int depth) {
	for (int i = 0; i < depth; i++)
		printf("  ");
	printf("c=%c[%02x],stridx=%d,strlen=%d\n",
		isprint((unsigned char)pnode->c) ? pnode->c : '?',
		(unsigned)pnode->c,
		pnode->stridx,
		pnode->strlen);
	for (int i = 0; i < 256; i++) {
		parse_trie_node_t* pnext = pnode->pnexts[i];
		if (pnext != NULL)
			parse_trie_print_aux(pnext, depth+1);
	}
}

// ----------------------------------------------------------------
void parse_trie_add_string(parse_trie_t* ptrie, char* string, int stridx) {
	int len = strlen(string);
	parse_trie_add_string_aux(ptrie->proot, string, stridx, strlen(string));
	if (len > ptrie->maxlen)
		ptrie->maxlen = len;
}

static void parse_trie_add_string_aux(parse_trie_node_t* pnode, char* string, int stridx, int len) {
	char c = string[0];
	if (c == 0) {
		pnode->stridx = stridx;
		pnode->strlen = len;
	} else {
		parse_trie_node_t* pnext = pnode->pnexts[(unsigned char)c];
		if (pnext == NULL) {
			pnext = parse_trie_node_alloc(c);
			pnext->c = c;
			pnext->stridx = -1;
			pnext->strlen = -1;
			pnode->pnexts[(unsigned char)c] = pnext;
		}
		parse_trie_add_string_aux(pnext, &string[1], stridx, len);
	}
}
