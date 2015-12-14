#ifndef BYTE_READER_H
#define BYTE_READER_H

struct _byte_reader_t;

// Abstract byte source for input. Impls are stdio from file/stdin, mmapped
// file, and string-backed (for unit test).

// pvstate will nominally hold FILE*, mmap fd and pointers, string-backing, etc.

// The open function should return TRUE on success and FALSE on failure.
// For the string reader, the char* argument is the backing string itself.
typedef int byte_reader_open_func_t(struct _byte_reader_t* pbr, char* prepipe, char* filename);

// The reader function should return a character, as an int. Reads past end of
// file should keep returning EOF, even if called multiple times.
typedef int   byte_reader_read_func_t(struct _byte_reader_t* pbr);

// The close function should close file pointers/descriptors, as well as any
// necessary heap-frees.
typedef void  byte_reader_close_func_t(struct _byte_reader_t* pbr);

typedef struct _byte_reader_t {
	void*                       pvstate;
	byte_reader_open_func_t*    popen_func;
	byte_reader_read_func_t*    pread_func;
	byte_reader_close_func_t*   pclose_func;
} byte_reader_t;

#endif // BYTE_READER_H
