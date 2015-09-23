noinst_LTLIBRARIES=	libmlr.la
libmlr_la_SOURCES=	minunit.h \
			mlr_globals.c \
			mlr_globals.h \
			mlrmath.c \
			mlrmath.h \
			mlrstat.c \
			mlrstat.h \
			mlrutil.c \
			mlrutil.h \
			context.c \
			context.h \
			mtrand.c \
			mtrand.h \
			string_builder.c \
			string_builder.h \
			mlr_test_util.c \
			mlr_test_util.h

TESTS=		$(check_PROGRAMS)
check_PROGRAMS=	test_mlrutil test_string_builder

test_mlrutil_LDADD=		libmlr.la
test_mlrutil_CPPFLAGS=		-D__TEST_MLRUTIL_MAIN__ ${AM_CPPFLAGS}

test_string_builder_LDADD=	libmlr.la
test_string_builder_CPPFLAGS=	-D__TEST_STRING_BUILDER_MAIN__ ${AM_CPPFLAGS}


AM_CPPFLAGS=	-I${srcdir}/../
AM_CFLAGS=	-std=gnu99
