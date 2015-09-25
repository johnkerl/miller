SUBDIRS=		lib cli stream input mapping output containers dsls experimental tools . reg_test

AM_CPPFLAGS=		-I${srcdir}
AM_CFLAGS=		-Wall -Werror -std=gnu99 -O3
bin_PROGRAMS=		mlr
noinst_PROGRAMS=	mlrg \
			mlrp

TESTS= $(check_PROGRAMS)
check_PROGRAMS = 	test_mlrutil \
			test_argparse \
			test_byte_readers \
			test_peek_file_reader \
			test_parse_trie \
			test_lrec \
			test_multiple_containers \
			test_string_builder \
			test_lrec_evaluators \
			test_join_bucket_keeper

mlr_SOURCES=		mlrmain.c

mlr_LDADD=		\
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

# Profile version. Usage:
# * make mlrp
# * mlrp {arguments>
# * gprof mlrp gmon.out > myfile.txt
# Note: works on Linux; not on OSX.
mlrp_CFLAGS=	-std=gnu99 -g -pg
mlrp_LDADD=	${mlr_LDADD}
mlrp_SOURCES=	${mlr_SOURCES}

# Unit-test mains
test_mlrutil_CFLAGS=              -std=gnu99 -g
test_mlrutil_LDADD=               ${mlr_LDADD}
test_mlrutil_SOURCES=             unit_test/test_mlrutil.c

test_argparse_CFLAGS=             -std=gnu99 -g
test_argparse_LDADD=              ${mlr_LDADD}
test_argparse_SOURCES=            unit_test/test_argparse.c

test_byte_readers_CFLAGS=         -std=gnu99 -g
test_byte_readers_LDADD=          ${mlr_LDADD}
test_byte_readers_SOURCES=        unit_test/test_byte_readers.c

test_peek_file_reader_CFLAGS=     -std=gnu99 -g
test_peek_file_reader_LDADD=      ${mlr_LDADD}
test_peek_file_reader_SOURCES=    unit_test/test_peek_file_reader.c

test_parse_trie_CFLAGS=           -std=gnu99 -g
test_parse_trie_LDADD=            ${mlr_LDADD}
test_parse_trie_SOURCES=          unit_test/test_parse_trie.c

test_lrec_CFLAGS=                 -std=gnu99 -g
test_lrec_LDADD=                  ${mlr_LDADD}
test_lrec_SOURCES=                unit_test/test_lrec.c

test_multiple_containers_CFLAGS=  -std=gnu99 -g
test_multiple_containers_LDADD=   ${mlr_LDADD}
test_multiple_containers_SOURCES= unit_test/test_multiple_containers.c

test_string_builder_CFLAGS=       -std=gnu99 -g
test_string_builder_LDADD=        ${mlr_LDADD}
test_string_builder_SOURCES=      unit_test/test_string_builder.c

test_lrec_evaluators_CFLAGS=      -std=gnu99 -g
test_lrec_evaluators_LDADD=       ${mlr_LDADD}
test_lrec_evaluators_SOURCES=     unit_test/test_lrec_evaluators.c

test_join_bucket_keeper_CFLAGS=   -std=gnu99 -g
test_join_bucket_keeper_LDADD=    ${mlr_LDADD}
test_join_bucket_keeper_SOURCES=  unit_test/test_join_bucket_keeper.c

# ================================================================
# Run this after unit-test expected output has changed, and is verified to be
# OK. (Example: after adding new test cases in test/run.)
regtest-copy:
	cp output/out test/expected

# ================================================================
perfclean profclean:
	@rm -vf gmon.out perf.data perf.data.old
