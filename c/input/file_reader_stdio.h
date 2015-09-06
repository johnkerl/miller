// ================================================================
// Abstraction layer for stdio file-read logic.
// ================================================================

#ifndef FILE_READER_STDIO_H
#define FILE_READER_STDIO_H

void* file_reader_stdio_vopen(void* pvstate, char* file_name);
void file_reader_stdio_vclose(void* pvstate, void* pvhandle);

#endif // FILE_READER_STDIO_H
