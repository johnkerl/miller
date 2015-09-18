# Miller file/record input

These are readers for Miller file formats, stdio and mmap versions. The stdio
and mmap record parsers are similar but not identical, due to inversion of
processing order: getting an entire mallocked line and then splitting it by
separators in the former case, versus splitting while discovering end of line in
the latter case. The code duplication could be largely removed by having the
mmap readers find end-of-lines, then split up the lines -- however that
requires two passes through input strings and for performance I want just a
single pass.

While there are separate record-writers for CSV and pretty-print, there is just
a common record-reader: pretty-print is CSV with field separator being a space,
and `allow_repeat_ifs` set to true.

Idea of `header_keeper` objects for CSV: each `header_keeper` object retains
the input-line backing and the `slls_t` for a CSV header line which is used by
one or more CSV data lines.  Meanwhile some mappers (e.g. `sort`, `tac`) retain
input records from the entire data stream, which may include header-schema
changes in the input stream. This means we need to keep headers intact as long
as any lrecs are pointing to them.  One option is reference-counting which I
experimented with; it was messy and error-prone. The approach used here is to
keep a hash map from header-schema to `header_keeper` object. The current
`pheader_keeper` is a pointer into one of those.  Then when the reader is
freed, all its header-keepers are freed.

There is some code duplication involving single-character and multi-character
IRS, IFS, and IPS. While single-character is a special case of multi-character,
keeping separate implementations for single-character and multi-character
versions is worthwhile for performance. The difference is betweeen `*p == ifs`
and `streqn(p, ifs, ifslen)`: even with function inlining, the latter is more
expensive than the former in the single-character case.

Example timing info for a million-line file is as follows:

```
TIME IN SECONDS 0.945 -- mlr --irs lf   --ifs ,  --ips =  check ../data/big.dkvp2
TIME IN SECONDS 1.139 -- mlr --irs crlf --ifs ,  --ips =  check ../data/big.dkvp2
TIME IN SECONDS 1.291 -- mlr --irs lf   --ifs /, --ips =: check ../data/big.dkvp2
TIME IN SECONDS 1.443 -- mlr --irs crlf --ifs /, --ips =: check ../data/big.dkvp2
```

i.e. (even when averaged over multiple runs) performance improvements of 20-30%
are obtained by special-casing single-character-separator code: this is worth
doing.
