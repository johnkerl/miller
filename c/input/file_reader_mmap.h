// ================================================================
// Abstraction layer for mmapped file-read logic.
// ================================================================

#ifndef FILE_READER_MMAP_H
#define FILE_READER_MMAP_H

typedef struct _file_reader_mmap_state_t {
	char* sol;
	char* eof;
	int   fd;
} file_reader_mmap_state_t;

file_reader_mmap_state_t* file_reader_mmap_open(char* prepipe, char* file_name);
void file_reader_mmap_close(file_reader_mmap_state_t* pstate, char* prepipe);

void* file_reader_mmap_vopen(void* pvstate, char* prepipe, char* file_name);
void file_reader_mmap_vclose(void* pvstate, void* pvhandle, char* prepipe);

#endif // FILE_READER_MMAP_H
