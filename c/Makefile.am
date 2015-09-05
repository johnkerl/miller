SUBDIRS=	lib cli stream input mapping output containers dsls

AM_CPPFLAGS=	-I${srcdir}
AM_CFLAGS=	-Wall -Werror -std=gnu99 -O3
bin_PROGRAMS=	mlr
mlr_SOURCES=	mlrmain.c mlrvers.h
mlr_LDADD=	\
		cli/libcli.la \
		containers/libcontainers.la \
		stream/libstream.la \
		input/libinput.la \
		mapping/libmapping.la \
		output/liboutput.la \
		lib/libmlr.la \
		dsls/libfdsl.la \
		dsls/libpdsl.la \
		-lm

EXPERIMENTAL_READER_SRCS = \
lib/mlrutil.c \
lib/mlr_globals.c \
lib/string_builder.c \
input/file_reader_mmap.c \
input/old_peek_file_reader.c \
experimental/getlines.c

# ================================================================
# User-make: creates the executable and runs unit & regression tests (without
# valgrind)
# ----> This is the default target for anyone pulling the repo and trying to
#       build it to be able to use it. It just needs flex and the C compiler.
top: mlr tests

# Dev-make: updates the tags file, creates the executable, and runs unit &
# regression tests (with valgrind)
# ----> This is the target for a developer to run before a commit.
#       It requires ctags and valgrind in addition to flex and the C compiler.
dev: tags mlr dev-tests

# ----------------------------------------------------------------
# Other executable variants

# Debug version
mlrg: .always dsls
	$(CCDEBUG) $(NON_DSL_SRCS) $(PDSL_OBJS) $(FDSL_OBJS) $(LFLAGS) -o mlrg

# Memcheck version
mlrk: .always dsls
	$(CC) -DUSE_MCHECK $(NON_DSL_SRCS) $(PDSL_OBJS) $(FDSL_OBJS) $(LFLAGS) -o mlrk

# Profile version. Usage:
# * make mlrp
# * mlrp {arguments>
# * gprof mlrp gmon.out > myfile.txt
# Note: works on Linux; not on OSX.
mlrp: .always dsls
	$(CC) -g -pg $(NON_DSL_SRCS) $(PDSL_OBJS) $(FDSL_OBJS) $(LFLAGS) -o mlrp

# ================================================================
# User-tests: don't use valgrind: they may not have it, and valgrind is most
# useful for detecting errors at build time. They don't need it to produce an
# executable.

tests: unit-test reg-test

unit-test: test-mlrutil test-byte-readers test-parse-trie test-lrec test-string-builder test-join-bucket-keeper
	./test-mlrutil
	./test-byte-readers
	./test-parse-trie
	./test-lrec
	./test-string-builder
	./test-join-bucket-keeper
	@echo
	@echo DONE

reg-test:
	./test/run

# ----------------------------------------------------------------
# Use valgrind at dev time for additional reassurance.

dev-tests: dev-unit-test reg-test

# Unfortunately --error-exitcode=1 doesn't work well since there are
# valgrind-detected errors in stdlibs. :(
dev-unit-test: test-mlrutil test-byte-readers test-parse-trie test-lrec test-string-builder test-join-bucket-keeper
#	valgrind --leak-check=full --error-exitcode=1 ./a.out
	valgrind --leak-check=full ./test-mlrutil
	valgrind --leak-check=full ./test-byte-readers
	valgrind --leak-check=full ./test-parse-trie
	valgrind --leak-check=full ./test-lrec
	valgrind --leak-check=full ./test-string-builder
	valgrind --leak-check=full ./test-join-bucket-keeper
	@echo
	@echo DONE

# Run this after unit-test expected output has changed, and is verified to be
# OK. (Example: after adding new test cases in test/run.)
regtest-copy:
	cp test/output/out test/expected

# ----------------------------------------------------------------
# Standalone mains
# TODO: replace the interesting content with unit tests; jettison the rest

dheap:
	$(CCDEBUG) -D__DHEAP_MAIN__ containers/dheap.c lib/mlrutil.c

slls:
	$(CCDEBUG) -D__SLLS_MAIN__ containers/slls.c lib/mlrutil.c

sllv:
	$(CCDEBUG) -D__SLLV_MAIN__ containers/sllv.c lib/mlrutil.c

hmss:
	$(CCDEBUG) -D__HMSS_MAIN__ containers/hmss.c  lib/mlrutil.c

lhms2v:
	$(CCDEBUG) -D__LHMS2V_MAIN__ containers/lhms2v.c  lib/mlrutil.c

lhmslv:
	$(CCDEBUG) -D__LHMSLV_MAIN__ containers/lhmslv.c  lib/mlrutil.c containers/slls.c

lhmss:
	$(CCDEBUG) -D__LHMSS_MAIN__ containers/lhmss.c  lib/mlrutil.c

lhmsi:
	$(CCDEBUG) -D__LHMSI_MAIN__ containers/lhmsi.c  lib/mlrutil.c

lhmsv:
	$(CCDEBUG) -D__LHMSV_MAIN__ containers/lhmsv.c  lib/mlrutil.c

lrec:
	$(CCDEBUG) -D__LREC_MAIN__ containers/lrec.c  lib/mlrutil.c

top_keeper:
	$(CCDEBUG) -D__TOP_KEEPER_MAIN__ containers/top_keeper.c lib/mlrutil.c

pk:
	$(CCDEBUG) -D__PERCENTILE_KEEPER_MAIN__ containers/percentile_keeper.c lib/mlrutil.c

ap:
	$(CCDEBUG) -D__AP_MAIN__ cli/argparse.c containers/sllv.c containers/slls.c lib/mlrutil.c

evl:
	$(CCDEBUG) -D__LREC_EVALUATORS_MAIN__ mapping/lrec_evaluators.c mapping/mlr_val.c containers/mlr_dsl_ast.c containers/sllv.c containers/slls.c lib/mlrutil.c lib/mtrand.c containers/lrec.c -lm

lrim:
	$(CCDEBUG) -D__LREC_READER_IN_MEMORY_MAIN__ input/lrec_reader_in_memory.c containers/sllv.c containers/lrec.c containers/slls.c lib/mlrutil.c lib/mlr_globals.c -lm

termcvt: tools/termcvt.c
	$(CCDEBUG) tools/termcvt.c -o termcvt

getl: .always
	$(CC) $(EXPERIMENTAL_READER_SRCS) -o getl

# ================================================================
perfclean profclean:
	@rm -vf gmon.out perf.data perf.data.old
.always:
	@true
