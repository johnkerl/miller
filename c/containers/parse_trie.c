#include <stdlib.h>
#include "lib/mlrutil.h"
#include "containers/parse_trie.h"

static parse_trie_node_t* parse_trie_node_alloc();
static void parse_trie_print_aux(parse_trie_node_t* pnode, int depth);

// ----------------------------------------------------------------
parse_trie_t* parse_trie_alloc() {
	parse_trie_t* ptrie = mlr_malloc_or_die(sizeof(parse_trie_t));
	return ptrie;
}

static parse_trie_node_t* parse_trie_node_alloc() {
	parse_trie_node_t* pnode = mlr_malloc_or_die(sizeof(parse_trie_node_t));
	for (int i = 0; i < 256; i++)
		pnode->pnexts[i] = NULL;
	pnode->c = 0;
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
	printf("[pnode=%p]\n", pnode);
	for (int i = 0; i < depth; i++)
		printf("  ");
	printf("c=%c,stridx=%d,strlen=%d\n", pnode->c, pnode->stridx, pnode->strlen);
	printf("[pnode=%p]\n", pnode);
	for (int i = 0; i < 256; i++) {
		parse_trie_node_t* pnext = pnode->pnexts[i];
		if (pnext != NULL) {
			parse_trie_print_aux(pnext, depth+1);
		}
	}
}

// ----------------------------------------------------------------
void parse_trie_add_string(parse_trie_t* ptrie, char* string) {
	if (ptrie->proot == NULL) {
		ptrie->proot = parse_trie_node_alloc();
	}
	//parse_trie_node_t* pnode = ptrie->proot;
	return; // xxx stub
}

// ----------------------------------------------------------------
// Example:
// * string 0 is "a"
// * string 1 is "aa"
// * buf is "aaabc"
int parse_trie_match(parse_trie_t* ptrie, char* buf, int buflen, int* pstridx, int* pmatchlen) {
	parse_trie_node_t* pnode = ptrie->proot;
	parse_trie_node_t* pnext;
	int i;
	for (i = 0; i < buflen; i++) {
		char c = buf[i];
		pnext = pnode->pnexts[(unsigned) c];
		// xxx not quite right
		if (pnode == NULL)
			return FALSE;
		pnode = pnext;
	}
	return TRUE; // xxx stub
}
