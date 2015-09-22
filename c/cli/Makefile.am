noinst_LTLIBRARIES=	libcli.la
libcli_la_SOURCES=	\
			argparse.c \
			argparse.h \
			mlrcli.c \
			mlrcli.h
libcli_la_LIBADD=	\
			../lib/libmlr.la
			../containers/libcontainers.la

# TODO: causes circular dependency
#noinst_PROGRAMS=	ap
#ap_CPPFLAGS=		-D__AP_MAIN__ ${AM_CPPFLAGS}
#ap_LDADD=		../containers/libcontainers.la ../lib/libmlr.la
#ap_SOURCES=		argparse.c

TESTS=	$(check_PROGRAMS)
check_PROGRAMS=	test_argparse

test_argparse_CPPFLAGS=	-D__TEST_ARGPARSE_MAIN__ ${AM_CPPFLAGS}
test_argparse_LDADD=	../lib/libmlr.la \
			../containers/libcontainers.la \
			libcli.la

AM_CPPFLAGS=	-I${srcdir}/../
AM_CFLAGS=	-std=gnu99
