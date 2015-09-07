#ifndef STRING_BUILDER_H
#define STRING_BUILDER_H

typedef struct _string_builder_t {
	int used_length;
	int alloc_length;
	char* buffer;
} string_builder_t;

void  sb_init(string_builder_t* psb, int alloc_length);
void _sb_enlarge(string_builder_t* psb); // private method

static inline void sb_append_char(string_builder_t* psb, char c) {
	if (psb->used_length >= psb->alloc_length)
		_sb_enlarge(psb);
	psb->buffer[psb->used_length++] = c;
}

void  sb_append_string(string_builder_t* psb, char* s);
int   sb_is_empty(string_builder_t* psb);
// The caller should free() the return value:
char* sb_finish(string_builder_t* psb);

#endif // STRING_BUILDER_H
