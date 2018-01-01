noinst_LTLIBRARIES=	libcli.la
libcli_la_SOURCES=	\
			argparse.c \
			argparse.h \
			comment_handling.h \
			json_array_ingest.h \
			mlrcli.c \
			mlrcli.h \
			quoting.h

# TODO: causes circular dependency
#noinst_PROGRAMS=	ap
#ap_CPPFLAGS=		-D__AP_MAIN__ ${AM_CPPFLAGS}
#ap_LDADD=		../containers/libcontainers.la ../lib/libmlr.la
#ap_SOURCES=		argparse.c

AM_CPPFLAGS=	-I${srcdir}/../
AM_CFLAGS=	-std=gnu99
