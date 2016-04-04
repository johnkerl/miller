noinst_LTLIBRARIES=	libmapping.la
libmapping_la_SOURCES=	\
			rval_evaluator.h \
			rval_evaluators.h \
			rval_expr_evaluators.c \
			rval_func_evaluators.c \
			rval_list_evaluators.c \
			mapper.h \
			mapper_bar.c \
			mapper_bootstrap.c \
			mapper_cat.c \
			mapper_check.c \
			mapper_cut.c \
			mapper_decimate.c \
			mapper_filter.c \
			mapper_grep.c \
			mapper_group_like.c \
			mapper_having_fields.c \
			mapper_head.c \
			mapper_histogram.c \
			mapper_join.c \
			mapper_label.c \
			mapper_merge_fields.c \
			mapper_nest.c \
			mapper_put.c \
			mapper_regularize.c \
			mapper_rename.c \
			mapper_reorder.c \
			mapper_repeat.c \
			mapper_reshape.c \
			mapper_sample.c \
			mapper_sec2gmt.c \
			mapper_shuffle.c \
			mapper_sort.c \
			mapper_stats1.c \
			mapper_stats2.c \
			mapper_step.c \
			mapper_tac.c \
			mapper_tail.c \
			mapper_top.c \
			mapper_uniq.c \
			mappers.h \
			mlr_dsl_cst.h \
			mlr_dsl_cst.c \
			stats1_accumulators.c \
			stats1_accumulators.h
libmapping_la_LIBADD=	../lib/libmlr.la ../cli/libcli.la ../input/libinput.la

# TODO: replace the interesting content with unit tests; jettison the rest
# TODO: Does not build
#noinst_PROGRAMS=	evl
#evl_SOURCES=		rval_evaluators.c
#evl_LDADD=		libmapping.la ../containers/libcontainers.la ../lib/libmlr.la
#evl_CPPFLAGS=		-D__LREC_EVALUATORS_MAIN__ ${AM_CPPFLAGS}

AM_CPPFLAGS=	-I${srcdir}/../
AM_CFLAGS=	-std=gnu99
