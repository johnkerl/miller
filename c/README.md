# Data flow

Miller data flow is records produced by a record-reader in `input/`, followed
by one or more mappers in `mapping/`, written by a record-writer in `output/`,
controlled by logic in `stream/`. Argument parsing for initial stream setup is
in `cli/`.

# Container names

The user-visible concept of *stream record* (or *srec*) is implemented in the
`lrec_t` (*linked-record type*) data structure. The user-visible concept of
*out-of-stream variables* is implemented using the `mlhmmv_t` (multi-level
hashmap of mlrvals) structure. Source-code comments and names within the code
refer to `srec`/`lrec` and `oosvar`/`mlhmmv` depending on the context.

While those two data structures contain user-visible data structures, others
are used in Miller's implementation: `slls` and `sllv` are singly-linked lists
of string and void-star respectively; `lhmss` is a linked hashmap from string
to string; `lhmsi` is a linked hashmap from string to int; and so on.

# Memory management

Miller is streaming and as near stateless as possible. For most Miller
functions, you can ingest a 20GB file with 4GB RAM, no problem.  For example,
`mlr cat` of a DKVP file retains no data in memory from one line to another;
`mlr cat` of a CSV file retains only the field names from the header line. The
`stats1` and `stats2` commands retain only aggregation state (e.g. count and
sum over specified fields needed to compute mean of specified fields). The `mlr
tac` and `mlr sort` commands, obviously, need to consume and retain all input
records before emitting any output records.

Miller classes are in general modular, following a constructor/destructor model
with minimal dependencies between classes.  As a general rule, void-star
payloads (`sllv`, `lhmslv`) must be freed by the callee (which has access to
the data type) whereas non-void-star payloads (`slls`, `hss`) are freed by the
container class.

One complication is for free-flags in `lrec` and `slls`: the idea is that an
entire line is mallocked and presented by the record reader; then individual
fields are split out and populated into linked list or records. To reduce the
amount of strduping there, free-flags are used to track which fields should be
freed by the destructor and which are freed elsewhere.

The `header_keeper` object is an elaboration on this theme: suppose there is a
CSV file with header line `a,b,c` and data lines `1,2,3`, then `4,5,6`, then
`7,8,9`. Then the keys `a`, `b`, and `c` are shared between all three records;
they are retained in a single `header_keeper` object.

A bigger complication to the otherwise modular nature of Miller is its
*baton-passing memory-management model*. Namely, one class may be responsible
for freeing memory allocated by another class.

For example, using `mlr cat`: The record-reader produces records and returns
pointers to them.  The record-mapper is just a pass-through; it returns the
record-pointers it receives.  The record-writer formats the records to stdout
and does not return them, so it is responsible for freeing them.

Similarly, `mlr cut -x` and any other mappers which modify record objects
without creating new ones. By contrast,`stats1` et al. produce their own
records; they free what they do not pass on.

# Null-lrec conventions

Record-readers return a null lrec-pointer to signify end of input stream.

Each mapper takes an lrec-pointer as input and returns a linked list of
lrec-pointer.

Null-lrec is input to mappers to signify end of stream: e.g. `sort` or `tac`
should use this as a signal to deliver the sorted/reversed list of rows.

When a mapper has no output before end of stream (e.g. `sort` or `tac` while
accumulating inputs) it returns a null lrec-pointer which is treated as
synonymous with returning an empty list.

At end of stream, a mapper returns a linked list of records ending in a null
lrec-pointer.

A null lrec-pointer at end of stream is passed to lrec writers so that they may
produce final output (e.g. pretty-print which produces no output until end of
stream).

# Performance optimizations

The initial implementation of Miller used `lhmss`
(insertion-ordered string-to-string hash map) for record objects.
Keys and values were strduped out of file-input lines. Each of the following
produced from 5 to 30 percent performance gains:
* The `lrec` object is a hashless map suited to low access-to-creation ratio.
See detailed comments in
https://github.com/johnkerl/miller/blob/master/c/containers/lrec.h.
* Free-flags as discussed above removed additional occurrences of string copies.
* Using `mmap` to read files gets rid of double passes on record parsing
(one to find end of line, and another to separate fields) as well as most use
of `malloc`. Note however that standard input cannot be mmapped, so both
record-reader options are retained.

# Source-code indexing

Please see https://sourcegraph.com/github.com/johnkerl/miller
