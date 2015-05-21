#ifndef FILE_READER_MMAP_H
#define FILE_READER_MMAP_H

// xxx rename to mmap_file_reader
typedef struct _mmap_reader_state_t {
	char* sol;
	char* eof;
	int   fd;
} mmap_reader_state_t;

mmap_reader_state_t mmap_reader_open(char* file_name);
void mmap_reader_close(mmap_reader_state_t* pstate);

#endif // FILE_READER_MMAP_H
