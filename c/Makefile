# ================================================================
# Travis does "export CC=gcc", "export CC=clang" so we can pick those up via
# "make -e" in ../.travis.yml.  Note that "CC?=gcc", without make -e, results
# in CC being expanded to cc on my OSX laptop, which is not OK.  Hence make -e.
CC=gcc
CFLAGS=-std=gnu99
IFLAGS=-I.

WFLAGS=-Wall -Werror
# Worth exploring ... but needs handling for unused parameters in functions which comply with interfaces.
# Best option I'm aware of is to replace "void foo(int bar) {...}" with "void foo(int) {...}" throughout.
# WFLAGS=-Wall -Wextra -Werror
# WFLAGS=-Wall -Wextra -pedantic-errors -Werror
# WFLAGS=-Wall -Wextra -pedantic-errors -Werror=unused-variable

LFLAGS=-lm

# You can do
#
#   make -e INSTALLDIR=/usr/local/bin
#
# pending https://github.com/johnkerl/miller/issues/9
INSTALLDIR=/usr/local/bin

CCOPT=$(CC) $(CFLAGS) $(IFLAGS) $(WFLAGS) -O3
#CCDEBUG=$(CC) -g -O1 $(CFLAGS) $(IFLAGS) $(WFLAGS)
CCDEBUG=$(CC) -g $(CFLAGS) $(IFLAGS) $(WFLAGS)

NON_DSL_SRCS = *.c cli/*.c lib/*.c containers/*.c stream/*.c input/*.c mapping/*.c output/*.c
PDSL_OBJS = ./dsls/put_dsl_parse.o    ./dsls/put_dsl_lexer.o    ./dsls/put_dsl_wrapper.o
FDSL_OBJS = ./dsls/filter_dsl_parse.o ./dsls/filter_dsl_lexer.o ./dsls/filter_dsl_wrapper.o

TEST_BYTE_READERS_SRCS = \
lib/mlrutil.c \
lib/mlr_test_util.c \
lib/mlr_globals.c \
input/string_byte_reader.c \
input/stdio_byte_reader.c \
input/mmap_byte_reader.c \
input/test_byte_readers.c

TEST_PEEK_FILE_READER_SRCS = \
lib/mlrutil.c \
lib/mlr_test_util.c \
lib/mlr_globals.c \
input/string_byte_reader.c \
input/peek_file_reader.c \
input/test_peek_file_reader.c

TEST_LREC_SRCS = lib/mlrutil.c lib/mlr_globals.c lib/string_builder.c \
containers/lrec.c containers/header_keeper.c containers/sllv.c \
containers/slls.c containers/lhmslv.c \
input/file_reader_mmap.c      input/file_reader_stdio.c \
input/lrec_reader_mmap_csvlite.c  input/lrec_reader_stdio_csvlite.c \
input/lrec_reader_mmap_dkvp.c input/lrec_reader_stdio_dkvp.c \
input/lrec_reader_mmap_nidx.c input/lrec_reader_stdio_nidx.c \
input/lrec_reader_mmap_xtab.c input/lrec_reader_stdio_xtab.c \
containers/test_lrec.c

TEST_MAPS_AND_SETS_SRCS = lib/mlrutil.c lib/mlr_globals.c lib/string_builder.c \
containers/lrec.c containers/header_keeper.c containers/sllv.c \
containers/slls.c \
containers/hss.c \
containers/lhmsi.c \
containers/lhms2v.c \
containers/lhmslv.c \
input/file_reader_mmap.c      input/file_reader_stdio.c \
input/lrec_reader_mmap_csvlite.c  input/lrec_reader_stdio_csvlite.c \
input/lrec_reader_mmap_dkvp.c input/lrec_reader_stdio_dkvp.c \
input/lrec_reader_mmap_nidx.c input/lrec_reader_stdio_nidx.c \
input/lrec_reader_mmap_xtab.c input/lrec_reader_stdio_xtab.c \
containers/test_maps_and_sets.c

TEST_JOIN_BUCKET_KEEPER_SRCS = \
lib/mlrutil.c lib/mlr_globals.c lib/string_builder.c \
mapping/context.c \
containers/parse_trie.c \
containers/lrec.c \
containers/sllv.c containers/slls.c containers/lhmslv.c containers/hss.c containers/mixutil.c \
containers/header_keeper.c \
containers/join_bucket_keeper.c \
input/mmap_byte_reader.c \
input/stdio_byte_reader.c \
input/lrec_reader_in_memory.c input/lrec_readers.c \
input/lrec_reader_csv.c \
input/lrec_reader_mmap_csvlite.c input/lrec_reader_stdio_csvlite.c \
input/lrec_reader_mmap_dkvp.c input/lrec_reader_stdio_dkvp.c \
input/lrec_reader_mmap_nidx.c input/lrec_reader_stdio_nidx.c \
input/lrec_reader_mmap_xtab.c input/lrec_reader_stdio_xtab.c \
input/file_reader_mmap.c input/file_reader_stdio.c \
input/peek_file_reader.c \
containers/test_join_bucket_keeper.c

EXPERIMENTAL_READER_SRCS = \
lib/mlrutil.c \
lib/mlr_globals.c \
lib/string_builder.c \
input/file_reader_mmap.c \
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

install: mlr tests
	cp mlr $(INSTALLDIR)
installhome: mlr tests
	cp mlr $(HOME)/bin

# ================================================================
tags: .always
	ctags -R .

mlr: .always dsls
	$(CCOPT) $(NON_DSL_SRCS) $(PDSL_OBJS) $(FDSL_OBJS) $(LFLAGS) -o mlr
mlr.static: .always dsls
	$(CCOPT) -static $(NON_DSL_SRCS) $(PDSL_OBJS) $(FDSL_OBJS) $(LFLAGS) -o mlr.static

dsls: .always
	make -C dsls put_dsl_parse.o
	make -C dsls put_dsl_lexer.o
	make -C dsls put_dsl_wrapper.o
	make -C dsls filter_dsl_parse.o
	make -C dsls filter_dsl_lexer.o
	make -C dsls filter_dsl_wrapper.o

# ----------------------------------------------------------------
# Other executable variants

# Debug version
mlrg: .always dsls
	$(CCDEBUG) $(NON_DSL_SRCS) $(PDSL_OBJS) $(FDSL_OBJS) $(LFLAGS) -o mlrg

# Memcheck version
mlrk: .always dsls
	$(CCDEBUG) -DUSE_MCHECK $(NON_DSL_SRCS) $(PDSL_OBJS) $(FDSL_OBJS) $(LFLAGS) -o mlrk

# Profile version. Usage:
# * make mlrp
# * mlrp {arguments>
# * gprof mlrp gmon.out > myfile.txt
# Note: works on Linux; not on OSX.
mlrp: .always dsls
	$(CCDEBUG) -g -pg $(NON_DSL_SRCS) $(PDSL_OBJS) $(FDSL_OBJS) $(LFLAGS) -o mlrp

# ================================================================
# User-tests: don't use valgrind: they may not have it, and valgrind is most
# useful for detecting errors at build time. They don't need it to produce an
# executable.

tests: unit-test reg-test

unit-test: test-mlrutil test-byte-readers test-peek-file-reader test-parse-trie test-lrec test-maps-and-sets test-string-builder test-join-bucket-keeper
	./test-mlrutil
	./test-byte-readers
	./test-peek-file-reader
	./test-parse-trie
	./test-lrec
	./test-maps-and-sets
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
dev-unit-test: test-mlrutil test-byte-readers test-peek-file-reader test-parse-trie test-lrec test-maps-and-sets test-string-builder test-join-bucket-keeper
	#valgrind --leak-check=full --error-exitcode=1 ./a.out
	valgrind --leak-check=full ./test-mlrutil
	valgrind --leak-check=full ./test-byte-readers
	valgrind --leak-check=full ./test-Peek-file-reader
	valgrind --leak-check=full ./test-parse-trie
	valgrind --leak-check=full ./test-lrec
	valgrind --leak-check=full ./test-maps-and-sets
	valgrind --leak-check=full ./test-string-builder
	valgrind --leak-check=full ./test-join-bucket-keeper
	@echo
	@echo DONE

# Run this after unit-test expected output has changed, and is verified to be
# OK. (Example: after adding new test cases in test/run.)
regtest-copy:
	cp test/output/out test/expected

# ----------------------------------------------------------------
# Unit-test executables

test-byte-readers: .always
	$(CCDEBUG) -D__TEST_BYTE_READERS_MAIN__ $(TEST_BYTE_READERS_SRCS) -o test-byte-readers

test-peek-file-reader: .always
	$(CCDEBUG) -D__TEST_PEEK_FILE_READER_MAIN__ $(TEST_PEEK_FILE_READER_SRCS) -o test-peek-file-reader

test-lrec: .always
	$(CCDEBUG) -D__TEST_LREC_MAIN__ $(TEST_LREC_SRCS) -o test-lrec

test-maps-and-sets: .always
	$(CCDEBUG) -D__TEST_MAPS_AND_SETS_MAIN__ $(TEST_MAPS_AND_SETS_SRCS) -o test-maps-and-sets

test-mlrutil: .always
	$(CCDEBUG) -D__TEST_MLRUTIL_MAIN__ lib/mlrutil.c lib/test_mlrutil.c -o test-mlrutil

test-string-builder: .always
	$(CCDEBUG) -D__TEST_STRING_BUILDER_MAIN__ lib/mlrutil.c lib/mlr_globals.c lib/string_builder.c lib/test_string_builder.c -o test-string-builder

test-parse-trie: .always
	$(CCDEBUG) -D__TEST_PARSE_TRIE_MAIN__ lib/mlrutil.c lib/mlr_globals.c containers/parse_trie.c containers/test_parse_trie.c -o test-parse-trie

test-join-bucket-keeper: .always
	$(CCDEBUG) -D__TEST_JOIN_BUCKET_KEEPER_MAIN__ $(TEST_JOIN_BUCKET_KEEPER_SRCS) -o test-join-bucket-keeper

# ----------------------------------------------------------------
# Standalone mains
# TODO: replace the interesting content with unit tests; jettison the rest

ap:
	$(CCDEBUG) -D__AP_MAIN__ cli/argparse.c containers/sllv.c containers/slls.c lib/mlrutil.c

evl:
	$(CCDEBUG) -D__LREC_EVALUATORS_MAIN__ mapping/lrec_evaluators.c mapping/mlr_val.c containers/mlr_dsl_ast.c containers/sllv.c containers/slls.c lib/mlrutil.c lib/mtrand.c containers/lrec.c -lm

lrim:
	$(CCDEBUG) -D__LREC_READER_IN_MEMORY_MAIN__ input/lrec_reader_in_memory.c containers/sllv.c containers/lrec.c containers/slls.c lib/mlrutil.c lib/mlr_globals.c -lm

termcvt: tools/termcvt.c
	$(CCDEBUG) tools/termcvt.c -o termcvt

getl: .always
	$(CCDEBUG) $(EXPERIMENTAL_READER_SRCS) -o getl

# ================================================================
clean:
	@rm -vf mlr mlrd mlrg mlrp tester
	@make -C dsls clean
perfclean profclean:
	@rm -vf gmon.out perf.data perf.data.old
.always:
	@true
