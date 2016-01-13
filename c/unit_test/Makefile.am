
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

AM_CPPFLAGS=		-I${srcdir}/..
AM_CFLAGS=		-Wall -std=gnu99

all_ldadd=		\
			../cli/libcli.la \
			../containers/libcontainers.la \
			../stream/libstream.la \
			../input/libinput.la \
			../mapping/libmapping.la \
			../output/liboutput.la \
			../lib/libmlr.la \
			../dsls/libdsl.la \
			-lm

# Unit-test mains
test_mlrutil_CFLAGS=              -std=gnu99 -g ${AM_CFLAGS}
test_mlrutil_LDADD=               ${all_ldadd}

test_argparse_CFLAGS=             -std=gnu99 -g ${AM_CFLAGS}
test_argparse_LDADD=              ${all_ldadd}

test_byte_readers_CFLAGS=         -std=gnu99 -g ${AM_CFLAGS}
test_byte_readers_LDADD=          ${all_ldadd}

test_peek_file_reader_CFLAGS=     -std=gnu99 -g ${AM_CFLAGS}
test_peek_file_reader_LDADD=      ${all_ldadd}

test_parse_trie_CFLAGS=           -std=gnu99 -g ${AM_CFLAGS}
test_parse_trie_LDADD=            ${all_ldadd}

test_lrec_CFLAGS=                 -std=gnu99 -g ${AM_CFLAGS}
test_lrec_LDADD=                  ${all_ldadd}

test_multiple_containers_CFLAGS=  -std=gnu99 -g ${AM_CFLAGS}
test_multiple_containers_LDADD=   ${all_ldadd}

test_string_builder_CFLAGS=       -std=gnu99 -g ${AM_CFLAGS}
test_string_builder_LDADD=        ${all_ldadd}

test_lrec_evaluators_CFLAGS=      -std=gnu99 -g ${AM_CFLAGS}
test_lrec_evaluators_LDADD=       ${all_ldadd}

test_join_bucket_keeper_CFLAGS=   -std=gnu99 -g ${AM_CFLAGS}
test_join_bucket_keeper_LDADD=    ${all_ldadd}
