* `src/miller/mapping` contains the abstract record-mapper interface datatype, as well as the Go-channel chaining mechanism for piping one mapper into the next.
* `src/miller/mappers` is all the concreate record-mappers such as `cat`, `tac`, `sort`, `put`, and so on. I put it here, not in `mapping`, so all files in `mappers` would be of the same type.
