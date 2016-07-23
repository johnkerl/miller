#ifndef FILE_OUTPUT_MODE_H
#define FILE_OUTPUT_MODE_H

typedef enum _file_output_mode_t {
	MODE_WRITE,
	MODE_APPEND,
} file_output_mode_t;

static inline char* get_mode_string(file_output_mode_t file_output_mode) {
	return file_output_mode == MODE_APPEND ? "a" : "w";
}
static inline char* get_mode_desc(file_output_mode_t file_output_mode) {
	return file_output_mode == MODE_APPEND ? "append" : "write";
}

#endif // FILE_OUTPUT_MODE_H
