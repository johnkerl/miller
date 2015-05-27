#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/stat.h>
#include <sys/mman.h>
#include "file_reader_mmap.h"

static char empty_buf[1];

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
	if (stat.st_size == 0) {
		// mmap doesn't allow us to map zero-length files but zero-length files do exist.
		state.sol = &empty_buf[0];
	} else {
		state.sol = mmap(NULL, (size_t)stat.st_size, PROT_READ|PROT_WRITE, MAP_FILE|MAP_PRIVATE, state.fd, (off_t)0);
		if (state.sol == MAP_FAILED) {
			perror("mmap");
			exit(1);
		}
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
