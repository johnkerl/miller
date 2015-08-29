noinst_LTLIBRARIES=	libcontainers.la
libcontainers_la_SOURCES=	\
			dheap.c \
			dheap.h \
			header_keeper.c \
			header_keeper.h \
			hss.c \
			hss.h \
			join_bucket_keeper.c \
			join_bucket_keeper.h \
			lhms2v.c \
			lhms2v.h \
			lhmsi.c \
			lhmsi.h \
			lhmslv.c \
			lhmslv.h \
			lhmss.c \
			lhmss.h \
			lhmsv.c \
			lhmsv.h \
			lrec.c \
			lrec.h \
			mixutil.c \
			mixutil.h \
			mlr_dsl_ast.c \
			mlr_dsl_ast.h \
			percentile_keeper.c \
			percentile_keeper.h \
			slls.c \
			slls.h \
			sllv.c \
			sllv.h \
			top_keeper.c \
			top_keeper.h
libcontainers_la_LIBADD=	\
			../lib/libmlr.la \
			../mapping/libmapping.la

noinst_PROGRAMS=		test_join_bucket_keeper test_lrec
test_join_bucket_keeper_LDADD=	../input/libinput.la \
				../mapping/libmapping.la \
				../lib/libmlr.la \
				libcontainers.la
test_join_bucket_keeper_CPPFLAGS=	-D__TEST_JOIN_BUCKET_KEEPER_MAIN__ ${AM_CPPFLAGS}
test_lrec_LDADD=		../input/libinput.la \
				../lib/libmlr.la \
				libcontainers.la
test_lrec_CPPFLAGS=	 -D__TEST_LREC_MAIN__ ${AM_CPPFLAGS}

AM_CPPFLAGS=	-I${srcdir}/../
AM_CFLAGS=	-std=gnu99
