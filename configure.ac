AC_PREREQ([2.60])
# Manually increment on updates to https://github.com/johnkerl/miller/releases
AC_INIT([mlr],[5.7.0-dev])
AC_CONFIG_SRCDIR([c/mlrmain.c])
AC_CONFIG_HEADERS([config.h])
AC_CONFIG_AUX_DIR([autotools])
AC_CONFIG_MACRO_DIR([m4])
AM_INIT_AUTOMAKE

AC_CANONICAL_HOST
AC_PROG_CC
m4_ifdef([AM_PROG_AR], [AM_PROG_AR])
AM_PROG_LEX
AC_EXEEXT
LT_INIT

# TODO: better source handling for lemon sources?
# perhaps lemon can be improved to survive being called from the build dir
AC_CONFIG_LINKS([c/parsing/lempar.c:c/parsing/lempar.c])
AC_CONFIG_LINKS([c/parsing/mlr_dsl_parse.y:c/parsing/mlr_dsl_parse.y])

AC_CONFIG_FILES([Makefile
		c/cli/Makefile
		c/containers/Makefile
		c/parsing/Makefile
		c/dsl/Makefile
		c/experimental/Makefile
		c/input/Makefile
		c/lib/Makefile
		c/mapping/Makefile
		c/output/Makefile
		c/stream/Makefile
		c/auxents/Makefile
		c/unit_test/Makefile
		c/reg_test/Makefile
		c/reg_test/expected/Makefile
		c/reg_test/input/Makefile
		c/reg_test/input/rfc-csv/Makefile
		c/reg_test/input/comments/Makefile
		c/Makefile
		doc/Makefile])
AC_OUTPUT
