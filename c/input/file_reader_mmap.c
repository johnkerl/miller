#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/stat.h>
#include <sys/mman.h>
#include "lib/mlrutil.h"
#include "file_reader_mmap.h"

static char empty_buf[1] = { 0 };

file_reader_mmap_state_t* file_reader_mmap_open(char* file_name) {
	file_reader_mmap_state_t* pstate = mlr_malloc_or_die(sizeof(file_reader_mmap_state_t));
	pstate->fd = open(file_name, O_RDONLY);
	if (pstate->fd < 0) {
		perror("open");
		exit(1);
	}
	struct stat stat;
	if (fstat(pstate->fd, &stat) < 0) {
		perror("fstat");
		exit(1);
	}
	if (stat.st_size == 0) {
		// mmap doesn't allow us to map zero-length files but zero-length files do exist.
		pstate->sol = &empty_buf[0];
	} else {
		pstate->sol = mmap(NULL, (size_t)stat.st_size, PROT_READ|PROT_WRITE, MAP_FILE|MAP_PRIVATE, pstate->fd, (off_t)0);
		if (pstate->sol == MAP_FAILED) {
			perror("mmap");
			exit(1);
		}
	}
	pstate->eof = pstate->sol + stat.st_size;
	return pstate;
}

void file_reader_mmap_close(file_reader_mmap_state_t* pstate) {
	if (close(pstate->fd) < 0) {
		perror("close");
		exit(1);
	}
	free(pstate);
}

void* file_reader_mmap_vopen(char* file_name) {
	return file_reader_mmap_open(file_name);
}
void file_reader_mmap_vclose(void* pvhandle) {
	file_reader_mmap_close(pvhandle);
}
