#ifndef LEMON_MEMORY_H
#define LEMON_MEMORY_H

extern void memory_error();
#define MemoryCheck(X) if ((X) == 0) { \
	memory_error(); \
}

#endif // LEMON_MEMORY_H
