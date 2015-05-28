# Miller file/record input

These are readers for Miller file formats, stdio and mmap versions. The stdio
and mmap record parsers are similar but not identical, due to inversion of
processing order: getting an entire mallocked line and then splitting it by
separators in the former case, versus spltting while discovering end of line in
the latter case. The code duplication could be largely removed by having the
mmap readers find end-of-lines, then split up the lines -- however that
requires two passes through input strings and for performance I want just a
single pass.

While there are separate record-writers for CSV and pretty-print, there is just
a common record-reader: pretty-print is CSV with field separator being a space,
and `allow_repeat_ifs` set to true.
