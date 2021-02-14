Logic for transforming input records into output records as requested by the user (sort, filter, etc.).

* `src/miller/transforming` contains the abstract record-transformer interface datatype, as well as the Go-channel chaining mechanism for piping one transformer into the next.
* `src/miller/transformers` is all the concrete record-transformers such as `cat`, `tac`, `sort`, `put`, and so on. I put it here, not in `transforming`, so all files in `transformers` would be of the same type.
