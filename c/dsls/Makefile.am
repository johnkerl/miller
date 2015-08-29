AM_CPPFLAGS=	-I${srcdir}
AM_CFLAGS=	-std=gnu99

noinst_PROGRAMS=	lemon pdm fdm
lemon_SOURCES=		lemon.c

noinst_LTLIBRARIES=	libfdsl.la libpdsl.la
libfdsl_la_SOURCES=	filter_dsl_parse.c filter_dsl_lexer.c filter_dsl_wrapper.c
libpdsl_la_SOURCES=	put_dsl_parse.c put_dsl_lexer.c put_dsl_wrapper.c
BUILT_SOURCES=		put_dsl_parse.c filter_dsl_parse.c
CLEANFILES=		${BUILT_SOURCES} put_dsl_parse.out filter_dsl_parse.out
EXTRA_DIST=		\
			filter_dsl_lexer.l \
			filter_dsl_parse.y \
			put_dsl_lexer.l \
			put_dsl_parse.y

fdm_SOURCES=		filter_dsl_parse.c filter_dsl_lexer.c filter_dsl_wrapper.c
fdm_LDADD=		../containers/libcontainers.la ../lib/libmlr.la libfdsl.la -lm
fdm_CFLAGS=		-D __FILTER_DSL_MAIN__ ${AM_CFLAGS}
pdm_SOURCES=		put_dsl_parse.c put_dsl_lexer.c put_dsl_wrapper.c
pdm_LDADD=		../containers/libcontainers.la ../lib/libmlr.la libpdsl.la -lm
pdm_CFLAGS=		-D __PUT_DSL_MAIN__ ${AM_CFLAGS}

put_dsl_parse.c put_dsl_parse.h: put_dsl_parse.y lemon
	./lemon put_dsl_parse.y
# lemon does not re-generate the file if it exists and is the same
# this causes problems with the dependency above
	touch put_dsl_parse.h
	mv put_dsl_parse.c put_dsl_parse.c.tmp
	sed \
			-e 's/ParseTrace/put_dsl_ParseTrace/g' \
			-e 's/ParseTokenName/put_dsl_ParseTokenName/g' \
			-e 's/lemon_parser_alloc/put_dsl_lemon_parser_alloc/g' \
			-e 's/lemon_parser_free/put_dsl_lemon_parser_free/g' \
			-e 's/lemon_parser_parse_token/put_dsl_lemon_parser_parse_token/g' \
			-e 's/yy_destructor/put_dsl_yy_destructor/g' \
		put_dsl_parse.c.tmp > put_dsl_parse.c
	rm -f put_dsl_parse.c.tmp


put_dsl_lexer.c put_dsl_lexer.h: put_dsl_lexer.l
	flex --prefix=put_dsl_lexer_ --outfile=put_dsl_lexer.c --header-file=put_dsl_lexer.h put_dsl_lexer.l

filter_dsl_parse.c filter_dsl_parse.h: filter_dsl_parse.y lemon
	./lemon filter_dsl_parse.y
# lemon does not re-generate the file if it exists and is the same
# this causes problems with the dependency above
	touch filter_dsl_parse.h
	mv filter_dsl_parse.c filter_dsl_parse.c.tmp
	sed \
			-e 's/ParseTrace/filter_dsl_ParseTrace/g' \
			-e 's/ParseTokenName/filter_dsl_ParseTokenName/g' \
			-e 's/lemon_parser_alloc/filter_dsl_lemon_parser_alloc/g' \
			-e 's/lemon_parser_free/filter_dsl_lemon_parser_free/g' \
			-e 's/lemon_parser_parse_token/filter_dsl_lemon_parser_parse_token/g' \
			-e 's/yy_destructor/filter_dsl_yy_destructor/g' \
		filter_dsl_parse.c.tmp > filter_dsl_parse.c
	rm -f filter_dsl_parse.c.tmp

filter_dsl_lexer.c filter_dsl_lexer.h: filter_dsl_lexer.l
	flex --prefix=filter_dsl_lexer_ --outfile=filter_dsl_lexer.c --header-file=filter_dsl_lexer.h filter_dsl_lexer.l
