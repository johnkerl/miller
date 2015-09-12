SUBDIRS=	lib cli stream input mapping output containers dsls experimental tools test

AM_CPPFLAGS=	-I${srcdir}
AM_CFLAGS=	-Wall -Werror -std=gnu99 -O3
bin_PROGRAMS=	mlr mlrg mlrk mlrp
mlr_SOURCES=	mlrmain.c
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

# Other executable variants
# Debug version
mlrg_CFLAGS=	-std=gnu99 -g
mlrg_LDADD=	${mlr_LDADD}
mlrg_SOURCES=	${mlr_SOURCES}

# Memcheck version
mlrk_CFLAGS=	-std=gnu99 -g -DUSE_MCHECK
mlrk_LDADD=	${mlr_LDADD}
mlrk_SOURCES=	${mlr_SOURCES}

# Profile version. Usage:
# * make mlrp
# * mlrp {arguments>
# * gprof mlrp gmon.out > myfile.txt
# Note: works on Linux; not on OSX.
mlrp_CFLAGS=	-std=gnu99 -g -pg
mlrp_LDADD=	${mlr_LDADD}
mlrp_SOURCES=	${mlr_SOURCES}

# ================================================================
# User-make: creates the executable and runs unit & regression tests (without
# valgrind)
# ----> This is the default target for anyone pulling the repo and trying to
#       build it to be able to use it. It just needs flex and the C compiler.
top: mlr check

# Dev-make: updates the tags file, creates the executable, and runs unit &
# regression tests (with valgrind)
# ----> This is the target for a developer to run before a commit.
#       It requires ctags and valgrind in addition to flex and the C compiler.
dev: tags mlr dev-tests

# ================================================================
# User-tests: don't use valgrind: they may not have it, and valgrind is most
# useful for detecting errors at build time. They don't need it to produce an
# executable.

tests: check reg-test

reg-test:
	./test/run

# ----------------------------------------------------------------
# Use valgrind at dev time for additional reassurance.

dev-tests: dev-unit-test reg-test

# Unfortunately --error-exitcode=1 doesn't work well since there are
# valgrind-detected errors in stdlibs. :(
dev-unit-test: test-mlrutil test-byte-readers test-peek-file-reader test-parse-trie test-lrec test-string-builder test-join-bucket-keeper
#	valgrind --leak-check=full --error-exitcode=1 ./a.out
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

# ================================================================
perfclean profclean:
	@rm -vf gmon.out perf.data perf.data.old
