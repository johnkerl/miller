Logic for transforming input records into output records as requested by the user (sort, filter, etc.).:

* The `IRecordTransformer` abstract record-transformer interface datatype, as well as the `ChainTransformer` Go-channel chaining mechanism for piping one transformer into the next.
* The transformer lookup table, used for Miller command-line parsing, verb construction, and online help.
* All the concrete record-transformers such as `cat`, `tac`, `sort`, `put`, and so on.
