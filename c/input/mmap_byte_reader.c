#include <string.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/stat.h>
#include <sys/mman.h>
#include "input/byte_readers.h"
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"

static char empty_buf[1] = { 0 };

typedef struct _mmap_byte_reader_state_t {
	char* filename;
	int   fd;
	char* sof;
	char* p;
	char* eof;
} mmap_byte_reader_state_t;

static int  mmap_byte_reader_open_func(struct _byte_reader_t* pbr, char* prepipe, char* filename);
static int  mmap_byte_reader_read_func(struct _byte_reader_t* pbr);
static void mmap_byte_reader_close_func(struct _byte_reader_t* pbr, char* prepipe);

// ----------------------------------------------------------------
byte_reader_t* mmap_byte_reader_alloc() {
	byte_reader_t* pbr = mlr_malloc_or_die(sizeof(byte_reader_t));

	pbr->pvstate     = NULL;
	pbr->popen_func  = mmap_byte_reader_open_func;
	pbr->pread_func  = mmap_byte_reader_read_func;
	pbr->pclose_func = mmap_byte_reader_close_func;

	return pbr;
}

void mmap_byte_reader_free(byte_reader_t* pbr) {
	mmap_byte_reader_state_t* pstate = pbr->pvstate;
	if (pstate != NULL) {
		free(pstate->filename); // null-ok semantics
	}
	free(pbr);
}

// ----------------------------------------------------------------
static int mmap_byte_reader_open_func(struct _byte_reader_t* pbr, char* prepipe, char* filename) {
	// popen is a stdio construct, not an mmap construct, and it can't be supported here.
	if (prepipe != NULL) {
		fprintf(stderr, "%s: coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}

	mmap_byte_reader_state_t* pstate = mlr_malloc_or_die(sizeof(mmap_byte_reader_state_t));
	pstate->filename = mlr_strdup_or_die(filename);
	pstate->fd = open(filename, O_RDONLY);
	if (pstate->fd < 0) {
		perror("open");
		fprintf(stderr, "%s: Couldn't open \"%s\" for read.\n", MLR_GLOBALS.argv0, filename);
		exit(1);
	}

	struct stat stat;
	if (fstat(pstate->fd, &stat) < 0) {
		perror("fstat");
		fprintf(stderr, "%s: could not fstat \"%s\"\n", MLR_GLOBALS.argv0, filename);
		exit(1);
	}
	if (stat.st_size == 0) {
		// mmap doesn't allow us to map zero-length files but zero-length files do exist.
		pstate->sof = &empty_buf[0];
	} else {
		pstate->sof = mmap(NULL, (size_t)stat.st_size, PROT_READ|PROT_WRITE, MAP_FILE|MAP_PRIVATE,
			pstate->fd, (off_t)0);
		if (pstate->sof == MAP_FAILED) {
			perror("mmap");
			fprintf(stderr, "%s: could not mmap \"%s\"\n", MLR_GLOBALS.argv0, filename);
			exit(1);
		}
	}
	pstate->eof = pstate->sof + stat.st_size;
	pstate->p = pstate->sof;
	pbr->pvstate = pstate;
	return TRUE;
}

static int mmap_byte_reader_read_func(struct _byte_reader_t* pbr) {
	mmap_byte_reader_state_t* pstate = pbr->pvstate;
	if (pstate->p >= pstate->eof) {
		return EOF;
	} else {
		int c = *pstate->p;
		pstate->p++;
		return c;
	}
}

static void mmap_byte_reader_close_func(struct _byte_reader_t* pbr, char* prepipe) {
	mmap_byte_reader_state_t* pstate = pbr->pvstate;
	if (close(pstate->fd) < 0) {
		perror("close");
		fprintf(stderr, "%s: close error on file \"%s\".\n", MLR_GLOBALS.argv0, pstate->filename);
		exit(1);
	}
}
