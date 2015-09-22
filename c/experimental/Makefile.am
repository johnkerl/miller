# TODO: replace the interesting content with unit tests; jettison the rest
noinst_PROGRAMS=	getl
AM_CFLAGS=		-std=gnu99
AM_CPPFLAGS=		-I${srcdir}/../

getl_SOURCES=	getlines.c
getl_LDADD=	../lib/libmlr.la ../input/libinput.la ../containers/libcontainers.la
