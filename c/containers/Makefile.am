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
			parse_trie.c \
			parse_trie.h \
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

TESTS=			$(check_PROGRAMS)
check_PROGRAMS=		test_join_bucket_keeper test_lrec test_parse_trie
# TODO: replace the interesting content with unit tests; jettison the rest
noinst_PROGRAMS=	\
			dheap \
			lhms2v \
			lhmsi \
			lhmslv \
			lhmss \
			lhmsv \
			pk \
			slls \
			sllv \
			top_keeper
dheap_CPPFLAGS=		-D__DHEAP_MAIN__ ${AM_CPPFLAGS}
dheap_LDADD=		../lib/libmlr.la
lhms2v_CPPFLAGS=	-D__LHMS2V_MAIN__ ${AM_CPPFLAGS}
lhms2v_LDADD=		../lib/libmlr.la
lhmsi_CPPFLAGS=		-D__LHMSI_MAIN__ ${AM_CPPFLAGS}
lhmsi_LDADD=		../lib/libmlr.la
lhmslv_CPPFLAGS=	-D__LHMSLV_MAIN__ ${AM_CPPFLAGS}
lhmslv_LDADD=		../lib/libmlr.la
lhmslv_SOURCES=		lhmslv.c slls.c
lhmss_CPPFLAGS=		-D__LHMSS_MAIN__ ${AM_CPPFLAGS}
lhmss_LDADD=		../lib/libmlr.la
lhmsv_CPPFLAGS=		-D__LHMSV_MAIN__ ${AM_CPPFLAGS}
lhmsv_LDADD=		../lib/libmlr.la
pk_CPPFLAGS=		-D__PERCENTILE_KEEPER_MAIN__ ${AM_CPPFLAGS}
pk_LDADD=		../lib/libmlr.la
pk_SOURCES=		percentile_keeper.c
slls_CPPFLAGS=		-D__SLLS_MAIN__ ${AM_CPPFLAGS}
slls_LDADD=		../lib/libmlr.la
sllv_CPPFLAGS=		-D__SLLV_MAIN__ ${AM_CPPFLAGS}
sllv_LDADD=		../lib/libmlr.la
top_keeper_CPPFLAGS=	-D__TOP_KEEPER_MAIN__ ${AM_CPPFLAGS}
top_keeper_LDADD=	../lib/libmlr.la

test_join_bucket_keeper_CPPFLAGS=	-D__TEST_JOIN_BUCKET_KEEPER_MAIN__ ${AM_CPPFLAGS}
test_join_bucket_keeper_LDADD=		../input/libinput.la \
					../mapping/libmapping.la \
					../lib/libmlr.la \
					libcontainers.la

test_lrec_CPPFLAGS=		-D__TEST_LREC_MAIN__ ${AM_CPPFLAGS}
test_lrec_LDADD=		../input/libinput.la \
				../lib/libmlr.la \
				libcontainers.la

test_parse_trie_CPPFLAGS=	-D__TEST_PARSE_TRIE_MAIN__ ${AM_CPPFLAGS}
test_parse_trie_LDADD=		libcontainers.la

AM_CPPFLAGS=	-I${srcdir}/../
AM_CFLAGS=	-std=gnu99
