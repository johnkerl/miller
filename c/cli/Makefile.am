noinst_LTLIBRARIES=	libcli.la
libcli_la_SOURCES=	\
			argparse.c \
			argparse.h \
			mlrcli.c \
			mlrcli.h

noinst_PROGRAMS=	ap
ap_CPPFLAGS=		-D__AP_MAIN__ ${AM_CPPFLAGS}
ap_LDADD=		../containers/libcontainers.la ../lib/libmlr.la
ap_SOURCES=		argparse.c

AM_CFLAGS=	-std=gnu99
AM_CPPFLAGS=	-I${srcdir}/../
