// ================================================================
// Abstraction layer for stdio file-read logic, when ingesting entirely into memory.
// ================================================================

#ifndef FILE_INGESTOR_STDIO_H
#define FILE_INGESTOR_STDIO_H

typedef struct _file_ingestor_stdio_state_t {
	char* sof;
	char* eof;
} file_ingestor_stdio_state_t;

void* file_ingestor_stdio_vopen(void* pvstate, char* prepipe, char* file_name);
void  file_ingestor_stdio_vclose(void* pvstate, void* pvhandle, char* prepipe);

#endif // FILE_INGESTOR_STDIO_H
