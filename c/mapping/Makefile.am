noinst_LTLIBRARIES=	libmapping.la
libmapping_la_SOURCES=	\
			context.c \
			context.h \
			lrec_evaluator.h \
			lrec_evaluators.c \
			lrec_evaluators.h \
			mapper.h \
			mapper_cat.c \
			mapper_check.c \
			mapper_cut.c \
			mapper_filter.c \
			mapper_group_like.c \
			mapper_having_fields.c \
			mapper_head.c \
			mapper_histogram.c \
			mapper_join.c \
			mapper_label.c \
			mapper_put.c \
			mapper_regularize.c \
			mapper_rename.c \
			mapper_reorder.c \
			mapper_sort.c \
			mapper_stats1.c \
			mapper_stats2.c \
			mapper_step.c \
			mapper_tac.c \
			mapper_tail.c \
			mapper_top.c \
			mapper_uniq.c \
			mappers.h \
			mlr_val.c \
			mlr_val.h
libmapping_la_LIBADD=	../lib/libmlr.la ../cli/libcli.la ../input/libinput.la

# TODO: replace the interesting content with unit tests; jettison the rest
# TODO: Does not build
#noinst_PROGRAMS=	evl
#evl_SOURCES=		lrec_evaluators.c
#evl_LDADD=		libmapping.la ../containers/libcontainers.la ../lib/libmlr.la
#evl_CPPFLAGS=		-D__LREC_EVALUATORS_MAIN__ ${AM_CPPFLAGS}

AM_CPPFLAGS=	-I${srcdir}/../
AM_CFLAGS=	-std=gnu99
