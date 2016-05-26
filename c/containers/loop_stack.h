// Holds broken/continued flags for loops (for-srec, for-oosvar, while, do-while).

#ifndef LOOP_STACK_H
#define LOOP_STACK_H

#define LOOP_BROKEN    0x8000
#define LOOP_CONTINUED 0x0100

typedef struct _loop_stack_t {
	int  num_used_minus_one;
	int  num_allocated;
	int* pframes;
} loop_stack_t;

loop_stack_t* loop_stack_alloc();
void loop_stack_free(loop_stack_t* pstack);

// To be used on entry to loop handler.
void loop_stack_push(loop_stack_t* pstack);

// To be used on exit from loop handler.
int loop_stack_pop(loop_stack_t* pstack);

// To be used by break/continue handler.
// NOTE: For efficiency the stack is **NOT** bounds-checked here. E.g. if set is done before a push,
// or after an emptying pop, behavior is unspecified.
void loop_stack_set(loop_stack_t* pstack, int bits);
void loop_stack_clear(loop_stack_t* pstack, int bits);

// To be used by loop handler.
// NOTE: For efficiency the stack is **NOT** bounds-checked here. E.g. if set is done before a push,
// or after an emptying pop, behavior is unspecified.
int loop_stack_get(loop_stack_t* pstack);

#endif // LOOP_STACK_H
