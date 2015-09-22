noinst_LTLIBRARIES=	libinput.la
libinput_la_SOURCES=	\
			byte_reader.h \
			byte_readers.h \
			file_reader_mmap.c \
			file_reader_mmap.h \
			file_reader_stdio.c \
			file_reader_stdio.h \
			line_readers.c \
			line_readers.h \
			lrec_reader.h \
			lrec_reader_csv.c \
			lrec_reader_in_memory.c \
			lrec_reader_mmap_csvlite.c \
			lrec_reader_mmap_dkvp.c \
			lrec_reader_mmap_nidx.c \
			lrec_reader_mmap_xtab.c \
			lrec_reader_stdio_csvlite.c \
			lrec_reader_stdio_dkvp.c \
			lrec_reader_stdio_nidx.c \
			lrec_reader_stdio_xtab.c \
			lrec_readers.c \
			lrec_readers.h \
			mmap_byte_reader.c \
			peek_file_reader.c \
			peek_file_reader.h \
			stdio_byte_reader.c \
			string_byte_reader.c
libinput_la_LIBADD=	../lib/libmlr.la
libinput_la_CPPFLAGS=	-I${srcdir}/../

TESTS=			$(check_PROGRAMS)
check_PROGRAMS=		test_byte_readers test_peek_file_reader
# TODO: replace the interesting content with unit tests; jettison the rest
# TODO: causes circular dependency
#noinst_PROGRAMS=	lrim

#lrim_SOURCES=		lrec_reader_in_memory.c
#lrim_CPPFLAGS=		-D__LREC_READER_IN_MEMORY_MAIN__ ${AM_CPPFLAGS}
#lrim_LDADD=		../containers/libcontainers.la ../lib/libmlr.la # -lm

test_byte_readers_CPPFLAGS=	-D__TEST_BYTE_READERS_MAIN__ ${AM_CPPFLAGS}
test_byte_readers_LDADD=	libinput.la \
				../lib/libmlr.la

test_peek_file_reader_CPPFLAGS=	-D__TEST_PEEK_FILE_READER_MAIN__ ${AM_CPPFLAGS}
test_peek_file_reader_LDADD=	libinput.la


AM_CPPFLAGS=	-I${srcdir}/../
AM_CFLAGS=	-std=gnu99
