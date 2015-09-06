noinst_LTLIBRARIES=	libinput.la
libinput_la_SOURCES=	\
			byte_reader.h \
			byte_readers.h \
			file_reader_mmap.c \
			file_reader_mmap.h \
			file_reader_stdio.c \
			file_reader_stdio.h \
			lrec_reader.h \
			lrec_reader_csvex.c \
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
			mmap_byte_reader.c \
			old_peek_file_reader.c \
			old_peek_file_reader.h \
			peek_file_reader.c \
			peek_file_reader.h \
			stdio_byte_reader.c \
			string_byte_reader.c
libinput_la_LIBADD=	../lib/libmlr.la
libinput_la_CPPFLAGS=	-I${srcdir}/../
libinput_la_CFLAGS=	-std=gnu99

TESTS=			$(check_PROGRAMS)
check_PROGRAMS=		test_byte_readers

test_byte_readers_LDADD=		libinput.la \
					../lib/libmlr.la
test_byte_readers_CPPFLAGS=	-D__TEST_BYTE_READERS_MAIN__ ${AM_CPPFLAGS}

AM_CPPFLAGS=	-I${srcdir}/../
