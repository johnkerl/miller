noinst_LTLIBRARIES=	libmapping.la
libmapping_la_SOURCES=	\
			mapper.h \
			mapper_altkv.c \
			mapper_bar.c \
			mapper_bootstrap.c \
			mapper_cat.c \
			mapper_check.c \
			mapper_clean_whitespace.c \
			mapper_count_similar.c \
			mapper_cut.c \
			mapper_decimate.c \
			mapper_fill_down.c \
			mapper_grep.c \
			mapper_group_like.c \
			mapper_having_fields.c \
			mapper_head.c \
			mapper_histogram.c \
			mapper_join.c \
			mapper_label.c \
			mapper_merge_fields.c \
			mapper_most_or_least_frequent.c \
			mapper_nest.c \
			mapper_nothing.c \
			mapper_fraction.c \
			mapper_put_or_filter.c \
			mapper_regularize.c \
			mapper_remove_empty_columns.c \
			mapper_rename.c \
			mapper_reorder.c \
			mapper_repeat.c \
			mapper_reshape.c \
			mapper_sample.c \
			mapper_sec2gmt.c \
			mapper_sec2gmtdate.c \
			mapper_seqgen.c \
			mapper_shuffle.c \
			mapper_skip_trivial_records.c \
			mapper_sort.c \
			mapper_stats1.c \
			mapper_stats2.c \
			mapper_step.c \
			mapper_tac.c \
			mapper_tail.c \
			mapper_tee.c \
			mapper_top.c \
			mapper_uniq.c \
			mapper_unsparsify.c \
			mappers.c \
			mappers.h \
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
