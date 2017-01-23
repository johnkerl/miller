noinst_LTLIBRARIES=	libdsl.la
libdsl_la_SOURCES=	\
			context_flags.h \
			function_manager.c \
			function_manager.h \
			mlr_dsl_ast.c \
			mlr_dsl_ast.h \
			mlr_dsl_blocked_ast.c \
			mlr_dsl_blocked_ast.h \
			mlr_dsl_cst.c \
			mlr_dsl_cst.h \
			mlr_dsl_cst_condish_statements.c \
			mlr_dsl_cst_for_map_statements.c \
			mlr_dsl_cst_for_srec_statements.c \
			mlr_dsl_cst_func_subr.c \
			mlr_dsl_cst_keywords.c \
			mlr_dsl_cst_loop_control_statements.c \
			mlr_dsl_cst_map_assignment_statements.c \
			mlr_dsl_cst_output_statements.c \
			mlr_dsl_cst_return_statements.c \
			mlr_dsl_cst_scalar_assignment_statements.c \
			mlr_dsl_cst_statements.c \
			mlr_dsl_cst_triple_for_statements.c \
			mlr_dsl_cst_unset_statements.c \
			mlr_dsl_stack_allocate.c \
			rval_evaluator.h \
			rxval_evaluator.h \
			rval_evaluators.h \
			rval_expr_evaluators.c \
			rval_func_evaluators.c \
			rval_list_evaluators.c \
			rxval_expr_evaluators.c \
			rxval_func_evaluators.c \
			type_inference.h
libdsl_la_LIBADD=	../lib/libmlr.la ../cli/libcli.la ../input/libinput.la

# TODO: replace the interesting content with unit tests; jettison the rest
# TODO: Does not build
#noinst_PROGRAMS=	evl
#evl_SOURCES=		rval_evaluators.c
#evl_LDADD=		libdsl.la ../containers/libcontainers.la ../lib/libmlr.la
#evl_CPPFLAGS=		-D__LREC_EVALUATORS_MAIN__ ${AM_CPPFLAGS}

AM_CPPFLAGS=	-I${srcdir}/../
AM_CFLAGS=	-std=gnu99
