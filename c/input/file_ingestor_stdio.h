// ================================================================
// Abstraction layer for stdio file-read logic, when ingesting entirely into memory.
// ================================================================

#ifndef FILE_INGESTOR_STDIO_H
#define FILE_INGESTOR_STDIO_H

typedef struct _file_ingestor_stdio_state_t {
	char* sof;
	char* eof;
} file_ingestor_stdio_state_t;

// The vclose method frees memory. The nop flavor does not. The latter is used by the JSON-input
// lrec reader: lrecs have keys/values pointing into parsed-JSON structures which in turn have
// values pointing to the ingested-file buffer.  This is done for the sake of performance, to reduce
// data-copies. But it also means we can't free ingested files after ingesting their contents into
// lrecs, since the lrecs in question might be retained after the input-file closes.  Example: mlr
// sort on multiple files.
void* file_ingestor_stdio_vopen(void* pvstate, char* prepipe, char* file_name);
void  file_ingestor_stdio_vclose(void* pvstate, void* pvhandle, char* prepipe);
void  file_ingestor_stdio_nop_vclose(void* pvstate, void* pvhandle, char* prepipe);

#endif // FILE_INGESTOR_STDIO_H
