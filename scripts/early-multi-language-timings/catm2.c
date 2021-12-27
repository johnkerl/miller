#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/stat.h>
#include <sys/mman.h>

typedef struct _file_reader_mmap_state_t {
	char* sol;
	char* eof;
	int   fd;
} file_reader_mmap_state_t;

file_reader_mmap_state_t file_reader_mmap_open(char* file_name) {
	file_reader_mmap_state_t state;
	state.fd = open(file_name, O_RDONLY);
	if (state.fd < 0) {
		perror("open");
		exit(1);
	}
	struct stat stat;
	if (fstat(state.fd, &stat) < 0) {
		perror("fstat");
		exit(1);
	}
	state.sol = mmap(NULL, (size_t)stat.st_size, PROT_READ|PROT_WRITE, MAP_FILE|MAP_PRIVATE, state.fd, (off_t)0);
	if (state.sol == MAP_FAILED) {
		perror("mmap");
		exit(1);
	}
	state.eof = state.sol + stat.st_size;
	return state;
}

void file_reader_mmap_close(file_reader_mmap_state_t* pstate) {
	if (close(pstate->fd) < 0) {
		perror("close");
		exit(1);
	}
}

// ----------------------------------------------------------------
static void emit(char* sol, char* eol, FILE* output_stream) {
	 size_t ntowrite = eol - sol;
     size_t nwritten = fwrite(sol, 1, ntowrite, output_stream);
	 if (nwritten != ntowrite) {
		perror("fwrite");
		exit(1);
	 }
	 fputc('\n', output_stream);
}

// ----------------------------------------------------------------
// xxx params/state:
// * ctor:  char*file_name
// * reads: currptr, eofptr
// * dtor:  int fd
static int do_stream(char* file_name) {
	FILE* output_stream = stdout;

	file_reader_mmap_state_t state = file_reader_mmap_open(file_name);

	char* eol;
	char* p = state.sol;

	while (p < state.eof) {
		if (*p == '\n') {
			*p = 0;
			eol = p;
			emit(state.sol, eol, output_stream);
			p++;
			state.sol = p;
		} else {
			p++;
		}
	}

	file_reader_mmap_close(&state);

	return 1;
}

// ================================================================
int main(int argc, char** argv) {
	int ok = 1;
	for (int argi = 1; argi < argc; argi++) {
	    ok = do_stream(argv[argi]);
	}
	return ok ? 0 : 1;
}
