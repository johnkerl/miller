#ifndef FILE_READER_MMAP_H
#define FILE_READER_MMAP_H

// xxx rename to mmap_file_reader
typedef struct _file_reader_mmap_state_t {
	char* sol;
	char* eof;
	int   fd;
} file_reader_mmap_state_t;

file_reader_mmap_state_t file_reader_mmap_open(char* file_name);
void file_reader_mmap_close(file_reader_mmap_state_t* pstate);

#endif // FILE_READER_MMAP_H
