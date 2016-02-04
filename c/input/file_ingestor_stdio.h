// ================================================================
// Abstraction layer for stdio file-read logic, when ingesting entirely into memory.
// ================================================================

#ifndef FILE_INGESTOR_STDIO_H
#define FILE_INGESTOR_STDIO_H

void* file_ingestor_stdio_vopen(void* pvstate, char* prepipe, char* file_name);
void  file_ingestor_stdio_vclose(void* pvstate, void* pvhandle, char* prepipe);

#endif // FILE_INGESTOR_STDIO_H
