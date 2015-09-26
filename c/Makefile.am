SUBDIRS=		lib cli stream input mapping output containers dsls experimental tools . unit_test reg_test

AM_CPPFLAGS=		-I${srcdir}
AM_CFLAGS=		-Wall -Werror -std=gnu99 -O3
bin_PROGRAMS=		mlr
noinst_PROGRAMS=	mlrg \
			mlrp

mlr_SOURCES=		mlrmain.c

mlr_LDADD=		\
			cli/libcli.la \
			containers/libcontainers.la \
			stream/libstream.la \
			input/libinput.la \
			mapping/libmapping.la \
			output/liboutput.la \
			lib/libmlr.la \
			dsls/libfdsl.la \
			dsls/libpdsl.la \
			-lm

# Other executable variants
# Debug version
mlrg_CFLAGS=	-std=gnu99 -g
mlrg_LDADD=	${mlr_LDADD}
mlrg_SOURCES=	${mlr_SOURCES}

# Profile version. Usage:
# * make mlrp
# * mlrp {arguments>
# * gprof mlrp gmon.out > myfile.txt
# Note: works on Linux; not on OSX.
mlrp_CFLAGS=	-std=gnu99 -g -pg
mlrp_LDADD=	${mlr_LDADD}
mlrp_SOURCES=	${mlr_SOURCES}

# ================================================================
# Run this after unit-test expected output has changed, and is verified to be
# OK. (Example: after adding new test cases in test/run.)
regtest-copy:
	cp output/out reg_test/expected

# ================================================================
perfclean profclean:
	@rm -vf gmon.out perf.data perf.data.old
