AM_CPPFLAGS=	-I${srcdir}
AM_CFLAGS=	-std=gnu99

noinst_PROGRAMS=	lemon
lemon_SOURCES=		lemon.c

noinst_LTLIBRARIES=	libpdsl.la
# if we don't want to distribute the generated sources, move them to a nodist_libfoo_la_SOURCES variable
libpdsl_la_SOURCES=	put_dsl_lexer.c put_dsl_lexer.h put_dsl_parse.c put_dsl_parse.h put_dsl_wrapper.c put_dsl_wrapper.h
CLEANFILES=		put_dsl_lexer.c put_dsl_lexer.h put_dsl_parse.c put_dsl_parse.h put_dsl_parse.out
EXTRA_DIST=		lempar.c put_dsl_lexer.l put_dsl_parse.y

put_dsl_wrapper.c: put_dsl_lexer.h
put_dsl_parse.h: put_dsl_parse.c
put_dsl_parse.c: put_dsl_parse.y lempar.c lemon
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


put_dsl_lexer.h: put_dsl_lexer.c
put_dsl_lexer.c: put_dsl_lexer.l put_dsl_parse.h
	flex --prefix=put_dsl_lexer_ --outfile=put_dsl_lexer.c --header-file=put_dsl_lexer.h ${srcdir}/put_dsl_lexer.l

distclean-local:
	if [ "${srcdir}" != "." ]; then \
		rm -f lempar.c put_dsl_parse.y; \
	fi
