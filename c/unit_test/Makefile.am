
TESTS= $(check_PROGRAMS)
check_PROGRAMS = 	test_mlrutil \
			test_mlrregex \
			test_argparse \
			test_byte_readers \
			test_line_readers \
			test_peek_file_reader \
			test_parse_trie \
			test_lrec \
			test_mlhmmv \
			test_multiple_containers \
			test_string_builder \
			test_rval_evaluators \
			test_join_bucket_keeper

AM_CPPFLAGS=		-I${srcdir}/..
AM_CFLAGS=		-Wall -std=gnu99

all_ldadd=		\
			../lib/libmlr.la \
			../containers/libcontainers.la \
			../parsing/libdsl.la \
			../cli/libcli.la \
			../input/libinput.la \
			../dsl/libdsl.la \
			../mapping/libmapping.la \
			../output/liboutput.la \
			../stream/libstream.la

all_ldflags=		-lm

# Unit-test mains
test_mlrutil_CFLAGS=              -std=gnu99 -g ${AM_CFLAGS}
test_mlrutil_LDFLAGS=             ${all_ldflags}
test_mlrutil_LDADD=               ${all_ldadd}

test_mlrregex_CFLAGS=             -std=gnu99 -g ${AM_CFLAGS}
test_mlrregex_LDFLAGS=            ${all_ldflags}
test_mlrregex_LDADD=              ${all_ldadd}

test_argparse_CFLAGS=             -std=gnu99 -g ${AM_CFLAGS}
test_argparse_LDFLAGS=            ${all_ldflags}
test_argparse_LDADD=              ${all_ldadd}

test_byte_readers_CFLAGS=         -std=gnu99 -g ${AM_CFLAGS}
test_byte_readers_LDFLAGS=        ${all_ldflags}
test_byte_readers_LDADD=          ${all_ldadd}

test_line_readers_CFLAGS=         -std=gnu99 -g ${AM_CFLAGS}
test_line_readers_LDFLAGS=        ${all_ldflags}
test_line_readers_LDADD=          ${all_ldadd}

test_peek_file_reader_CFLAGS=     -std=gnu99 -g ${AM_CFLAGS}
test_peek_file_reader_LDFLAGS=    ${all_ldflags}
test_peek_file_reader_LDADD=      ${all_ldadd}

test_parse_trie_CFLAGS=           -std=gnu99 -g ${AM_CFLAGS}
test_parse_trie_LDFLAGS=          ${all_ldflags}
test_parse_trie_LDADD=            ${all_ldadd}

test_lrec_CFLAGS=                 -std=gnu99 -g ${AM_CFLAGS}
test_lrec_LDFLAGS=                ${all_ldflags}
test_lrec_LDADD=                  ${all_ldadd}

test_mlhmmv_CFLAGS=               -std=gnu99 -g ${AM_CFLAGS}
test_mlhmmv_LDFLAGS=              ${all_ldflags}
test_mlhmmv_LDADD=                ${all_ldadd}

test_multiple_containers_CFLAGS=  -std=gnu99 -g ${AM_CFLAGS}
test_multiple_containers_LDFLAGS= ${all_ldflags}
test_multiple_containers_LDADD=   ${all_ldadd}

test_string_builder_CFLAGS=       -std=gnu99 -g ${AM_CFLAGS}
test_string_builder_LDFLAGS=      ${all_ldflags}
test_string_builder_LDADD=        ${all_ldadd}

test_rval_evaluators_CFLAGS=      -std=gnu99 -g ${AM_CFLAGS}
test_rval_evaluators_LDFLAGS=     ${all_ldflags}
test_rval_evaluators_LDADD=       ${all_ldadd}

test_join_bucket_keeper_CFLAGS=   -std=gnu99 -g ${AM_CFLAGS}
test_join_bucket_keeper_LDFLAGS=  ${all_ldflags}
test_join_bucket_keeper_LDADD=    ${all_ldadd}
