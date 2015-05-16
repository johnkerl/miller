CCOMP=gcc
CFLAGS=-std=gnu99
IFLAGS=-I.
WFLAGS=-Wall -Werror
LFLAGS=-lm

CC=$(CCOMP) $(CFLAGS) $(IFLAGS) $(WFLAGS) -O3
#CCDEBUG=$(CCOMP) -g -O1 $(CFLAGS) $(IFLAGS) $(WFLAGS)
CCDEBUG=$(CCOMP) -g $(CFLAGS) $(IFLAGS) $(WFLAGS)

NON_DSL_SRCS = *.c cli/*.c lib/*.c containers/*.c stream/*.c input/*.c mapping/*.c output/*.c
PDSL_OBJS = ./dsls/put_dsl_parse.o    ./dsls/put_dsl_lexer.o    ./dsls/put_dsl_wrapper.o
FDSL_OBJS = ./dsls/filter_dsl_parse.o ./dsls/filter_dsl_lexer.o ./dsls/filter_dsl_wrapper.o

# ----------------------------------------------------------------
# mlr  for default use.
# mlrg for valgrind.

top: tags mlr test

install: top
	cp mlr $(HOME)/bin/

all: tags mlr mlrd mlrg mlrp catc catc0 catm catd catgo catrust cutd cutgo

tags: .always
	ctags -R .

mlr: .always tags dsls
	$(CC) $(NON_DSL_SRCS) $(PDSL_OBJS) $(FDSL_OBJS) $(LFLAGS) -o mlr

# Memcheck version
mlrk: .always tags dsls
	$(CC) -DUSE_MCHECK $(NON_DSL_SRCS) $(PDSL_OBJS) $(FDSL_OBJS) $(LFLAGS) -o mlrk

# Debug version
mlrg: .always tags dsls
	$(CCDEBUG) $(NON_DSL_SRCS) $(PDSL_OBJS) $(FDSL_OBJS) $(LFLAGS) -o mlrg

# Profile version
mlrp: .always tags dsls
	$(CC) -g -pg $(NON_DSL_SRCS) $(PDSL_OBJS) $(FDSL_OBJS) $(LFLAGS) -o mlrp

# ----------------------------------------------------------------
dsls: .always
	make -C dsls put_dsl_parse.o
	make -C dsls put_dsl_lexer.o
	make -C dsls put_dsl_wrapper.o
	make -C dsls filter_dsl_parse.o
	make -C dsls filter_dsl_lexer.o
	make -C dsls filter_dsl_wrapper.o

# ----------------------------------------------------------------
test: unit-test reg-test

# Unfortunately --error-exitcode=1 doesn't work well since there are
# valgrind-detected errors in stdlibs. :(

unit-test: test-mlrutil test-lrec
	#valgrind --leak-check=full --error-exitcode=1 ./a.out
	valgrind --leak-check=full ./test-mlrutil
	valgrind --leak-check=full ./test-lrec
	@echo
	@echo DONE

reg-test:
	./test/run

# Run this after unit-test expected output has changed, and is verified to be
# OK. (Example: after adding new test cases in test/run.)
regtest-copy:
	cp test/output/out test/expected

# ----------------------------------------------------------------
# Unit-test executables
test-lrec: .always
	$(CCDEBUG) -D__TEST_LREC_MAIN__ lib/mlrutil.c containers/lrec.c containers/hdr_keeper.c containers/sllv.c containers/slls.c containers/lrec_parsers.c containers/test_lrec.c -o test-lrec

test-mlrutil: .always
	$(CCDEBUG) -D__TEST_MLRUTIL_MAIN__ lib/mlrutil.c lib/test_mlrutil.c -o test-mlrutil

# ----------------------------------------------------------------
# Standalone mains
dheap:
	$(CCDEBUG) -D__DHEAP_MAIN__ containers/dheap.c lib/mlrutil.c

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

top_keeper: $(CCDEBUG) -D__TOP_KEEPER_MAIN__ containers/top_keeper.c lib/mlrutil.c

pk:
	$(CCDEBUG) -D__PERCENTILE_KEEPER_MAIN__ containers/percentile_keeper.c lib/mlrutil.c

ap:
	$(CCDEBUG) -D__AP_MAIN__ cli/argparse.c containers/sllv.c containers/slls.c lib/mlrutil.c

evl:
	$(CCDEBUG) -D__LREC_EVALUATORS_MAIN__ mapping/lrec_evaluators.c mapping/mlr_val.c containers/mlr_dsl_ast.c containers/sllv.c containers/slls.c lib/mlrutil.c lib/mtrand.c containers/lrec.c -lm

# ----------------------------------------------------------------
# xxx to do: move these to ../language-comparisons.

catc: tools/catc.c
	$(CC) tools/catc.c -o catc
catc0: tools/catc0.c
	$(CC) tools/catc0.c -o catc0
catm: tools/catm.c
	$(CC) tools/catm.c -o catm
catd: catd.d
	dmd -O catd.d
catgo: catgo.go
	go build catgo.go
catrust: catrust.rs
	rustc catrust.rs

cutd: cutd.d
	dmd -O cutd.d
cutgo: cutgo.go
	go build cutgo.go

# ----------------------------------------------------------------
clean:
	@rm -vf mlr mlrd mlrg mlrp tester catc catc0 catm catd catgo catrust cutd cutgo mlr-unit-tester
	@make -C dsls clean
perfclean profclean:
	@rm -vf gmon.out perf.data perf.data.old
tar:
	tar zcf mlrc.tgz *.c *.h Makefile
.always:
	@true
