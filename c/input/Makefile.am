noinst_LTLIBRARIES=	libinput.la
libinput_la_SOURCES=	\
			file_reader_mmap.c \
			file_reader_mmap.h \
			file_reader_stdio.c \
			file_reader_stdio.h \
			lrec_reader.h \
			lrec_reader_in_memory.c \
			lrec_reader_mmap_csvlite.c \
			lrec_reader_mmap_dkvp.c \
			lrec_reader_mmap_nidx.c \
			lrec_reader_mmap_xtab.c \
			lrec_reader_stdio_csv.c \
			lrec_reader_stdio_csvlite.c \
			lrec_reader_stdio_dkvp.c \
			lrec_reader_stdio_nidx.c \
			lrec_reader_stdio_xtab.c \
			lrec_readers.c \
			lrec_readers.h \
			peek_file_reader.c \
			peek_file_reader.h \
			string_byte_reader.c \
			string_byte_reader.h
libinput_la_LIBADD=	../lib/libmlr.la
libinput_la_CPPFLAGS=	-I${srcdir}/../
libinput_la_CFLAGS=	-std=gnu99

TESTS=			$(check_PROGRAMS)
check_PROGRAMS=		test_string_byte_reader

test_string_byte_reader_LDADD=		libinput.la \
					../lib/libmlr.la
test_string_byte_reader_CPPFLAGS=	-D__TEST_STRING_BYTE_READER_MAIN__ ${AM_CPPFLAGS}

AM_CPPFLAGS=	-I${srcdir}/../
