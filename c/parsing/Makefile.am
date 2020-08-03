AM_CPPFLAGS=	-I${srcdir}
AM_CFLAGS=	-std=gnu99

noinst_PROGRAMS=	lemon
lemon_SOURCES=		lemon_action.c \
			lemon_action.h \
			lemon_assert.c \
			lemon_assert.h \
			lemon_config_list.c \
			lemon_config_list.h \
			lemon_dims.h \
			lemon_error.c \
			lemon_error.h \
			lemon_fsm.c \
			lemon_fsm.h \
			lemon_main.c \
			lemon_memory.c \
			lemon_memory.h \
			lemon_msort.c \
			lemon_msort.h \
			lemon_option.c \
			lemon_option.h \
			lemon_parse.c \
			lemon_parse.h \
			lemon_plink.c \
			lemon_plink.h \
			lemon_report.c \
			lemon_report.h \
			lemon_set.c \
			lemon_set.h \
			lemon_state_table.c \
			lemon_state_table.h \
			lemon_string.c \
			lemon_string.h \
			lemon_structs.h \
			lemon_symbol.c \
			lemon_symbol.h

noinst_LTLIBRARIES=	libdsl.la
# if we don't want to distribute the generated sources, move them to a nodist_libfoo_la_SOURCES variable
libdsl_la_SOURCES=	mlr_dsl_lexer.c mlr_dsl_lexer.h mlr_dsl_parse.c mlr_dsl_parse.h mlr_dsl_wrapper.c mlr_dsl_wrapper.h
CLEANFILES=		mlr_dsl_lexer.c mlr_dsl_lexer.h mlr_dsl_parse.c mlr_dsl_parse.h mlr_dsl_parse.out
EXTRA_DIST=		lempar.c mlr_dsl_lexer.l mlr_dsl_parse.y

mlr_dsl_wrapper.c: mlr_dsl_lexer.h
mlr_dsl_parse.h: mlr_dsl_parse.c
mlr_dsl_parse.c: mlr_dsl_parse.y lempar.c lemon${EXEEXT}
	./lemon -s mlr_dsl_parse.y
# lemon does not re-generate the file if it exists and is the same
# this causes problems with the dependency above
	touch mlr_dsl_parse.h
	mv mlr_dsl_parse.c mlr_dsl_parse.c.tmp
	sed \
			-e 's/ParseTrace/mlr_dsl_ParseTrace/g' \
			-e 's/ParseTokenName/mlr_dsl_ParseTokenName/g' \
			-e 's/lemon_parser_alloc/mlr_dsl_lemon_parser_alloc/g' \
			-e 's/lemon_parser_free/mlr_dsl_lemon_parser_free/g' \
			-e 's/lemon_parser_parse_token/mlr_dsl_lemon_parser_parse_token/g' \
			-e 's/yy_destructor/mlr_dsl_yy_destructor/g' \
		mlr_dsl_parse.c.tmp > mlr_dsl_parse.c
	rm -f mlr_dsl_parse.c.tmp


mlr_dsl_lexer.h: mlr_dsl_lexer.c
mlr_dsl_lexer.c: mlr_dsl_lexer.l mlr_dsl_parse.h
	flex --prefix=mlr_dsl_lexer_ --outfile=mlr_dsl_lexer.c --header-file=mlr_dsl_lexer.h ${srcdir}/mlr_dsl_lexer.l

distclean-local:
	if [ "${srcdir}" != "." ]; then \
		rm -f lempar.c mlr_dsl_parse.y; \
	fi
