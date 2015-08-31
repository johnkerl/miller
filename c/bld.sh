#gcc -Wall -Werror -I.  lib/mlrutil.c lib/string_builder.c lib/mlr_globals.c containers/slls.c input/peek_file_reader.c experimental/csv0.c

OPT=-O3
#OPT=-g
gcc -Wall -Werror -I. $OPT \
  lib/mlrutil.c \
  lib/mlr_globals.c \
  lib/string_builder.c \
  input/file_reader_mmap.c \
  input/peek_file_reader.c \
  experimental/getlines.c -o getl
