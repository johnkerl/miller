#gcc -Wall -Werror -I.  lib/mlrutil.c lib/string_builder.c lib/mlr_globals.c containers/slls.c input/peek_file_reader.c experimental/csv0.c
gcc -Wall -Werror -I. -O3 lib/mlrutil.c lib/mlr_globals.c experimental/getline_for_profile.c -o get1
gcc -Wall -Werror -I. -O3 lib/mlrutil.c lib/mlr_globals.c \
  lib/string_builder.c \
  input/peek_file_reader.c \
  experimental/getmulti_for_profile.c -o get2
